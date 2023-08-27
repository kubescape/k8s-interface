package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubescape/k8s-interface/cloudsupport/apis"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
)

const (
	TypeCloudProviderDescribe                workloadinterface.ObjectType = "CloudProviderDescribe"
	TypeCloudProviderDescribeRepositories    workloadinterface.ObjectType = "CloudProviderDescribeRepositories"
	TypeCloudProviderListEntitiesForPolicies workloadinterface.ObjectType = "CloudProviderListEntitiesForPolicies"
	TypeCloudProviderPolicyVersion           workloadinterface.ObjectType = "CloudProviderPolicyVersion"
)

const (
	AKS string = "aks"
	GKE string = "gke"
	EKS string = "eks"
)

const (
	Version         = "v1"
	NotSupportedMsg = "Not supported"
)

// GetCloudProvider get cloud provider name from gitVersion/server URL
func GetCloudProvider() string {
	if IsEKS() {
		return EKS
	}
	if IsGKE() {
		return GKE
	}
	if IsAKS() {
		return AKS
	}
	return ""
}

// NewDescriptiveInfoFromCloudProvider construct a CloudProviderDescribe from map[string]interface{}. If the map does not match the object, will return nil
func NewDescriptiveInfoFromCloudProvider(object map[string]interface{}) *CloudProviderDescribe {
	if !apis.IsTypeDescriptiveInfoFromCloudProvider(object) {
		return nil
	}

	description := &CloudProviderDescribe{}
	if b := workloadinterface.MapToBytes(object); b != nil {
		if err := json.Unmarshal(b, description); err != nil {
			return nil
		}
	} else {
		return nil
	}
	return description
}

// ================================ ListEntitiesForPolicies ================================

func GetListEntitiesForPoliciesEKS(eksSupport IEKSSupport, cluster string, region string) (*CloudProviderListEntitiesForPolicies, error) {
	cluster = eksSupport.GetContextName(cluster)
	// get cluster describe just to get cluster name
	clusterDescribe, err := eksSupport.GetClusterDescribe(cluster, region)
	if err != nil {
		return nil, err
	}
	listEntitiesForPolicies, err := eksSupport.GetListEntitiesForPolicies(region)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(listEntitiesForPolicies)
	if err != nil {
		return nil, err
	}
	// set listEntitiesForPoliciesInfo object
	listEntitiesForPoliciesInfo := &CloudProviderListEntitiesForPolicies{}
	listEntitiesForPoliciesInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version))
	listEntitiesForPoliciesInfo.SetName(eksSupport.GetName(clusterDescribe))
	listEntitiesForPoliciesInfo.SetProvider(EKS)
	listEntitiesForPoliciesInfo.SetKind(apis.CloudProviderListEntitiesForPoliciesKind)

	data := map[string]interface{}{}
	if err := json.Unmarshal(resultInBytes, &data); err != nil {
		return nil, err
	}
	listEntitiesForPoliciesInfo.SetData(data)

	return listEntitiesForPoliciesInfo, nil
}

// GetListEntitiesForPoliciesAKS gets a list of entities for policies (role assignments)
func GetListEntitiesForPoliciesAKS(aksSupport IAKSSupport, cluster string, subscriptionId string, resourceGroup string) (*CloudProviderListEntitiesForPolicies, error) {
	// get cluster describe just to get cluster name
	clusterDescribe, err := aksSupport.GetClusterDescribe(subscriptionId, cluster, resourceGroup)
	if err != nil {
		return nil, err
	}
	scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)
	listEntitiesForPolicies, err := aksSupport.ListAllRolesForScope(subscriptionId, scope)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(listEntitiesForPolicies)
	if err != nil {
		return nil, err
	}
	// set listEntitiesForPoliciesInfo object
	listEntitiesForPoliciesInfo := &CloudProviderListEntitiesForPolicies{}
	listEntitiesForPoliciesInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version))
	listEntitiesForPoliciesInfo.SetName(aksSupport.GetContextName(clusterDescribe))
	listEntitiesForPoliciesInfo.SetProvider(AKS)
	listEntitiesForPoliciesInfo.SetKind(apis.CloudProviderListEntitiesForPoliciesKind)

	data := map[string]interface{}{}
	if err := json.Unmarshal(resultInBytes, &data); err != nil {
		return nil, err
	}
	listEntitiesForPoliciesInfo.SetData(data)

	return listEntitiesForPoliciesInfo, nil
}

// ================================ DescribeRepositories ================================

func GetDescribeRepositoriesEKS(eksSupport IEKSSupport, cluster string, region string) (*CloudProviderDescribeRepositories, error) {
	cluster = eksSupport.GetContextName(cluster)
	// get cluster describe just to get cluster name
	clusterDescribe, err := eksSupport.GetClusterDescribe(cluster, region)
	if err != nil {
		return nil, err
	}
	describeRepositories, err := eksSupport.GetDescribeRepositories(region)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(describeRepositories)
	if err != nil {
		return nil, err
	}
	// set repositoriesInfo object
	repositoriesInfo := &CloudProviderDescribeRepositories{}
	repositoriesInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version))
	repositoriesInfo.SetName(eksSupport.GetName(clusterDescribe))
	repositoriesInfo.SetProvider(EKS)
	repositoriesInfo.SetKind(apis.CloudProviderDescribeRepositoriesKind)

	data := map[string]interface{}{}
	if err := json.Unmarshal(resultInBytes, &data); err != nil {
		return nil, err
	}
	repositoriesInfo.SetData(data)

	return repositoriesInfo, nil
}

// ============================== ClusterDescribe ==============================

