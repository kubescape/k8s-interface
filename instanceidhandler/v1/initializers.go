package instanceidhandler

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/kubescape/k8s-interface/workloadinterface"
	appsv1 "k8s.io/api/apps/v1"
	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/dump"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/util/jsonpath"
)

const (
	Container          = "container"
	InitContainer      = "initContainer"
	EphemeralContainer = "ephemeralContainer"
)

// GenerateInstanceID generates instance ID from workload
// Pods created by a CronJob have an alternate naming convention, using a calculated pod template hash,
// jsonPaths can be specified to drop fields from the pod spec before hashing.
func GenerateInstanceID(w workloadinterface.IWorkload, jsonPaths []string) ([]instanceidhandler.IInstanceID, error) {
	var templateHash string
	if podHash, ok := w.GetLabel("pod-template-hash"); ok && podHash != "" {
		templateHash = podHash
	}
	ownerReferences, err := w.GetOwnerReferences()
	if err != nil {
		return nil, fmt.Errorf("failed to get owner references: %v", err)
	}
	var ownerReference *metav1.OwnerReference
	var alternateName string
	if len(ownerReferences) > 0 {
		ownerReference = &ownerReferences[0]
		switch ownerReference.Kind {
		case "DaemonSet":
			if label, ok := w.GetLabel(appsv1.StatefulSetRevisionLabel); ok {
				alternateName = ownerReference.Name + "-" + label
			}
		case "Job":
			// if the Pod is created by a CronJob, its parent is a Job named after the CronJob
			// with the scheduled timestamp appended to it (unix time in minutes).
			// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/cronjob/utils.go#L277
			// TODO add a json path to exclude some fields
			s := strings.Split(ownerReference.Name, "-")
			if len(s) > 1 && isUnixTimeInMinutes(s[len(s)-1]) {
				// calculate pod template hash
				// https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/controller_utils.go#L1200
				podTemplateSpecHasher := fnv.New32a()
				spec, err := w.GetPodSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to get pod spec: %v", err)
				}
				DeepHashObject(podTemplateSpecHasher, spec, jsonPaths)
				templateHash = rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32()))
				s[len(s)-1] = templateHash
				alternateName = strings.Join(s, "-")
			}
		case "StatefulSet":
			if label, ok := w.GetLabel(appsv1.StatefulSetRevisionLabel); ok {
				alternateName = label
			}
		}
	}

	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	c, err := containerinstance.ListInstanceIDs(ownerReference, containers, Container, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}

	initContainers, err := w.GetInitContainers()
	if err != nil {
		return nil, err
	}

	initC, err := containerinstance.ListInstanceIDs(ownerReference, initContainers, InitContainer, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}
	c = append(c, initC...)

	ephemeralContainers, err := w.GetEphemeralContainers()
	if err != nil {
		return nil, err
	}

	ephemeralC, err := containerinstance.ListInstanceIDs(ownerReference, convertEphemeralToContainers(ephemeralContainers), EphemeralContainer, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}
	c = append(c, ephemeralC...)

	return convertContainersToIInstanceID(c), nil
}

func DeepHashObject(hasher hash.Hash32, pod *core1.PodSpec, extraJsonPaths []string) {
	// sanitize pod sped
	dropProjectedVolumesAndMounts(pod)
	jsonPaths := []string{
		".hostname",
		".containers[*].env[?(@.name==\"DD_INJECT_START_TIME\")]",
	}
	jsonPaths = append(jsonPaths, extraJsonPaths...)
	if err := dropFieldsByJSONPath(pod, jsonPaths); err != nil {
		logger.L().Error("failed to drop fields by JSONPath", helpers.Error(err), helpers.Interface("jsonPath", jsonPaths))
	}
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

func dropFieldsByJSONPath(podSpec *core1.PodSpec, jsonPaths []string) error {
	for _, path := range jsonPaths {
		j := jsonpath.New("dropFields")
		if err := j.Parse(fmt.Sprintf("{%s}", path)); err != nil {
			return fmt.Errorf("failed to parse JSONPath %s: %v", path, err)
		}
		results, err := j.FindResults(podSpec)
		if err != nil {
			return fmt.Errorf("failed to find results for JSONPath %s: %v", path, err)
		}
		for _, values := range results {
			for i := range values {
				v := values[i]
				if v.CanSet() {
					// Get the zero value for the element's type
					zeroVal := reflect.Zero(v.Type())
					// Set the element to its zero value
					v.Set(zeroVal)
				} else {
					fmt.Printf("Element at index %d is not settable\n", i)
				}
			}
		}
	}
	return nil
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
	return GenerateInstanceID(w, nil)
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
