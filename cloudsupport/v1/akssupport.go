package v1

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-04-30/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	AZURE_SUBSCRIPTION_ID_ENV_VAR = "AZURE_SUBSCRIPTION_ID"
	AZURE_RESOURCE_GROUP_ENV_VAR  = "AZURE_RESOURCE_GROUP"
)

type IAKSSupport interface {
	GetClusterDescribe(subscriptionId string, clusterName string, resourceGroup string) (*containerservice.ManagedCluster, error)
	GetContextName(*containerservice.ManagedCluster) string
	GetSubscriptionID() (string, error)
	GetResourceGroup() (string, error)
}
type AKSSupport struct {
}

func NewAKSSupport() *AKSSupport {
	return &AKSSupport{}
}

// Get descriptive info about cluster running in AKS.
func (AKSSupport *AKSSupport) GetClusterDescribe(subscriptionId string, clusterName string, resourceGroup string) (*containerservice.ManagedCluster, error) {
	// see https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	// set AKS managed cluster client
	aksClient := containerservice.NewManagedClustersClient(subscriptionId)
	aksClient.Authorizer = authorizer

	// Get cluster
	resp, err := aksClient.Get(context.Background(), resourceGroup, clusterName)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}

func (AKSSupport *AKSSupport) GetContextName(managedCluster *containerservice.ManagedCluster) string {
	if managedCluster != nil {
		if managedCluster.Name != nil {
			return *managedCluster.Name
		}
	}
	return ""
}

func (AKSSupport *AKSSupport) GetSubscriptionID() (string, error) {
	if subscriptionId, ok := os.LookupEnv(AZURE_SUBSCRIPTION_ID_ENV_VAR); ok {
		return subscriptionId, nil
	}
	return "", fmt.Errorf("error retrieving azure subscription id: environment variable %s not set", AZURE_SUBSCRIPTION_ID_ENV_VAR)
}

func (AKSSupport *AKSSupport) GetResourceGroup() (string, error) {
	if subscriptionId, ok := os.LookupEnv(AZURE_RESOURCE_GROUP_ENV_VAR); ok {
		return subscriptionId, nil
	}
	return "", fmt.Errorf("error retrieving azure subscription id: environment variable %s not set", AZURE_RESOURCE_GROUP_ENV_VAR)
}
