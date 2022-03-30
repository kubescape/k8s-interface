package v1

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-04-30/containerservice"
	"github.com/armosec/k8s-interface/cloudsupport/mockobjects"
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

func (AKSSupportM *AKSSupportMock) GetName(managedCluster containerservice.ManagedCluster) string {
	return "daniel"
}

func (AKSSupportM *AKSSupportMock) GetSubscriptionID() (string, error) {
	return "e053c6a9-157e-49c0-818b-461019cb7fef", nil
}

func (AKSSupportM *AKSSupportMock) GetResourceGroup() (string, error) {
	return "armo-dev", nil
}
