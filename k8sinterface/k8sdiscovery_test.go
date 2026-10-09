package k8sinterface

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceGroupToString(t *testing.T) {
	InitializeMapResourcesMock()

	allResources := ResourceGroupToString("*", "*", "*")
	expectedTotal := 0
	for _, gvs := range GetAllResourceGroupMappings() {
		expectedTotal += len(gvs)
	}
	if len(allResources) != expectedTotal {
		t.Errorf("Expected len: %d, received: %d", expectedTotal, len(allResources))
	}
	pod := ResourceGroupToString("*", "*", "Pod")
	if len(pod) == 0 || pod[0] != "/v1/pods" {
		t.Errorf("pod: %v", pod)
	}
	deployments := ResourceGroupToString("*", "*", "Deployment")
	if len(deployments) == 0 || deployments[0] != "apps/v1/deployments" {
		t.Errorf("deployments: %v", deployments)
	}
	cronjobs := ResourceGroupToString("*", "*", "cronjobs")
	if len(cronjobs) == 0 || cronjobs[0] != "batch/v1/cronjobs" {
		t.Errorf("cronjobs: %v", cronjobs)
	}
}

func TestGetGroupVersionResource(t *testing.T) {
	InitializeMapResourcesMock()
	wlid := "wlid://cluster-david-v1/namespace-default/deployment-nginx-deployment"
	r, err := GetGroupVersionResource(wlidpkg.GetKindFromWlid(wlid))
	if err != nil {
		t.Error(err)
		return
	}
	if r.Group != "apps" {
		t.Errorf("wrong group")
	}
	if r.Version != "v1" {
		t.Errorf("wrong Version")
	}
	if r.Resource != "deployments" {
		t.Errorf("wrong Resource")
	}

	r2, err := GetGroupVersionResource("NetworkPolicy")
	if err != nil {
		t.Error(err)
		return
	}
	if r2.Resource != "networkpolicies" {
		t.Errorf("wrong Resource")
	}
}

func TestIsNamespaceScope(t *testing.T) {
	InitializeMapResourcesMock()
	assert.True(t, IsResourceInNamespaceScope("pods"))
	assert.False(t, IsResourceInNamespaceScope("nodes"))
	assert.True(t, IsNamespaceScope(&schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}))
	assert.False(t, IsNamespaceScope(&schema.GroupVersionResource{Group: "", Version: "", Resource: "pods"}))
	assert.True(t, IsNamespaceScope(&schema.GroupVersionResource{Version: "v1", Resource: "serviceaccounts"}))
	assert.False(t, IsNamespaceScope(&schema.GroupVersionResource{Version: "v1", Resource: "nodes"}))
}

func TestInitializeMapResourcesMock(t *testing.T) {

	InitializeMapResourcesMock()
	sampleMap := map[string]string{
		"services":                        "/v1",
		"pods":                            "/v1",
		"replicationcontrollers":          "/v1",
		"podtemplates":                    "/v1",
		"namespaces":                      "/v1",
		"nodes":                           "/v1",
		"configmaps":                      "/v1",
		"secrets":                         "/v1",
		"serviceaccounts":                 "/v1",
		"persistentvolumeclaims":          "/v1",
		"limitranges":                     "/v1",
		"resourcequotas":                  "/v1",
		"daemonsets":                      "apps/v1",
		"deployments":                     "apps/v1",
		"replicasets":                     "apps/v1",
		"statefulsets":                    "apps/v1",
		"controllerrevisions":             "apps/v1",
		"jobs":                            "batch/v1",
		"cronjobs":                        "batch/v1",
		"horizontalpodautoscalers":        "autoscaling/v1",
		"podsecuritypolicies":             "policy/v1beta1",
		"poddisruptionbudgets":            "policy/v1beta1",
		"ingresses":                       "networking.k8s.io/v1",
		"networkpolicies":                 "networking.k8s.io/v1",
		"clusterroles":                    "rbac.authorization.k8s.io/v1",
		"clusterrolebindings":             "rbac.authorization.k8s.io/v1",
		"roles":                           "rbac.authorization.k8s.io/v1",
		"rolebindings":                    "rbac.authorization.k8s.io/v1",
		"mutatingwebhookconfigurations":   "admissionregistration.k8s.io/v1",
		"validatingwebhookconfigurations": "admissionregistration.k8s.io/v1",
	}

	for k, v := range sampleMap {
		v2, ok := GetSingleResourceFromGroupMapping(k)
		assert.True(t, ok)
		assert.Equal(t, v, v2, fmt.Sprintf("resource: %s", k))
	}
}

