package apis

import (
	"strings"

	"golang.org/x/exp/slices"
)

const (
	ApiVersionEKS                            = "eks.amazonaws.com"
	ApiVersionAKS                            = "management.azure.com"
	ApiVersionGKE                            = "container.googleapis.com"
	CloudProviderDescribeKind                = "ClusterDescribe"
	CloudProviderDescribeRepositoriesKind    = "DescribeRepositories"
	CloudProviderListEntitiesForPoliciesKind = "ListEntitiesForPolicies"
)

// IsTypeDescriptiveInfoFromCloudProvider return true if the object apiVersion kind match the CloudProviderDescribeKind struct
func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	// "Describe" is deprecated
	return IsCloudProviderType(object, []string{CloudProviderDescribeKind, "Describe"})
}

// IsTypeDescribeRepositories return true if the object apiVersion kind match the CloudProviderDescribeRepositoriesKind struct
func IsTypeDescribeRepositories(object map[string]interface{}) bool {
	return IsCloudProviderType(object, []string{CloudProviderDescribeRepositoriesKind})
}

// IsTypeListEntitiesForPolicies return true if the object apiVersion kind match the CloudProviderListEntitiesForPoliciesKind struct
func IsTypeListEntitiesForPolicies(object map[string]interface{}) bool {
	return IsCloudProviderType(object, []string{CloudProviderListEntitiesForPoliciesKind})
}

func IsCloudProviderType(object map[string]interface{}, acceptableTypes []string) bool {
	if object == nil {
		return false
	}
	if apiVersion, ok := object["apiVersion"]; ok {
		if p, k := apiVersion.(string); k {
			if group := strings.Split(p, "/"); len(group) == 2 {
				if group[0] == ApiVersionGKE || group[0] == ApiVersionEKS || group[0] == ApiVersionAKS {
					if kind, ok := object["kind"]; ok {
						if k, kk := kind.(string); kk && slices.Contains(acceptableTypes, k) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
