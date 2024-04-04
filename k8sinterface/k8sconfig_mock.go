package k8sinterface

import (
	"context"

	// DO NOT REMOVE - load cloud providers auth
	"k8s.io/apimachinery/pkg/runtime"
	discoveryfake "k8s.io/client-go/discovery/fake"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
)

func NewKubernetesApiMock() *KubernetesApi {
	InitializeMapResourcesMock()
	return &KubernetesApi{
		KubernetesClient: kubernetesfake.NewSimpleClientset(),
		DynamicClient:    dynamicfake.NewSimpleDynamicClient(&runtime.Scheme{}),
		DiscoveryClient:  &discoveryfake.FakeDiscovery{},
		Context:          context.Background(),
	}
}
