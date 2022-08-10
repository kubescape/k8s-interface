package v1

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-04-30/containerservice"
	"github.com/kubescape/k8s-interface/cloudsupport/mockobjects"
)

func NewAKSSupportMock() *AKSSupportMock {
	return &AKSSupportMock{}
}

type AKSSupportMock struct {
}

// Get descriptive info about cluster running in AKS.
func (AKSSupportM *AKSSupportMock) GetClusterDescribe(subscriptionId string, clusterName string, resourceGroup string) (*containerservice.ManagedCluster, error) {
	c := &containerservice.ManagedCluster{}
	err := json.Unmarshal([]byte(mockobjects.AksDescriptor), c)
	return c, err
}

func (AKSSupportM *AKSSupportMock) GetContextName(managedCluster *containerservice.ManagedCluster) string {
	return "daniel"
}

func (AKSSupportM *AKSSupportMock) GetSubscriptionID() (string, error) {
	return "XXXXX", nil
}

func (AKSSupportM *AKSSupportMock) GetResourceGroup() (string, error) {
	return "armo-dev", nil
}
