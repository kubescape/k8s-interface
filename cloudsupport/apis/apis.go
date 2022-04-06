package apis

import "strings"

const (
	ApiVersionEKS             = "eks.amazonaws.com"
	ApiVersionAKS             = "management.azure.com"
	ApiVersionGKE             = "container.googleapis.com"
	CloudProviderDescribeKind = "ClusterDescribe"
)

// IsTypeDescriptiveInfoFromCloudProvider return true if the object apiVersion kind match the CloudProviderDescribeKind struct
func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if apiVersion, ok := object["apiVersion"]; ok {
		if p, k := apiVersion.(string); k {
			if group := strings.Split(p, "/"); group[0] == ApiVersionGKE || group[0] == ApiVersionEKS || group[0] == ApiVersionAKS {
				if kind, ok := object["kind"]; ok {
					// "Describe" is deprecated
					if k, kk := kind.(string); kk && k == CloudProviderDescribeKind || k == "Describe" {
						return true
					}
				}
			}
		}
	}
	return false
}
