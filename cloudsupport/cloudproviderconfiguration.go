package cloudsupport

import (
	"errors"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"

	cloudsupportv1 "github.com/kubescape/k8s-interface/cloudsupport/v1"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
)

const (
	TypeCloudProviderDescription workloadinterface.ObjectType = "CloudProviderDescribe" // DEPRECATED
	CloudProviderDescriptionKind                              = "ClusterDescription"    // DEPRECATED
	KS_KUBE_CLUSTER_ENV_VAR                                   = "KS_KUBE_CLUSTER"
	KS_OFFLINE_ENV_VAR                                        = "KS_OFFLINE"
)

// ErrCloudDescribeUnavailable is returned by the Get*FromCloudProvider entry
// points when cluster cloud-describe data cannot be obtained but the failure
// must not abort the scan (e.g. air-gapped environments, missing creds, or
// KS_OFFLINE=true). Callers should recognise it via errors.Is and continue
// collecting the remaining host-scanner / node-agent data.
var ErrCloudDescribeUnavailable = errors.New("cloud describe unavailable")

// cloudDescribeDisabled reports whether cloud-describe should be skipped
// entirely. Set by the Helm chart when capabilities.kubescapeOffline=enable.
func cloudDescribeDisabled() bool {
	return os.Getenv(KS_OFFLINE_ENV_VAR) == "true"
}

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

// GetCloudProvider returns the cloud provider name
func GetCloudProvider(nodeList *corev1.NodeList) string {
	return cloudsupportv1.GetCloudProvider(nodeList)
}

// GetDescriptiveInfoFromCloudProvider returns the cluster description from the cloud provider wrapped in IMetadata obj
func GetDescriptiveInfoFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	if cloudDescribeDisabled() {
		return nil, fmt.Errorf("%w: %s=true", ErrCloudDescribeUnavailable, KS_OFFLINE_ENV_VAR)
	}

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
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		resourceGroup, err := aksSupport.GetResourceGroup()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		clusterInfo, err = cloudsupportv1.GetClusterDescribeAKS(aksSupport, cluster, subscriptionID, resourceGroup)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
	default:
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	}

	return clusterInfo, nil
}

// GetDescribeRepositoriesFromCloudProvider returns image repository descriptions from the cloud provider wrapped in IMetadata obj
func GetDescribeRepositoriesFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	if cloudDescribeDisabled() {
		return nil, fmt.Errorf("%w: %s=true", ErrCloudDescribeUnavailable, KS_OFFLINE_ENV_VAR)
	}

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
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	case cloudsupportv1.AKS:
		//TODO - implement AKS support
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	default:
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	}

	return clusterInfo, nil
}

// GetListEntitiesForPoliciesFromCloudProvider returns EntitiesForpolicies from the cloud provider wrapped in IMetadata obj
func GetListEntitiesForPoliciesFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	if cloudDescribeDisabled() {
		return nil, fmt.Errorf("%w: %s=true", ErrCloudDescribeUnavailable, KS_OFFLINE_ENV_VAR)
	}

	var listEntitiesForPolicies *cloudsupportv1.CloudProviderListEntitiesForPolicies

	switch cloudProvider {
	case cloudsupportv1.EKS:
		eksSupport := cloudsupportv1.NewEKSSupport()
		region, err := eksSupport.GetRegion(cluster)
		if err != nil {
			return nil, err
		}
		listEntitiesForPolicies, err = cloudsupportv1.GetListEntitiesForPoliciesEKS(eksSupport, cluster, region)
		if err != nil {
			return nil, err
		}
	case cloudsupportv1.GKE:
		//TODO - implement GKE support
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	case cloudsupportv1.AKS:
		aksSupport := cloudsupportv1.NewAKSSupport()
		subscriptionID, err := aksSupport.GetSubscriptionID()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		resourceGroup, err := aksSupport.GetResourceGroup()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		listEntitiesForPolicies, err = cloudsupportv1.GetListEntitiesForPoliciesAKS(aksSupport, cluster, subscriptionID, resourceGroup)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
	default:
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	}

	return listEntitiesForPolicies, nil
}

// GetPolicyVersionFromCloudProvider returns PolicyVersion from the cloud provider wrapped in IMetadata obj
func GetPolicyVersionFromCloudProvider(cluster string, cloudProvider string) (workloadinterface.IMetadata, error) {
	if cloudDescribeDisabled() {
		return nil, fmt.Errorf("%w: %s=true", ErrCloudDescribeUnavailable, KS_OFFLINE_ENV_VAR)
	}

	var policyVersion *cloudsupportv1.CloudProviderPolicyVersion

	switch cloudProvider {
	case cloudsupportv1.EKS:
		eksSupport := cloudsupportv1.NewEKSSupport()
		region, err := eksSupport.GetRegion(cluster)
		if err != nil {
			return nil, err
		}
		policyVersion, err = cloudsupportv1.GetPolicyVersionEKS(eksSupport, cluster, region)
		if err != nil {
			return nil, err
		}
	case cloudsupportv1.GKE:
		//TODO - implement GKE support
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	case cloudsupportv1.AKS:
		aksSupport := cloudsupportv1.NewAKSSupport()
		subscriptionID, err := aksSupport.GetSubscriptionID()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		resourceGroup, err := aksSupport.GetResourceGroup()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
		policyVersion, err = cloudsupportv1.GetPolicyVersionAKS(aksSupport, cluster, subscriptionID, resourceGroup)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCloudDescribeUnavailable, err)
		}
	default:
		return nil, fmt.Errorf(cloudsupportv1.NotSupportedMsg)
	}

	return policyVersion, nil
}
