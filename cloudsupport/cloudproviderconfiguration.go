package cloudsupport

import (
	"fmt"
	"os"
	"strings"

	cloudsupportv1 "github.com/kubescape/k8s-interface/cloudsupport/v1"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
)

const (
	TypeCloudProviderDescription workloadinterface.ObjectType = "CloudProviderDescribe" // DEPRECATED
	CloudProviderDescriptionKind                              = "ClusterDescription"    // DEPRECATED
	KS_CLOUD_PROVIDER_ENV_VAR                                 = "KS_CLOUD_PROVIDER"
	KS_KUBE_CLUSTER_ENV_VAR                                   = "KS_KUBE_CLUSTER"
)

func IsRunningInCloudProvider(cluster string) bool {
	if cluster == "" {
		return false
	}
	if strings.Contains(cluster, strings.ToLower(cloudsupportv1.EKS)) || strings.Contains(cluster, strings.ToLower(cloudsupportv1.GKE)) || strings.Contains(cluster, strings.ToLower(cloudsupportv1.AKS)) {
		return true
	}
	return false
}

func GetKubeContextName() string {
	val, present := os.LookupEnv(KS_KUBE_CLUSTER_ENV_VAR)
	if present {
		return val
	}

	return k8sinterface.GetContextName()
}

// Try to lookup from env var and then from current context
func GetCloudProvider(currContext string) string {
	val, ok := os.LookupEnv(KS_CLOUD_PROVIDER_ENV_VAR)
	if ok {
		return val
	}
	if strings.Contains(currContext, strings.ToLower(cloudsupportv1.EKS)) {
		return cloudsupportv1.EKS
	} else if strings.Contains(currContext, strings.ToLower(cloudsupportv1.GKE)) {
		return cloudsupportv1.GKE
	} else if strings.Contains(currContext, strings.ToLower(cloudsupportv1.AKS)) {
		return cloudsupportv1.AKS
	}
	return ""
}

// GetDescriptiveInfoFromCloudProvider returns the cluster description from the cloud provider wrapped in IMetadata obj
func GetDescriptiveInfoFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	var clusterInfo *cloudsupportv1.CloudProviderDescribe

	switch cloudProvider {
	case cloudsupportv1.EKS:
		eksSupport := cloudsupportv1.NewEKSSupport()
		region, err := eksSupport.GetRegion(cluster)
		if err != nil {
			return nil, err
		}
		clusterInfo, err = cloudsupportv1.GetClusterDescribeEKS(eksSupport, cluster, region)
		if err != nil {
			return nil, err
		}
	case cloudsupportv1.GKE:
		gkeSupport := cloudsupportv1.NewGKESupport()
		project, err := gkeSupport.GetProject(cluster)
		if err != nil {
			return nil, err
		}
		region, err := gkeSupport.GetRegion(cluster)
		if err != nil {
			return nil, err
		}
		clusterInfo, err = cloudsupportv1.GetClusterDescribeGKE(gkeSupport, cluster, region, project)
		if err != nil {
			return nil, err
		}
	case cloudsupportv1.AKS:
		aksSupport := cloudsupportv1.NewAKSSupport()
		subscriptionID, err := aksSupport.GetSubscriptionID()
		if err != nil {
			return nil, err
		}
		resourceGroup, err := aksSupport.GetResourceGroup()
		if err != nil {
			return nil, err
		}
		clusterInfo, err = cloudsupportv1.GetClusterDescribeAKS(aksSupport, cluster, subscriptionID, resourceGroup)
		if err != nil {
			return nil, err
		}
	}

	return clusterInfo, nil
}

// GetDescribeRepositoriesFromCloudProvider returns image repository descriptions from the cloud provider wrapped in IMetadata obj
func GetDescribeRepositoriesFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	var clusterInfo *cloudsupportv1.CloudProviderDescribeRepositories

	switch cloudProvider {
	case cloudsupportv1.EKS:
		eksSupport := cloudsupportv1.NewEKSSupport()
		region, err := eksSupport.GetRegion(cluster)
		if err != nil {
			return nil, err
		}
		clusterInfo, err = cloudsupportv1.GetDescribeRepositoriesEKS(eksSupport, cluster, region)
		if err != nil {
			return nil, err
		}
	case cloudsupportv1.GKE:
		//TODO - implement GKE support
		return nil, fmt.Errorf("GKE is not supported yet")
	case cloudsupportv1.AKS:
		//TODO - implement AKS support
		return nil, fmt.Errorf("AKS is not supported yet")
	}

	return clusterInfo, nil
}
