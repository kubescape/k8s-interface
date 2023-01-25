package v1

import (
	"encoding/json"
	"fmt"

	"github.com/kubescape/k8s-interface/cloudsupport/apis"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
)

// ==========================================================================================================
// ================================= CloudProviderMetadata ==================================================
// ==========================================================================================================

// getters
func (cloudProviderMetadata *CloudProviderMetadata) GetName() string {
	return cloudProviderMetadata.Name
}

func (cloudProviderMetadata *CloudProviderMetadata) GetProvider() string {
	return cloudProviderMetadata.Provider
}

// setters
func (cloudProviderMetadata *CloudProviderMetadata) SetName(name string) {
	cloudProviderMetadata.Name = name
}

func (cloudProviderMetadata *CloudProviderMetadata) SetProvider(provider string) {
	cloudProviderMetadata.Provider = provider
}

// ==========================================================================================================
// ============================== CloudProviderDescribe ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderDescribe) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderDescribe) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderDescribe) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderDescribe) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderDescribe) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderDescribe) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderDescribe) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderDescribe) SetObject(object map[string]interface{}) {
	if !apis.IsTypeDescriptiveInfoFromCloudProvider(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderDescribe{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderDescribe) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderDescribe) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderDescribe
}
func (description *CloudProviderDescribe) GetKind() string {
	return description.Kind
}

func (description *CloudProviderDescribe) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderDescribe) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderDescribe) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderDescribe) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderDescribe) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderDescribe) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderDescribe) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}

// ==========================================================================================================
// ============================== CloudProviderDescribeRepositories ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderDescribeRepositories) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderDescribeRepositories) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderDescribeRepositories) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderDescribeRepositories) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderDescribeRepositories) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderDescribeRepositories) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderDescribeRepositories) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderDescribeRepositories) SetObject(object map[string]interface{}) {
	if !apis.IsTypeDescribeRepositories(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderDescribeRepositories{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderDescribeRepositories) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderDescribeRepositories) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderDescribeRepositories
}
func (description *CloudProviderDescribeRepositories) GetKind() string {
	return description.Kind
}

func (description *CloudProviderDescribeRepositories) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderDescribeRepositories) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderDescribeRepositories) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderDescribeRepositories) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderDescribeRepositories) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderDescribeRepositories) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderDescribeRepositories) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}

// ==========================================================================================================
// ============================== CloudProviderListRolePolicies ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderListRolePolicies) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderListRolePolicies) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderListRolePolicies) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderListRolePolicies) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderListRolePolicies) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderListRolePolicies) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderListRolePolicies) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderListRolePolicies) SetObject(object map[string]interface{}) {
	if !apis.IsTypeDescribeRepositories(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderListRolePolicies{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderListRolePolicies) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderListRolePolicies) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderListRolePolicies
}
func (description *CloudProviderListRolePolicies) GetKind() string {
	return description.Kind
}

func (description *CloudProviderListRolePolicies) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderListRolePolicies) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderListRolePolicies) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderListRolePolicies) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderListRolePolicies) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderListRolePolicies) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderListRolePolicies) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}

// ==========================================================================================================
// ============================== CloudProviderListUserPolicies ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderListUserPolicies) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderListUserPolicies) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderListUserPolicies) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderListUserPolicies) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderListUserPolicies) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderListUserPolicies) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderListUserPolicies) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderListUserPolicies) SetObject(object map[string]interface{}) {
	if !apis.IsTypeDescribeRepositories(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderListUserPolicies{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderListUserPolicies) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderListUserPolicies) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderListUserPolicies
}
func (description *CloudProviderListUserPolicies) GetKind() string {
	return description.Kind
}

func (description *CloudProviderListUserPolicies) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderListUserPolicies) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderListUserPolicies) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderListUserPolicies) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderListUserPolicies) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderListUserPolicies) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderListUserPolicies) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}

// ==========================================================================================================
// ============================== CloudProviderListGroupPolicies ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderListGroupPolicies) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderListGroupPolicies) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderListGroupPolicies) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderListGroupPolicies) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderListGroupPolicies) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderListGroupPolicies) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderListGroupPolicies) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderListGroupPolicies) SetObject(object map[string]interface{}) {
	if !apis.IsTypeDescribeRepositories(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderListGroupPolicies{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderListGroupPolicies) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderListGroupPolicies) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderListGroupPolicies
}
func (description *CloudProviderListGroupPolicies) GetKind() string {
	return description.Kind
}

func (description *CloudProviderListGroupPolicies) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderListGroupPolicies) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderListGroupPolicies) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderListGroupPolicies) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderListGroupPolicies) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderListGroupPolicies) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderListGroupPolicies) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}
