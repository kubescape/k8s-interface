package k8sinterface

import (
	"context"

	"github.com/armosec/utils-go/boolutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func IsLabel(labels map[string]string, key string) *bool {
	if len(labels) == 0 {
		return nil
	}
	var k bool
	if l, ok := labels[key]; ok {
		if boolutils.StringToBool(l) {
			k = true
		} else if !boolutils.StringToBool(l) {
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
	labels[key] = boolutils.BoolToString(val)
}

func (k8sAPI *KubernetesApi) ListPods(namespace string, podLabels map[string]string) (*corev1.PodList, error) {
	listOptions := metav1.ListOptions{}
	if len(podLabels) > 0 {
		set := labels.Set(podLabels)
		listOptions.LabelSelector = set.AsSelector().String()
	}
	pods, err := k8sAPI.KubernetesClient.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	return pods, nil
}