// TestMultiGroupResource covers resources that are served under more than one
// API group. The mock data exposes "ingresses" under both networking.k8s.io/v1
// and extensions/v1beta1; both must be discoverable.
func TestMultiGroupResource(t *testing.T) {
	InitializeMapResourcesMock()

	gvs, ok := GetResourceFromGroupMapping("ingresses")
	assert.True(t, ok)
	assert.Contains(t, gvs, "networking.k8s.io/v1")
	assert.Contains(t, gvs, "extensions/v1beta1")

	// wildcarded group lookup should return one triplet per group serving the resource
	triplets := ResourceGroupToString("*", "*", "Ingress")
	assert.Contains(t, triplets, "networking.k8s.io/v1/ingresses")
	assert.Contains(t, triplets, "extensions/v1beta1/ingresses")

	// pinning the group should narrow the result to that group only
	pinned := ResourceGroupToString("extensions", "", "Ingress")
	assert.Equal(t, []string{"extensions/v1beta1/ingresses"}, pinned)

	// pinning the version with a wildcarded group must only emit groups whose
	// discovered version matches. extensions only serves v1beta1 in the mock,
	// so a v1 lookup must skip it and only return networking.k8s.io/v1.
	versionPinned := ResourceGroupToString("*", "v1", "Ingress")
	assert.Equal(t, []string{"networking.k8s.io/v1/ingresses"}, versionPinned)

	// asking for a version no group serves should return nothing.
	missingVersion := ResourceGroupToString("*", "v2", "Ingress")
	assert.Empty(t, missingVersion)

	// the v1beta1 lookup must hit extensions and skip networking.k8s.io's v1 entry.
	betaPinned := ResourceGroupToString("*", "v1beta1", "Ingress")
	assert.Equal(t, []string{"extensions/v1beta1/ingresses"}, betaPinned)
}

