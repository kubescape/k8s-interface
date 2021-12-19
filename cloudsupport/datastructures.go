package cloudsupport

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
