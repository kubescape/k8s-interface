package cloudsupport

import (
	"encoding/json"
	"fmt"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"k8s.io/apimachinery/pkg/version"
)

type ApiServerMetadata struct {
	Name     string `json:"name"`               // Name of the info object. e.g. "version"
	Provider string `json:"provider,omitempty"` // Provider of the info object. e.g. "kubernetes"
}

const (
	TypeApiServerInfo        workloadinterface.ObjectType = "APIServerInfo"
	ApiServerInfoGroup       string                       = "apiserverinfo.kubescape.cloud"
	ApiServerInfoVersion     string                       = "v1beta0"
	ApiServerInfoApiVersion  string                       = ApiServerInfoGroup + "/" + ApiServerInfoVersion
	apiServerInfoVersionName string                       = "version" // Name of the ApiServerInfo object that holds the version
)

// ApiServerInfo is a struct that holds the information about the api server.
// It implements the IMetadata interface.
//
// An example of the object is:
// 		{
// 		    "apiVersion": "apiserverinfo.kubescape.cloud/v1beta0",
// 		    "kind": "APIServerInfo",
// 		    "metadata": {
// 		        "name": "version"
// 		    },
// 		    "data": {
// 		        "major": "1",
// 		        "minor": "22",
// 		        "gitVersion": "v1.22.11-gke.400",
// 		        "gitCommit": "b4e1ab06be827438def8aee0021791b413a0961d",
// 		        "gitTreeState": "clean",
// 		        "buildDate": "2022-06-24T09:27:38Z",
// 		        "goVersion": "go1.16.15b7",
// 		        "compiler": "gc",
// 		        "platform": "linux/amd64"
// 		    }
// 		}
type ApiServerInfo struct {
	ApiVersion string            `json:"apiVersion"` // apiVersion: apiserverinfo.kubescape.cloud/v1beta0
	Kind       string            `json:"kind"`       // kind: APIServerInfo
	Metadata   ApiServerMetadata `json:"metadata"`
	Data       interface{}       `json:"data"` // currently only *version.Info
}

// ApiserverMetadata getters / setters
func (apiServerMetadata *ApiServerMetadata) GetName() string     { return apiServerMetadata.Name }
func (apiServerMetadata *ApiServerMetadata) GetProvider() string { return apiServerMetadata.Provider }
func (apiServerMetadata *ApiServerMetadata) SetName(name string) { apiServerMetadata.Name = name }
func (apiServerMetadata *ApiServerMetadata) SetProvider(provider string) {
	apiServerMetadata.Provider = provider
}

// ApiServerInfo implements the Imetadata interface
// Getters
func (apiServerInfo *ApiServerInfo) GetNamespace() string  { return "" }
func (apiServerInfo *ApiServerInfo) GetName() string       { return apiServerInfo.Metadata.GetName() }
func (apiServerInfo *ApiServerInfo) GetKind() string       { return apiServerInfo.Kind }
func (apiServerInfo *ApiServerInfo) GetApiVersion() string { return apiServerInfo.ApiVersion }
func (apiServerInfo *ApiServerInfo) GetWorkload() map[string]interface{} {
	return apiServerInfo.GetObject()
}

func (apiServerInfo *ApiServerInfo) GetObjectType() workloadinterface.ObjectType {
	return TypeApiServerInfo
}

func (apiServerInfo *ApiServerInfo) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*apiServerInfo)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

func (apiServerInfo *ApiServerInfo) GetID() string {
	return fmt.Sprintf("%s/%s/%s",
		k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(apiServerInfo.GetApiVersion())),
		apiServerInfo.GetKind(),
		apiServerInfo.GetName(),
	)
}

// Setters
func (apiServerInfo *ApiServerInfo) SetNamespace(namespace string)               {}
func (apiServerInfo *ApiServerInfo) SetName(name string)                         { apiServerInfo.Metadata.SetName(name) }
func (apiServerInfo *ApiServerInfo) SetKind(kind string)                         { apiServerInfo.Kind = kind }
func (apiServerInfo *ApiServerInfo) SetWorkload(workload map[string]interface{}) { /* Deprecated */ }

func (apiServerInfo *ApiServerInfo) SetApiVersion(apiVersion string) {
	apiServerInfo.ApiVersion = apiVersion
}

func (apiServerInfo *ApiServerInfo) SetObject(object map[string]interface{}) {
	apiServerInfo.Data = object
}

func (ApiServerInfo *ApiServerInfo) SetProvider(provider string) {
	ApiServerInfo.Metadata.SetProvider(provider)
}

func (apiServerInfo *ApiServerInfo) SetApiServerVersion(version *version.Info) {
	apiServerInfo.Data = version
	apiServerInfo.SetName(apiServerInfoVersionName)
	// TODO: Parse the gitVersion and set the provider
}

func NewApiServerInfo() *ApiServerInfo {
	return &ApiServerInfo{
		ApiVersion: ApiServerInfoApiVersion,
		Kind:       string(TypeApiServerInfo),
		Metadata:   ApiServerMetadata{},
	}
}

func NewApiServerVersionInfo(version *version.Info) *ApiServerInfo {
	apiServerInfo := NewApiServerInfo()
	apiServerInfo.SetApiServerVersion(version)
	return apiServerInfo
}
