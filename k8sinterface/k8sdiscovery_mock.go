package k8sinterface

// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#-strong-api-groups-strong-
var ResourceGroupMappingMock = map[string]string{
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
	"ingresses":                       "extensions/v1beta1",
	"podsecuritypolicies":             "policy/v1beta1",
	"poddisruptionbudgets":            "policy/v1",
	"networkpolicies":                 "networking.k8s.io/v1",
	"clusterroles":                    "rbac.authorization.k8s.io/v1",
	"clusterrolebindings":             "rbac.authorization.k8s.io/v1",
	"roles":                           "rbac.authorization.k8s.io/v1",
	"rolebindings":                    "rbac.authorization.k8s.io/v1",
	"mutatingwebhookconfigurations":   "admissionregistration.k8s.io/v1",
	"validatingwebhookconfigurations": "admissionregistration.k8s.io/v1",
}

var ResourceClusterScopeMock = []string{"nodes", "namespaces", "podsecuritypolicies", "clusterroles", "clusterrolebindings", "validatingwebhookconfigurations", "mutatingwebhookconfigurations"}

func InitializeMapResourcesMock() {
	ResourceClusterScope = ResourceClusterScopeMock
	ResourceGroupMapping = ResourceGroupMappingMock
}
