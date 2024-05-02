package v1

import (
	_ "embed"
	"encoding/json"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/kubescape/k8s-interface/cloudsupport/apis"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/stretchr/testify/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	//go:embed kubeconfig_mock.json
	kubeConfigMock string
	//go:embed kubeconfig_mock_context_not_exist.json
	kubeConfigContextNotExistMock string
)

func getKubeConfigMock() *clientcmdapi.Config {
	kubeConfig := clientcmdapi.Config{}
	if err := json.Unmarshal([]byte(kubeConfigMock), &kubeConfig); err != nil {
		panic(err)
	}
	return &kubeConfig
}

func getkubeConfigContextNotExistMock() *clientcmdapi.Config {
	kubeConfig := clientcmdapi.Config{}
	if err := json.Unmarshal([]byte(kubeConfigContextNotExistMock), &kubeConfig); err != nil {
		panic(err)
	}
	return &kubeConfig
}

func TestGetClusterDescribeGKE(t *testing.T) {
	g := NewGKESupportMock()
	des, err := GetClusterDescribeGKE(g, "kubescape-demo-01", "", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderDescribeKind, des.GetKind())
	assert.Equal(t, "container.googleapis.com/v1/ClusterDescribe/kubescape-demo-01", des.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionGKE, Version), des.GetApiVersion())
	assert.Equal(t, "kubescape-demo-01", des.GetName())
	assert.Equal(t, 34, len(des.GetData()))
}

func TestGetClusterDescribeEKS(t *testing.T) {
	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderDescribeKind, des.GetKind())
	assert.Equal(t, "eks.amazonaws.com/v1/ClusterDescribe/ca-terraform-eks-dev-stage", des.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version), des.GetApiVersion())
	assert.Equal(t, "ca-terraform-eks-dev-stage", des.GetName())
	//assert.Equal(t, 1, len(des.GetData()))
}

// ==================== GetPolicyVersion ====================

func TestGetPolicyVersionAKS(t *testing.T) {
	g := NewAKSSupportMock()
	repos, err := GetPolicyVersionAKS(g, "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderPolicyVersionKind, repos.GetKind())
	assert.Equal(t, "management.azure.com/v1/PolicyVersion/daniel", repos.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version), repos.GetApiVersion())
	assert.Equal(t, "daniel", repos.GetName())
	assert.Equal(t, TypeCloudProviderPolicyVersion, repos.GetObjectType())
	assert.NotEmpty(t, repos.GetData()["roleDefinitions"])
}

func TestGetPolicyVersionEKS(t *testing.T) {
	//TODO: Add more tests
	g := NewEKSSupportMock()
	repos, err := GetPolicyVersionEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderPolicyVersionKind, repos.GetKind())
	assert.Equal(t, "eks.amazonaws.com/v1/PolicyVersion/ca-terraform-eks-dev-stage", repos.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version), repos.GetApiVersion())
	assert.Equal(t, "ca-terraform-eks-dev-stage", repos.GetName())
	assert.Equal(t, TypeCloudProviderPolicyVersion, repos.GetObjectType())
}

func TestSetApiVersionGetPolicyVersion(t *testing.T) {
	repos := &CloudProviderPolicyVersion{ApiVersion: "eks.amazonaws.com/v1"}
	repos.SetApiVersion(apis.ApiVersionEKS)
	assert.Equal(t, repos.ApiVersion, apis.ApiVersionEKS)
}

func TestSetNameGetPolicyVersion(t *testing.T) {
	repos := &CloudProviderPolicyVersion{Metadata: CloudProviderMetadata{Name: "ca-terraform-eks-dev-stage"}}
	repos.SetName("new-name")
	assert.Equal(t, repos.Metadata.Name, "new-name")
}

func TestSetKindGetPolicyVersion(t *testing.T) {
	repos := &CloudProviderPolicyVersion{Kind: "PolicyVersion"}
	repos.SetKind("new-kind")
	assert.Equal(t, repos.Kind, "new-kind")
}

func TestSetObjectGetPolicyVersion(t *testing.T) {
	repos := &CloudProviderPolicyVersion{}
	repos.SetObject(map[string]interface{}{
		"kind": "PolicyVersion", "apiVersion": "eks.amazonaws.com/v1", "metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "PolicyVersion", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
}

func TestSetWorkloadGetPolicyVersion(t *testing.T) {
	repos := &CloudProviderPolicyVersion{}
	repos.SetWorkload(map[string]interface{}{
		"kind": "PolicyVersion", "apiVersion": "eks.amazonaws.com/v1",
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "PolicyVersion", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})

}

// ==================== ListEntitiesForPolicies ====================
func TestGetListEntitiesForPoliciesAKS(t *testing.T) {
	g := NewAKSSupportMock()
	repos, err := GetListEntitiesForPoliciesAKS(g, "XXXXXX", "armo-testing", "armo-dev")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderListEntitiesForPoliciesKind, repos.GetKind())
	assert.Equal(t, "management.azure.com/v1/ListEntitiesForPolicies/daniel", repos.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version), repos.GetApiVersion())
	assert.Equal(t, "daniel", repos.GetName())
	assert.Equal(t, TypeCloudProviderListEntitiesForPolicies, repos.GetObjectType())
	assert.NotEmpty(t, repos.GetData()["roleAssignments"])
}

