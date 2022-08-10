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
	TypeCloudProviderDescribe workloadinterface.ObjectType = "CloudProviderDescribe"
	Version                                                = "v1"
	AKS                                                    = "aks"
	GKE                                                    = "gke"
	EKS                                                    = "eks"
)

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

// DEPRECATED - Use apis.IsTypeDescriptiveInfoFromCloudProvider instead
func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if apiVersion, ok := object["apiVersion"]; ok {
		if p, k := apiVersion.(string); k {
			if group := strings.Split(p, "/"); group[0] == apis.ApiVersionGKE || group[0] == apis.ApiVersionEKS {
				if kind, ok := object["kind"]; ok {
					if k, kk := kind.(string); kk && k == apis.CloudProviderDescribeKind || k == "Describe" {
						return true
					}
				}
			}
		}
	}
	return false
}

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
