package k8sinterface

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	//
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// NewKubernetesApi -
func NewKubernetesApiMock() *KubernetesApi {
	InitializeMapResourcesMock()
	return &KubernetesApi{
		KubernetesClient: kubernetesfake.NewSimpleClientset(),
		DynamicClient:    dynamicfake.NewSimpleDynamicClient(&runtime.Scheme{}),
		Context:          context.Background(),
	}
}

func TestListDynamic(t *testing.T) {
	if !IsConnectedToCluster() {
		return
	}
	k8s := NewKubernetesApi()
	// resource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	ww, err := k8s.ListWorkloads2("nginx-ingress", "Deployment")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ww) == 0 {
		t.Error("empty list")
		return
	}
	s, _ := ww[0].GetSelector()
	g, _ := GetGroupVersionResource("pods")
	w, err := k8s.ListWorkloads(&g, "nginx-ingress", s.MatchLabels, nil)
	if len(w) != 1 {
		t.Error("empty list")
		return
	}

}
