package cloudsupport

import (
	"fmt"
	"strings"

	cloudsupportv1 "github.com/armosec/k8s-interface/cloudsupport/v1"
	"github.com/armosec/k8s-interface/workloadinterface"
)

const TypeCloudProviderDescription workloadinterface.ObjectType = "CloudProviderDescribe" // DEPRECATED

const (
	CloudProviderDescriptionKind = "ClusterDescription" // DEPRECATED
)

func IsRunningInCloudProvider(cluster string) bool {
	if cluster == "" {
		return false
	}
	if strings.Contains(cluster, strings.ToLower("eks")) || strings.Contains(cluster, strings.ToLower("gke")) || strings.Contains(cluster, strings.ToLower("aks")) {
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

func GetDescriptiveInfoFromCloudProvider(cluster string, cloudProvider string, region string, project string) (workloadinterface.IMetadata, error) {
	var clusterInfo *cloudsupportv1.CloudProviderDescribe
	var err error
	switch cloudProvider {
	case "eks":
		clusterInfo, err = cloudsupportv1.GetClusterDescribeEKS(cloudsupportv1.NewEKSSupport(), cluster, region)
	case "gke":
		clusterInfo, err = cloudsupportv1.GetClusterDescribeGKE(cloudsupportv1.NewGKESupport(), cluster, region, project)
	case "aks":
		return nil, fmt.Errorf("we currently do not support reading cloud provider description from aks")
	}

	if err != nil {
		return nil, err
	}

	return clusterInfo, nil
}
