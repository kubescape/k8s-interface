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
// ============================== CloudProviderListEntitiesForPolicies ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderListEntitiesForPolicies) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderListEntitiesForPolicies) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderListEntitiesForPolicies) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderListEntitiesForPolicies) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderListEntitiesForPolicies) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderListEntitiesForPolicies) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderListEntitiesForPolicies) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderListEntitiesForPolicies) SetObject(object map[string]interface{}) {
	if !apis.IsTypeListEntitiesForPolicies(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderListEntitiesForPolicies{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderListEntitiesForPolicies) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderListEntitiesForPolicies) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderListEntitiesForPolicies
}
func (description *CloudProviderListEntitiesForPolicies) GetKind() string {
	return description.Kind
}

func (description *CloudProviderListEntitiesForPolicies) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderListEntitiesForPolicies) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderListEntitiesForPolicies) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderListEntitiesForPolicies) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderListEntitiesForPolicies) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderListEntitiesForPolicies) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderListEntitiesForPolicies) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}

// ==========================================================================================================
// ============================== CloudProviderPolicyVersion ==================================================
// ==========================================================================================================
// Setters
func (description *CloudProviderPolicyVersion) SetNamespace(namespace string) {
	description.SetProvider(namespace)
}

func (description *CloudProviderPolicyVersion) SetApiVersion(apiVersion string) {
	description.ApiVersion = apiVersion
}

func (description *CloudProviderPolicyVersion) SetName(name string) {
	description.Metadata.SetName(name)
}

func (description *CloudProviderPolicyVersion) SetProvider(provider string) {
	description.Metadata.SetProvider(provider)
}

func (description *CloudProviderPolicyVersion) SetKind(kind string) {
	description.Kind = kind
}

func (description *CloudProviderPolicyVersion) SetData(data map[string]interface{}) {
	description.Data = data
}

func (description *CloudProviderPolicyVersion) SetWorkload(object map[string]interface{}) {
	description.SetObject(object)
}

func (description *CloudProviderPolicyVersion) SetObject(object map[string]interface{}) {
	if !apis.IsTypePolicyVersion(object) {
		return
	}
	if b := workloadinterface.MapToBytes(object); len(b) > 0 {
		d := &CloudProviderPolicyVersion{}
		if err := json.Unmarshal(b, d); err == nil {
			description.SetApiVersion(d.GetApiVersion())
			description.SetKind(d.GetKind())
			description.SetData(d.GetData())
			description.Metadata = d.Metadata
		}
	}
}

// Getters

func (description *CloudProviderPolicyVersion) GetApiVersion() string {
	return description.ApiVersion
}

func (description *CloudProviderPolicyVersion) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderPolicyVersion
}
func (description *CloudProviderPolicyVersion) GetKind() string {
	return description.Kind
}

func (description *CloudProviderPolicyVersion) GetName() string {
	return description.Metadata.GetName()
}

// provider -> eks/gke/etc.
func (description *CloudProviderPolicyVersion) GetProvider() string {
	return description.Metadata.GetProvider()
}

// Compatible with the IMetadata interface
func (description *CloudProviderPolicyVersion) GetNamespace() string {
	return description.GetProvider()
}

func (description *CloudProviderPolicyVersion) GetWorkload() map[string]interface{} {
	return description.GetObject()
}

func (description *CloudProviderPolicyVersion) GetData() map[string]interface{} {
	return description.Data
}

func (description *CloudProviderPolicyVersion) GetObject() map[string]interface{} {
	m := map[string]interface{}{}
	b, err := json.Marshal(*description)
	if err != nil {
		return m
	}
	return workloadinterface.BytesToMap(b)
}

// ApiVersion/Kind/Name
func (description *CloudProviderPolicyVersion) GetID() string {
	return fmt.Sprintf("%s/%s/%s", k8sinterface.JoinGroupVersion(k8sinterface.SplitApiVersion(description.GetApiVersion())), description.GetKind(), description.GetName())
}