func TestGetListEntitiesForPoliciesEKS(t *testing.T) {
	//TODO: Add more tests
	g := NewEKSSupportMock()
	repos, err := GetListEntitiesForPoliciesEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderListEntitiesForPoliciesKind, repos.GetKind())
	assert.Equal(t, "eks.amazonaws.com/v1/ListEntitiesForPolicies/ca-terraform-eks-dev-stage", repos.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version), repos.GetApiVersion())
	assert.Equal(t, "ca-terraform-eks-dev-stage", repos.GetName())
	assert.Equal(t, TypeCloudProviderListEntitiesForPolicies, repos.GetObjectType())
}

func TestSetApiVersionListEntitiesForPolicies(t *testing.T) {
	repos := &CloudProviderListEntitiesForPolicies{ApiVersion: "eks.amazonaws.com/v1"}
	repos.SetApiVersion(apis.ApiVersionGKE)
	assert.Equal(t, repos.ApiVersion, apis.ApiVersionGKE)
}

func TestSetNameListEntitiesForPolicies(t *testing.T) {
	repos := &CloudProviderListEntitiesForPolicies{Metadata: CloudProviderMetadata{Name: "ca-terraform-eks-dev-stage"}}
	repos.SetName("new-name")
	assert.Equal(t, repos.Metadata.Name, "new-name")
}

func TestSetKindListEntitiesForPolicies(t *testing.T) {
	repos := &CloudProviderListEntitiesForPolicies{Kind: "ListEntitiesForPolicies"}
	repos.SetKind("new-kind")
	assert.Equal(t, repos.Kind, "new-kind")
}

func TestSetObjectListEntitiesForPolicies(t *testing.T) {
	repos := &CloudProviderListEntitiesForPolicies{}
	repos.SetObject(map[string]interface{}{
		"kind": "ListEntitiesForPolicies", "apiVersion": "eks.amazonaws.com/v1", "metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "ListEntitiesForPolicies", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
}

func TestSetWorkloadListEntitiesForPolicies(t *testing.T) {
	repos := &CloudProviderListEntitiesForPolicies{}
	repos.SetWorkload(map[string]interface{}{
		"kind": "ListEntitiesForPolicies", "apiVersion": "eks.amazonaws.com/v1",
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "ListEntitiesForPolicies", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})

}

// ==================== DescribeRepositories ====================

func TestGetDescribeRepositoriesEKS(t *testing.T) {
	g := NewEKSSupportMock()
	repos, err := GetDescribeRepositoriesEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderDescribeRepositoriesKind, repos.GetKind())
	assert.Equal(t, "eks.amazonaws.com/v1/DescribeRepositories/ca-terraform-eks-dev-stage", repos.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version), repos.GetApiVersion())
	assert.Equal(t, "ca-terraform-eks-dev-stage", repos.GetName())
	assert.Equal(t, TypeCloudProviderDescribeRepositories, repos.GetObjectType())
}

func TestSetApiVersionDescribeRepositories(t *testing.T) {
	repos := &CloudProviderDescribeRepositories{ApiVersion: "eks.amazonaws.com/v1"}
	repos.SetApiVersion(apis.ApiVersionGKE)
	assert.Equal(t, repos.ApiVersion, apis.ApiVersionGKE)
}

func TestSetNameDescribeRepositories(t *testing.T) {
	repos := &CloudProviderDescribeRepositories{Metadata: CloudProviderMetadata{Name: "ca-terraform-eks-dev-stage"}}
	repos.SetName("new-name")
	assert.Equal(t, repos.Metadata.Name, "new-name")
}

func TestSetKindDescribeRepositories(t *testing.T) {
	repos := &CloudProviderDescribeRepositories{Kind: "DescribeRepositories"}
	repos.SetKind("new-kind")
	assert.Equal(t, repos.Kind, "new-kind")
}

func TestSetObjectDescribeRepositories(t *testing.T) {
	repos := &CloudProviderDescribeRepositories{}
	repos.SetObject(map[string]interface{}{
		"kind": "DescribeRepositories", "apiVersion": "eks.amazonaws.com/v1", "metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "DescribeRepositories", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-object", "provider": "b"}})
}

func TestSetWorkloadDescribeRepositories(t *testing.T) {
	repos := &CloudProviderDescribeRepositories{}
	repos.SetWorkload(map[string]interface{}{
		"kind": "DescribeRepositories", "apiVersion": "eks.amazonaws.com/v1",
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})
	assert.Equal(t, repos.GetObject(), map[string]interface{}{
		"kind": "DescribeRepositories", "apiVersion": "eks.amazonaws.com/v1", "data": interface{}(nil),
		"metadata": map[string]interface{}{"name": "new-workload", "provider": "bla"}})

}

// ==================== DescribeCluster ====================

