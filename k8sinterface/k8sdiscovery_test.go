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
	if len(allResources) != len(ResourceGroupMapping) {
		t.Errorf("Expected len: %d, received: %d", len(ResourceGroupMapping), len(allResources))
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
	if len(cronjobs) == 0 || cronjobs[0] != "batch/v1beta1/cronjobs" {
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
		"cronjobs":                        "batch/v1beta1",
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
		v2, ok := ResourceGroupMapping[k]
		assert.True(t, ok)
		assert.Equal(t, v, v2, fmt.Sprintf("resource: %s", k))
	}
}

func TestIsTypeWorkload(t *testing.T) {
	InitializeMapResourcesMock()
	assert.True(t, IsTypeWorkload(cronJobObjectMock()))
}
