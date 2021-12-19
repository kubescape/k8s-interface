package v1

import (
	"encoding/json"

	"github.com/armosec/k8s-interface/cloudsupport/apis"
	"github.com/armosec/k8s-interface/k8sinterface"
	"github.com/armosec/k8s-interface/workloadinterface"
	"k8s.io/client-go/tools/clientcmd/api"
)

const TypeCloudProviderDescribe workloadinterface.ObjectType = "CloudProviderDescribe"
const Version = "v1"

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
	return description
}

func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if kind, ok := object["kind"]; ok && kind != apis.CloudProviderDescribeKind {
		return true
	}
	return false
}

// Get descriptive info about cluster running in EKS.
func GetClusterDescribeEKS(eksSupport IEKSSupport, currContext *api.Context) (*CloudProviderDescribe, error) {

	clusterDescribe, err := eksSupport.GetClusterDescribe(currContext)
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
	clusterInfo.SetProvider("eks")
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
func GetClusterDescribeGKE(gkeSupport IGKESupport) (*CloudProviderDescribe, error) {

	clusterDescribe, err := gkeSupport.GetClusterDescribe()
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
	clusterInfo.SetProvider("gke")
	clusterInfo.SetKind(apis.CloudProviderDescribeKind)

	data := map[string]interface{}{}
	err = json.Unmarshal(resultInBytes, &data)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetData(data)

	return clusterInfo, nil
}
