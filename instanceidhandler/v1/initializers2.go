package instanceidhandler

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
)

func GetGvkFromRuntimeObj(w runtime.Object) (schema.GroupVersionKind, error) {
	switch object := w.(type) {
	case *appsv1.DaemonSet:
		return schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}, nil
	case *appsv1.Deployment:
		return schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, nil
	case *appsv1.ReplicaSet:
		return schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"}, nil
	case *appsv1.StatefulSet:
		return schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}, nil
	case *batchv1.CronJob:
		return schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"}, nil
	case *batchv1.Job:
		return schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}, nil
	case *v1.Pod:
		return schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, nil
	case *v1.PodTemplate:
		return schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PodTemplate"}, nil
	default:
		return schema.GroupVersionKind{}, fmt.Errorf("unsupported workload type: %s", object)
	}
}

func GetPodSpecFromRuntimeObj(w runtime.Object) (v1.PodSpec, error) {
	switch object := w.(type) {
	case *appsv1.DaemonSet:
		return object.Spec.Template.Spec, nil
	case *appsv1.Deployment:
		return object.Spec.Template.Spec, nil
	case *appsv1.ReplicaSet:
		return object.Spec.Template.Spec, nil
	case *appsv1.StatefulSet:
		return object.Spec.Template.Spec, nil
	case *batchv1.CronJob:
		return object.Spec.JobTemplate.Spec.Template.Spec, nil
	case *batchv1.Job:
		return object.Spec.Template.Spec, nil
	case *v1.Pod:
		return object.Spec, nil
	case *v1.PodTemplate:
		return object.Template.Spec, nil
	default:
		return v1.PodSpec{}, fmt.Errorf("unsupported workload type: %s", object)
	}
}

// GenerateInstanceIDFromRuntimeObj generates instance ID from pod
// Pods created by a CronJob have an alternate naming convention, using a calculated pod template hash,
// jsonPaths can be specified to drop fields from the pod spec before hashing.
func GenerateInstanceIDFromRuntimeObj(w runtime.Object, jsonPaths []string) ([]instanceidhandler.IInstanceID, error) {
	m := w.(metav1.Object)
	gvk, err := GetGvkFromRuntimeObj(w)
	if err != nil {
		return nil, fmt.Errorf("failed to get gvk from workload: %v", err)
	}
	podSpec, err := GetPodSpecFromRuntimeObj(w)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod spec from workload %s/%s: %v", m.GetNamespace(), m.GetName(), err)
	}

	var templateHash string
	if podHash, ok := m.GetLabels()["pod-template-hash"]; ok && podHash != "" {
		templateHash = podHash
	}
	ownerReferences := m.GetOwnerReferences()
	var ownerReference *metav1.OwnerReference
	var alternateName string
	if len(ownerReferences) > 0 {
		ownerReference = &ownerReferences[0]
		switch ownerReference.Kind {
		case "DaemonSet":
			if label, ok := m.GetLabels()[appsv1.StatefulSetRevisionLabel]; ok {
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
				DeepHashObject(podTemplateSpecHasher, &podSpec, jsonPaths)
				templateHash = rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32()))
				s[len(s)-1] = templateHash
				alternateName = strings.Join(s, "-")
			}
		case "StatefulSet":
			if label, ok := m.GetLabels()[appsv1.StatefulSetRevisionLabel]; ok {
				alternateName = label
			}
		}
	}

	c, err := containerinstance.ListInstanceIDs(ownerReference, podSpec.Containers, Container, gvk.GroupVersion().String(), m.GetNamespace(), gvk.Kind, m.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}

	initC, err := containerinstance.ListInstanceIDs(ownerReference, podSpec.InitContainers, InitContainer, gvk.GroupVersion().String(), m.GetNamespace(), gvk.Kind, m.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}
	c = append(c, initC...)

	ephemeralC, err := containerinstance.ListInstanceIDs(ownerReference, convertEphemeralToContainers(podSpec.EphemeralContainers), EphemeralContainer, gvk.GroupVersion().String(), m.GetNamespace(), gvk.Kind, m.GetName(), alternateName, templateHash)
	if err != nil {
		return nil, err
	}
	c = append(c, ephemeralC...)

	return convertContainersToIInstanceID(c), nil
}