// Get descriptive info about cluster running in EKS.
func GetClusterDescribeEKS(eksSupport IEKSSupport, cluster string, region string) (*CloudProviderDescribe, error) {
	cluster = eksSupport.GetContextName(cluster)
	clusterDescribe, err := eksSupport.GetClusterDescribe(cluster, region)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(clusterDescribe)
	if err != nil {
		return nil, err
	}

	// set descriptor object
	clusterInfo := &CloudProviderDescribe{}
	clusterInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version))
	clusterInfo.SetName(eksSupport.GetName(clusterDescribe))
	clusterInfo.SetProvider(EKS)
	clusterInfo.SetKind(apis.CloudProviderDescribeKind)

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}

// Get descriptive info about cluster running in GKE.
func GetClusterDescribeGKE(gkeSupport IGKESupport, clusterName string, region string, project string) (*CloudProviderDescribe, error) {
	cluster := gkeSupport.GetContextName(clusterName)
	clusterDescribe, err := gkeSupport.GetClusterDescribe(cluster, region, project)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(clusterDescribe)
	if err != nil {
		return nil, err
	}

	clusterInfo := &CloudProviderDescribe{}
	clusterInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionGKE, Version))
	clusterInfo.SetName(gkeSupport.GetName(clusterDescribe))
	clusterInfo.SetProvider(GKE)
	clusterInfo.SetKind(apis.CloudProviderDescribeKind)

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}

// Get descriptive info about cluster running in AKS.
func GetClusterDescribeAKS(aksSupport IAKSSupport, cluster string, subscriptionId string, resourceGroup string) (*CloudProviderDescribe, error) {
	clusterDescribe, err := aksSupport.GetClusterDescribe(subscriptionId, cluster, resourceGroup)
	if err != nil {
		return nil, err
	}
	if clusterDescribe == nil {
		return nil, fmt.Errorf("error getting cluster descriptive information")
	}

	resultInBytes, err := json.Marshal(clusterDescribe)
	if err != nil {
		return nil, err
	}

	// set descriptor object
	clusterInfo := &CloudProviderDescribe{}
	clusterInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version))
	clusterInfo.SetName(aksSupport.GetContextName(clusterDescribe))
	clusterInfo.SetProvider(AKS)
	clusterInfo.SetKind(apis.CloudProviderDescribeKind)

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}

func GetPolicyVersionEKS(eksSupport IEKSSupport, cluster string, region string) (*CloudProviderPolicyVersion, error) {
	cluster = eksSupport.GetContextName(cluster)
	// get cluster describe just to get cluster name
	clusterDescribe, err := eksSupport.GetClusterDescribe(cluster, region)
	if err != nil {
		return nil, err
	}
	listPolicyVersion, err := eksSupport.GetPolicyVersion(region)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(listPolicyVersion)
	if err != nil {
		return nil, err
	}
	// set listEntitiesForPoliciesInfo object
	listPolicyInfo := &CloudProviderPolicyVersion{}
	listPolicyInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version))
	listPolicyInfo.SetName(eksSupport.GetName(clusterDescribe))
	listPolicyInfo.SetProvider(EKS)
	listPolicyInfo.SetKind(apis.CloudProviderPolicyVersionKind)

	data := map[string]interface{}{}
	if err := json.Unmarshal(resultInBytes, &data); err != nil {
		return nil, err
	}
	listPolicyInfo.SetData(data)

	return listPolicyInfo, nil
}

// GetPolicyVersionAKS returns a list of all the role definitions that are assigned in this scope.
func GetPolicyVersionAKS(aksSupport IAKSSupport, cluster string, subscriptionId string, resourceGroup string) (*CloudProviderPolicyVersion, error) {
	// get cluster describe just to get cluster name
	clusterDescribe, err := aksSupport.GetClusterDescribe(subscriptionId, cluster, resourceGroup)
	if err != nil {
		return nil, err
	}

	scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)
	listPolicyVersion, err := aksSupport.ListAllRoleDefinitions(subscriptionId, scope)
	if err != nil {
		return nil, err
	}
	if listPolicyVersion == nil {
		return nil, fmt.Errorf("error getting cluster descriptive information")
	}

	resultInBytes, err := json.Marshal(listPolicyVersion)
	if err != nil {
		return nil, err
	}

	// set listEntitiesForPoliciesInfo object
	listPolicyInfo := &CloudProviderPolicyVersion{}
	listPolicyInfo.SetApiVersion(k8sinterface.JoinGroupVersion(apis.ApiVersionAKS, Version))
	listPolicyInfo.SetName(aksSupport.GetContextName(clusterDescribe))
	listPolicyInfo.SetProvider(AKS)
	listPolicyInfo.SetKind(apis.CloudProviderPolicyVersionKind)

	data := map[string]interface{}{}
	if err := json.Unmarshal(resultInBytes, &data); err != nil {
		return nil, err
	}
	listPolicyInfo.SetData(data)

	return listPolicyInfo, nil
}

// check if the server is AKS. e.g. https://XXX.XX.XXX.azmk8s.io:443
func IsAKS() bool {
	const serverIdentifierAKS = "azmk8s.io"
	clusterServerName := k8sinterface.GetK8sConfigClusterServerName()
	return strings.Contains(clusterServerName, serverIdentifierAKS)
}

// check if the server is EKS. e.g. arn:aws:eks:eu-west-1:xxx:cluster/xxxx
func IsEKS() bool {
	version, err := k8sinterface.GetK8SServerGitVersion()
	if err != nil {
		return false
	}
	return strings.Contains(version, EKS)
}

// check if the server is GKE. e.g. gke_xxx-xx-0000_us-central1-c_xxxx-1
func IsGKE() bool {
	version, err := k8sinterface.GetK8SServerGitVersion()
	if err != nil {
		return false
	}
	return strings.Contains(version, GKE)
}
