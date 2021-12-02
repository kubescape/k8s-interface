package k8sinterface

import (
	"context"

	cautils "github.com/armosec/utils-go/utils"
	"github.com/armosec/utils-k8s-go/armometadata"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func IsAttached(labels map[string]string) *bool {
	return IsLabel(labels, armometadata.ArmoAttach)
}

func IsAgentCompatibleLabel(labels map[string]string) *bool {
	return IsLabel(labels, armometadata.ArmoCompatibleLabel)
}
func IsAgentCompatibleAnnotation(annotations map[string]string) *bool {
	return IsLabel(annotations, armometadata.ArmoCompatibleAnnotation)
}
func SetAgentCompatibleLabel(labels map[string]string, val bool) {
	SetLabel(labels, armometadata.ArmoCompatibleLabel, val)
}
func SetAgentCompatibleAnnotation(annotations map[string]string, val bool) {
	SetLabel(annotations, armometadata.ArmoCompatibleAnnotation, val)
}
func IsLabel(labels map[string]string, key string) *bool {
	if len(labels) == 0 {
		return nil
	}
	var k bool
	if l, ok := labels[key]; ok {
		if cautils.StringToBool(l) {
			k = true
		} else if !cautils.StringToBool(l) {
			k = false
		}
		return &k
	}
	return nil
}
func SetLabel(labels map[string]string, key string, val bool) {
	if labels == nil {
		return
	}
	labels[key] = cautils.BoolToString(val)
}
func (k8sAPI *KubernetesApi) ListAttachedPods(namespace string) ([]corev1.Pod, error) {
	return k8sAPI.ListPods(namespace, map[string]string{armometadata.ArmoAttach: cautils.BoolToString(true)})
}

func (k8sAPI *KubernetesApi) ListPods(namespace string, podLabels map[string]string) ([]corev1.Pod, error) {
	listOptions := metav1.ListOptions{}
	if len(podLabels) > 0 {
		set := labels.Set(podLabels)
		listOptions.LabelSelector = set.AsSelector().String()
	}
	pods, err := k8sAPI.KubernetesClient.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		return []corev1.Pod{}, err
	}
	return pods.Items, nil
}
