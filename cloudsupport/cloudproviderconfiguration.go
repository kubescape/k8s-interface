package cloudsupport

import (
	"fmt"
	"strings"

	cloudsupportv1 "github.com/armosec/k8s-interface/cloudsupport/v1"
	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/armosec/k8s-interface/workloadinterface"
)

const TypeCloudProviderDescription workloadinterface.ObjectType = "CloudProviderDescribe" // DEPRECATED

const (
	CloudProviderDescriptionKind = "ClusterDescription" // DEPRECATED
)

func IsRunningInCloudProvider() bool {
	currContext := k8sinterface.GetClusterName()
	if currContext == "" {
		return false
	}
	if strings.Contains(currContext, strings.ToLower("eks")) || strings.Contains(currContext, strings.ToLower("gke")) || strings.Contains(currContext, strings.ToLower("aks")) {
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
	currContext := k8sinterface.GetClusterName()
	var clusterInfo *cloudsupportv1.CloudProviderDescribe
	var err error
	if currContext == "" {
		return nil, nil
	}
	cloudProvider := GetCloudProvider(currContext)
	switch cloudProvider {
	case "eks":
		clusterInfo, err = cloudsupportv1.GetClusterDescribeEKS(cloudsupportv1.NewEKSSupport(), currContext)
	case "gke":
		clusterInfo, err = cloudsupportv1.GetClusterDescribeGKE(cloudsupportv1.NewGKESupport())
	case "aks":
		return nil, fmt.Errorf("we currently do not support reading cloud provider description from aks")
	}

	if err != nil {
		return nil, err
	}

	return clusterInfo, nil
}
