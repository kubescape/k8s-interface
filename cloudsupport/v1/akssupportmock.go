package v1

import (
	"encoding/json"

	armcontainerservice "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/kubescape/k8s-interface/cloudsupport/mockobjects"
	"github.com/kubescape/k8s-interface/k8sinterface"
)

func NewAKSSupportMock() *AKSSupportMock {
	return &AKSSupportMock{}
}

type AKSSupportMock struct {
}

// Get descriptive info about cluster running in AKS.
func (AKSSupportM *AKSSupportMock) GetClusterDescribe(subscriptionId string, clusterName string, resourceGroup string) (*armcontainerservice.ManagedCluster, error) {
	c := &armcontainerservice.ManagedCluster{}
	err := json.Unmarshal([]byte(mockobjects.AksDescriptor), c)
	return c, err
}

func (AKSSupportM *AKSSupportMock) ListAllRolesForScope(subscriptionId string, scope string) (*ListRoleAssignment, error) {
	c := &ListRoleAssignment{}
	err := json.Unmarshal([]byte(mockobjects.AKSListRoleAssignments), c)
	return c, err
}

func (AKSSupportM *AKSSupportMock) ListAllRoleDefinitions(subscriptionId string, scope string) (*ListRoleDefinition, error) {
	c := &ListRoleDefinition{}
	err := json.Unmarshal([]byte(mockobjects.AKSListRoleDefinitions), c)
	return c, err
}
func (AKSSupportM *AKSSupportMock) GetContextName(managedCluster *armcontainerservice.ManagedCluster) string {
	return "daniel"
}

func (AKSSupportM *AKSSupportMock) GetSubscriptionID() (string, error) {
	return "XXXXX", nil
}

func (AKSSupportM *AKSSupportMock) GetResourceGroup() (string, error) {
	return "armo-dev", nil
}

func (AKSSupportM *AKSSupportMock) GetGroupIdsRoleBindings(kapi *k8sinterface.KubernetesApi, namespace string) ([]string, error) {
	return []string{"e808215d-d159-49ba-8bb6-9661ba478842", "unexpected comma, expecting type"}, nil
}
