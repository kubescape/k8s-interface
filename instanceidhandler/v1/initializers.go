package instanceidhandler

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"strconv"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/kubescape/k8s-interface/workloadinterface"
	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/dump"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	Container          = "container"
	InitContainer      = "initContainer"
	EphemeralContainer = "ephemeralContainer"
)

// GenerateInstanceID generates instance ID from workload
func GenerateInstanceID(w workloadinterface.IWorkload) ([]instanceidhandler.IInstanceID, error) {
	ownerReferences, err := w.GetOwnerReferences()
	if err != nil {
		return nil, fmt.Errorf("failed to get owner references: %v", err)
	}
	var ownerReference *metav1.OwnerReference
	var alternateName string
	if len(ownerReferences) > 0 {
		ownerReference = &ownerReferences[0]
		// if the Pod is created by a CronJob, its parent is a Job named after the CronJob
		// with the scheduled timestamp appended to it (unix time in minutes).
		// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/cronjob/utils.go#L277
		if ownerReference.Kind == "Job" {
			s := strings.Split(ownerReference.Name, "-")
			if len(s) > 1 && isUnixTimeInMinutes(s[len(s)-1]) {
				// calculate pod template hash
				// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/controller_utils.go#L1200
				podTemplateSpecHasher := fnv.New32a()
				spec, err := w.GetPodSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to get pod spec: %v", err)
				}
				DeepHashObject(podTemplateSpecHasher, spec)
				s[len(s)-1] = rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32()))
				alternateName = strings.Join(s, "-")
			}
		}
	}

	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	c, err := containerinstance.ListInstanceIDs(ownerReference, containers, Container, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName)
	if err != nil {
		return nil, err
	}

	initContainers, err := w.GetInitContainers()
	if err != nil {
		return nil, err
	}

	initC, err := containerinstance.ListInstanceIDs(ownerReference, initContainers, InitContainer, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName)
	if err != nil {
		return nil, err
	}
	c = append(c, initC...)

	ephemeralContainers, err := w.GetEphemeralContainers()
	if err != nil {
		return nil, err
	}

	ephemeralC, err := containerinstance.ListInstanceIDs(ownerReference, convertEphemeralToContainers(ephemeralContainers), EphemeralContainer, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName)
	if err != nil {
		return nil, err
	}
	c = append(c, ephemeralC...)

	return convertContainersToIInstanceID(c), nil
}

func DeepHashObject(hasher hash.Hash32, pod *core1.PodSpec) {
	// sanitize pod sped
	dropProjectedVolumesAndMounts(pod)
	dropHostname(pod)
	hasher.Reset()
	_, _ = fmt.Fprintf(hasher, "%v", dump.ForHash(pod))
}

// https://github.com/kubernetes/autoscaler/blob/5077f4817bc3105bd61659950833bd6d66edc556/cluster-autoscaler/utils/utils.go#L76
func dropProjectedVolumesAndMounts(podSpec *core1.PodSpec) {
	projectedVolumeNames := map[string]bool{}
	var volumes []core1.Volume
	for _, v := range podSpec.Volumes {
		if v.Projected == nil {
			volumes = append(volumes, v)
		} else {
			projectedVolumeNames[v.Name] = true
		}
	}
	podSpec.Volumes = volumes

	for i := range podSpec.Containers {
		var volumeMounts []core1.VolumeMount
		for _, mount := range podSpec.Containers[i].VolumeMounts {
			if ok := projectedVolumeNames[mount.Name]; !ok {
				volumeMounts = append(volumeMounts, mount)
			}
		}
		podSpec.Containers[i].VolumeMounts = volumeMounts
	}

	for i := range podSpec.InitContainers {
		var volumeMounts []core1.VolumeMount
		for _, mount := range podSpec.InitContainers[i].VolumeMounts {
			if ok := projectedVolumeNames[mount.Name]; !ok {
				volumeMounts = append(volumeMounts, mount)
			}
		}
		podSpec.InitContainers[i].VolumeMounts = volumeMounts
	}
}

// https://github.com/kubernetes/autoscaler/blob/5077f4817bc3105bd61659950833bd6d66edc556/cluster-autoscaler/utils/utils.go#L109
func dropHostname(podSpec *core1.PodSpec) {
	podSpec.Hostname = ""
}

func isUnixTimeInMinutes(s string) bool {
	if i, err := strconv.Atoi(s); err == nil {
		return i > 0 && i < math.MaxInt64/60
	}
	return false
}

// GenerateInstanceIDFromPod generates instance ID from pod
// Deprecated: use GenerateInstanceID instead
func GenerateInstanceIDFromPod(pod *core1.Pod) ([]instanceidhandler.IInstanceID, error) {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	if err != nil {
		return nil, fmt.Errorf("convert pod to unstructured: %v", err)
	}
	w := workloadinterface.NewWorkloadObj(unstructuredObj)
	return GenerateInstanceID(w)
}

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/...
func GenerateInstanceIDFromString(input string) (instanceidhandler.IInstanceID, error) {
	return containerinstance.GenerateInstanceIDFromString(input)
}

// convert list containerinstance.InstanceID to instanceidhandler.IInstanceID
func convertContainersToIInstanceID(l []containerinstance.InstanceID) []instanceidhandler.IInstanceID {
	li := make([]instanceidhandler.IInstanceID, len(l))
	for i := range l {
		li[i] = &l[i]
	}
	return li
}

func convertEphemeralToContainers(e []core1.EphemeralContainer) []core1.Container {
	c := make([]core1.Container, len(e))
	for i := range e {
		c[i] = core1.Container(e[i].EphemeralContainerCommon)
	}
	return c
}