func TestGetClusterDescribeAKS(t *testing.T) {
	defer tearDown()

	g := NewAKSSupportMock()
	clusterDescribe, err := GetClusterDescribeAKS(g, "XXXXXX", "armo-testing", "armo-dev")
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, err)
	d := CloudProviderDescribe{}
	d.SetObject(clusterDescribe.GetObject())

	assert.Equal(t, apis.CloudProviderDescribeKind, d.GetKind())
	assert.Equal(t, "management.azure.com/v1/ClusterDescribe/daniel", d.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version), d.GetApiVersion())
	assert.Equal(t, "daniel", d.GetName())
	assert.Equal(t, AKS, d.GetProvider())

}
func TestNewDescriptiveInfoFromCloudProvider(t *testing.T) {
	defer tearDown()

	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)

	assert.True(t, apis.IsTypeDescriptiveInfoFromCloudProvider(des.GetObject()))
	d := NewDescriptiveInfoFromCloudProvider(des.GetObject())
	assert.NotNil(t, d)

	assert.Equal(t, d.GetID(), des.GetID())

}
func TestSetObjectClusterDescribe(t *testing.T) {
	defer tearDown()

	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, "ca-terraform-eks-dev-stage", "")
	assert.NoError(t, err)

	d := CloudProviderDescribe{}
	d.SetObject(des.GetObject())

	assert.Equal(t, d.GetID(), des.GetID())
	assertMap(t, d.GetObject(), des.GetObject())
}

func assertMap(t *testing.T, expected, actual map[string]interface{}) {
	for k0, v0 := range expected {
		if v1, ok := actual[k0]; ok {
			switch f := v1.(type) {
			case string:
				assert.Equal(t, v0.(string), f)
			case map[string]interface{}:
				assertMap(t, v0.(map[string]interface{}), f)
			}
		}
	}
}

func Test_GetK8sConfigClusterServerName2(t *testing.T) {
	defer tearDown()

	tests := []struct {
		name               string
		clusterContextName string
		clientConfigAPI    *clientcmdapi.Config
		expectedResult     string
	}{
		{
			name:               "Test_GetK8sConfigClusterServerName with valid config",
			clusterContextName: "xxxx-2",
			clientConfigAPI:    getKubeConfigMock(),
			expectedResult:     "https://XXX.XX.XXX.azmk8s.io:443",
		},
		{
			name:               "Test_GetK8sConfigClusterServerName with empty config",
			clusterContextName: "",
			clientConfigAPI:    &clientcmdapi.Config{},
			expectedResult:     "",
		},
		{
			name:               "Test_GetK8sConfigClusterServerName with non-existent context",
			clusterContextName: "non-existent-context",
			clientConfigAPI:    getkubeConfigContextNotExistMock(),
			expectedResult:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sinterface.SetClusterContextName(tt.clusterContextName)
			k8sinterface.SetClientConfigAPI(tt.clientConfigAPI)
			k8sClusterConfigName := k8sinterface.GetK8sConfigClusterServerName()
			assert.Equal(t, tt.expectedResult, k8sClusterConfigName)
		})
	}

}

func Test_GetK8SServerGitVersionNotConnectedToCluster(t *testing.T) {
	defer tearDown()

	k8sinterface.SetK8SGitServerVersion("")
	k8sinterface.SetClusterContextName("no-such-cluster")
	K8SGitServerVersion, err := k8sinterface.GetK8SServerGitVersion()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", K8SGitServerVersion)
}

func TestGetCloudProviderFromNode(t *testing.T) {
	tests := []struct {
		name       string
		providerID string
		want       string
	}{
		{
			name:       "Test AWS",
			providerID: "aws://instance-id",
			want:       EKS,
		},
		{
			name:       "Test GCE",
			providerID: "gce://instance-id",
			want:       GKE,
		},
		{
			name:       "Test Azure",
			providerID: "azure://instance-id",
			want:       AKS,
		},
		{
			name:       "Test DigitalOcean",
			providerID: "digitalocean://instance-id",
			want:       DigitalOcean,
		},
		{
			name:       "Test OpenStack",
			providerID: "openstack://instance-id",
			want:       OpenStack,
		},
		{
			name:       "Test vSphere",
			providerID: "vsphere://instance-id",
			want:       vSphere,
		},
		{
			name:       "Test Oracle",
			providerID: "oci://instance-id",
			want:       Oracle,
		},
		{
			name:       "Test IBM",
			providerID: "ibm://instance-id",
			want:       IBM,
		},
		{
			name:       "Test Unknown Provider",
			providerID: "unknown://instance-id",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &corev1.Node{
				Spec: corev1.NodeSpec{
					ProviderID: tt.providerID,
				},
			}
			if got := GetCloudProviderFromNode(node); got != tt.want {
				t.Errorf("GetCloudProviderFromNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func tearDown() {
	k8sinterface.SetClientConfigAPI(nil)
	k8sinterface.SetClusterContextName("")
	k8sinterface.SetConfigClusterServerName("")
	k8sinterface.SetK8SGitServerVersion("")
}
