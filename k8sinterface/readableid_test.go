package k8sinterface

import (
	"testing"

	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
)

func TestGetReadableID(t *testing.T) {
	tests := []struct {
		name       string
		obj        map[string]interface{}
		expectedID string
	}{
		{
			name: "namespaced resource with group in apiVersion",
			obj: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name":      "nginx",
					"namespace": "default",
				},
			},
			expectedID: "apps/v1/default/Deployment/nginx",
		},
		{
			name: "cluster-scoped resource has no namespace segment",
			obj: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
				"metadata": map[string]interface{}{
					"name": "admin",
				},
			},
			expectedID: "rbac.authorization.k8s.io/v1/ClusterRole/admin",
		},
		{
			name: "missing apiVersion is omitted",
			obj: map[string]interface{}{
				"kind": "Pod",
				"metadata": map[string]interface{}{
					"name":      "nginx",
					"namespace": "default",
				},
			},
			expectedID: "default/Pod/nginx",
		},
		{
			name:       "empty object still returns a slash-separated ID",
			obj:        map[string]interface{}{},
			expectedID: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := workloadinterface.NewWorkloadObj(tt.obj)
			assert.Equal(t, tt.expectedID, GetReadableID(obj))
		})
	}
}