func TestIsTypeWorkload(t *testing.T) {
	InitializeMapResourcesMock()

	tests := []struct {
		name     string
		object   map[string]interface{}
		expected bool
	}{
		// Canonical workloads
		{
			name:     "CronJob mock object",
			object:   cronJobObjectMock(),
			expected: true,
		},
		{
			name: "Pod core v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
			},
			expected: true,
		},
		{
			name: "Deployment apps/v1",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
			},
			expected: true,
		},
		{
			name: "DaemonSet apps/v1",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "DaemonSet",
			},
			expected: true,
		},
		{
			name: "StatefulSet apps/v1",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "StatefulSet",
			},
			expected: true,
		},
		{
			name: "ReplicaSet apps/v1",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "ReplicaSet",
			},
			expected: true,
		},
		{
			name: "Job batch/v1",
			object: map[string]interface{}{
				"apiVersion": "batch/v1",
				"kind":       "Job",
			},
			expected: true,
		},
		{
			name: "CronJob batch/v1",
			object: map[string]interface{}{
				"apiVersion": "batch/v1",
				"kind":       "CronJob",
			},
			expected: true,
		},
		{
			name: "ReplicationController v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ReplicationController",
			},
			expected: true,
		},
		{
			name: "Deployment extensions/v1beta1 legacy",
			object: map[string]interface{}{
				"apiVersion": "extensions/v1beta1",
				"kind":       "Deployment",
			},
			expected: true,
		},
		{
			name: "Plural kind deployments",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "deployments",
			},
			expected: true,
		},
		// Custom resource workloads with container specs
		{
			name: "CRD with spec.template.spec.containers (e.g. Argo Rollout)",
			object: map[string]interface{}{
				"apiVersion": "argoproj.io/v1alpha1",
				"kind":       "Rollout",
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{"name": "web", "image": "nginx"},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "CRD with direct spec.containers",
			object: map[string]interface{}{
				"apiVersion": "example.com/v1",
				"kind":       "CustomPod",
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{"name": "worker", "image": "alpine"},
					},
				},
			},
			expected: true,
		},
		{
			name: "CRD with spec.jobTemplate.spec.template.spec.containers",
			object: map[string]interface{}{
				"apiVersion": "custom.io/v1",
				"kind":       "CustomScheduledJob",
				"spec": map[string]interface{}{
					"jobTemplate": map[string]interface{}{
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"containers": []interface{}{
										map[string]interface{}{"name": "task", "image": "busybox"},
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		// Non-workload resources (must return false)
		{
			name: "ConfigMap v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
			},
			expected: false,
		},
		{
			name: "Secret v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
			},
			expected: false,
		},
		{
			name: "Service v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
			},
			expected: false,
		},
		{
			name: "ServiceAccount v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ServiceAccount",
			},
			expected: false,
		},
		{
			name: "Namespace v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Namespace",
			},
			expected: false,
		},
		{
			name: "PersistentVolume v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolume",
			},
			expected: false,
		},
		{
			name: "PersistentVolumeClaim v1",
			object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolumeClaim",
			},
			expected: false,
		},
		{
			name: "Role rbac.authorization.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "Role",
			},
			expected: false,
		},
		{
			name: "ClusterRole rbac.authorization.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
			},
			expected: false,
		},
		{
			name: "RoleBinding rbac.authorization.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "RoleBinding",
			},
			expected: false,
		},
		{
			name: "ClusterRoleBinding rbac.authorization.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRoleBinding",
			},
			expected: false,
		},
		{
			name: "NetworkPolicy networking.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "networking.k8s.io/v1",
				"kind":       "NetworkPolicy",
			},
			expected: false,
		},
		{
			name: "Ingress networking.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "networking.k8s.io/v1",
				"kind":       "Ingress",
			},
			expected: false,
		},
		{
			name: "CustomResourceDefinition apiextensions.k8s.io/v1",
			object: map[string]interface{}{
				"apiVersion": "apiextensions.k8s.io/v1",
				"kind":       "CustomResourceDefinition",
			},
			expected: false,
		},
		{
			name: "CRD without containers (e.g. Certificate)",
			object: map[string]interface{}{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"spec": map[string]interface{}{
					"secretName": "cert-secret",
					"dnsNames":   []interface{}{"example.com"},
				},
			},
			expected: false,
		},
		// Invalid / Malformed inputs
		{
			name:     "nil object",
			object:   nil,
			expected: false,
		},
		{
			name:     "empty object",
			object:   map[string]interface{}{},
			expected: false,
		},
		{
			name: "missing kind",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
			},
			expected: false,
		},
		{
			name: "missing apiVersion",
			object: map[string]interface{}{
				"kind": "Deployment",
			},
			expected: false,
		},
		{
			name: "non-string apiVersion",
			object: map[string]interface{}{
				"apiVersion": 123,
				"kind":       "Deployment",
			},
			expected: false,
		},
		{
			name: "non-string kind",
			object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       []interface{}{"Deployment"},
			},
			expected: false,
		},
		{
			name: "empty strings",
			object: map[string]interface{}{
				"apiVersion": "",
				"kind":       "",
			},
			expected: false,
		},
		{
			name: "canonical kind in mismatched group",
			object: map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "Pod",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsTypeWorkload(tt.object)
			assert.Equal(t, tt.expected, actual, "failed for %s", tt.name)
		})
	}
}

func TestUpdateResourceKind(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Pod", "pods"},
		{"Service", "services"},
		{"Node", "nodes"},
		{"Deployment", "deployments"},
		{"NetworkPolicy", "networkpolicies"},
		{"Ingress", "ingresses"},

		{"pod", "pods"},
		{"service", "services"},
		{"node", "nodes"},
		{"deployment", "deployments"},
		{"networkPolicy", "networkpolicies"},
		{"ingress", "ingresses"},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Input: %s", test.input), func(t *testing.T) {
			result := updateResourceKind(test.input)
			if result != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, result)
			}
		})
	}
}
