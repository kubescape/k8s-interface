package cloudsupport

import (
	"encoding/json"
	"fmt"
	"strings"

	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/armosec/k8s-interface/workloadinterface"
	"k8s.io/client-go/tools/clientcmd/api"
)

const TypeCloudProviderDescription workloadinterface.ObjectType = "CloudProviderDescribe" // DEPRECATED
const TypeCloudProviderDescribe workloadinterface.ObjectType = "CloudProviderDescribe"

const (
	ApiVersionEKS                = "eks.amazonaws.com"
	ApiVersionGKE                = "cloud.google.com"
	CloudProviderDescribeKind    = "Describe"
	CloudProviderDescriptionKind = "ClusterDescription" // DEPRECATED
)

// NewDescriptiveInfoFromCloudProvider construct a CloudProviderDescribe from map[string]interface{}. If the map does not match the object, will return nil
func NewDescriptiveInfoFromCloudProvider(object map[string]interface{}) *CloudProviderDescribe {
	if !IsTypeDescriptiveInfoFromCloudProvider(object) {
		return nil
	}

	description := &CloudProviderDescribe{}
	if b := workloadinterface.MapToBytes(object); b != nil {
		if err := json.Unmarshal(b, &description); err != nil {
			return nil
		}
	} else {
		return nil
	}
	description.setProviderFromApiVersion()

	return description
}

func (description *CloudProviderDescribe) setProviderFromApiVersion() {
	if provider := GetCloudProvider(description.GetApiVersion()); provider != "" {
		description.SetProvider(provider)
	}
}
func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if kind, ok := object["kind"]; ok && kind != CloudProviderDescribeKind {
		return true
	}
	return false
}

func IsRunningInCloudProvider() bool {
	currContext := k8sinterface.GetCurrentContext()
	if currContext == nil {
		return false
	}
	if strings.Contains(currContext.Cluster, strings.ToLower("eks")) || strings.Contains(currContext.Cluster, strings.ToLower("gke")) || strings.Contains(currContext.Cluster, strings.ToLower("aks")) {
		return true
	}
	return false
}

func GetCloudProvider(currContext string) string {
	if strings.Contains(currContext, strings.ToLower("eks")) {
		return "eks"
	} else if strings.Contains(currContext, strings.ToLower("gke")) {
		return "gke"
	} else if strings.Contains(currContext, strings.ToLower("aks")) {
		return "aks"
	}
	return ""
}

func GetDescriptiveInfoFromCloudProvider() (workloadinterface.IMetadata, error) {
	currContext := k8sinterface.GetCurrentContext()
	var clusterInfo *CloudProviderDescribe
	var err error
	if currContext == nil {
		return nil, nil
	}
	cloudProvider := GetCloudProvider(currContext.Cluster)
	switch cloudProvider {
	case "eks":
		clusterInfo, err = GetClusterDescribeEKS(currContext)
	case "gke":
		clusterInfo, err = GetClusterDescribeGKE()
	case "aks":
		return nil, fmt.Errorf("we currently do not support reading cloud provider description from aks")
	}

	if err != nil {
		return nil, err
	}
	clusterInfo.SetKind(CloudProviderDescribeKind)
	return clusterInfo, nil
}

// Get descriptive info about cluster running in EKS.
func GetClusterDescribeEKS(currContext *api.Context) (*CloudProviderDescribe, error) {
	eksSupport := NewEKSSupport()

	clusterDescribe, err := eksSupport.getClusterDescribe(currContext)
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(clusterDescribe)
	if err != nil {
		return nil, err
	}

	// set descriptor object
	clusterInfo := &CloudProviderDescribe{}
	clusterInfo.SetApiVersion(ApiVersionEKS)
	clusterInfo.SetName(eksSupport.getName(clusterDescribe))
	clusterInfo.setProviderFromApiVersion()

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}

// Get descriptive info about cluster running in GKE.
func GetClusterDescribeGKE() (*CloudProviderDescribe, error) {
	gkeSupport := newGKESupport()

	clusterDescribe, err := gkeSupport.getClusterDescribe()
	if err != nil {
		return nil, err
	}

	resultInBytes, err := json.Marshal(clusterDescribe)
	if err != nil {
		return nil, err
	}

	clusterInfo := &CloudProviderDescribe{}
	clusterInfo.SetApiVersion(ApiVersionGKE)
	clusterInfo.SetName(gkeSupport.getName(clusterDescribe))
	clusterInfo.setProviderFromApiVersion()

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}
