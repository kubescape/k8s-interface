package v1

/*
CloudProviderMetadata:
=====================
Metadata of a cloud provider object.
This object may be any configuration object supported by the cloud provider

Name: Object name
Provider: CloudProvider name eks/gke/etc.
*/
type CloudProviderMetadata struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

/*
CloudProviderDescribe:
=========================

CloudProviderDescribe is the desc
*/
type CloudProviderDescribe struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   CloudProviderMetadata  `json:"metadata"`
	Data       map[string]interface{} `json:"data"`
}

const (
	KS_CLOUD_REGION_ENV_VAR = "KS_CLOUD_REGION"
)

/*
CloudProviderDescribeRepositories:
=========================

CloudProviderDescribeRepositories has a list of the image repositories in the cloud provider
*/
type CloudProviderDescribeRepositories struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   CloudProviderMetadata  `json:"metadata"`
	Data       map[string]interface{} `json:"data"`
}

/*
CloudProviderListRolePolicies:
=========================

CloudProviderListRolePolicies has a list of the RolePolicies in the cloud provider (EKS)
*/
type CloudProviderListRolePolicies struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   CloudProviderMetadata  `json:"metadata"`
	Data       map[string]interface{} `json:"data"`
}

/*
CloudProviderListUserPolicies:
=========================

CloudProviderListUserPolicies has a list of the UserPolicies in the cloud provider (EKS)
*/
type CloudProviderListUserPolicies struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   CloudProviderMetadata  `json:"metadata"`
	Data       map[string]interface{} `json:"data"`
}

/*
CloudProviderListGroupPolicies:
=========================

CloudProviderListGroupPolicies has a list of the GroupPolicies in the cloud provider (EKS)
*/
type CloudProviderListGroupPolicies struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   CloudProviderMetadata  `json:"metadata"`
	Data       map[string]interface{} `json:"data"`
}
