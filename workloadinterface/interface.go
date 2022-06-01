package workloadinterface

import (
	"github.com/armosec/armoapi-go/apis"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ObjectType string

type IMetadata interface {
	// Set
	SetNamespace(string)
	SetName(string)
	SetKind(string)
	SetWorkload(map[string]interface{}) // DEPRECATED
	SetObject(map[string]interface{})
	// TODO - AetApiVersion

	// Get
	GetNamespace() string
	GetName() string
	GetKind() string
	GetApiVersion() string
	GetWorkload() map[string]interface{} // DEPRECATED
	GetObject() map[string]interface{}
	GetID() string // Get object unique ID

	GetObjectType() ObjectType // Get struct type

}
type IBasicWorkload interface {
	IMetadata

	// Set
	SetLabel(key, value string)
	SetAnnotation(key, value string)

	// Get

	GetVersion() string
	GetGroup() string
	GetGenerateName() string
	GetInnerAnnotation(string) (string, bool)
	GetPodAnnotation(string) (string, bool)
	GetAnnotation(string) (string, bool)
	GetLabel(string) (string, bool)
	GetAnnotations() map[string]string
	GetInnerAnnotations() map[string]string
	GetPodAnnotations() map[string]string
	GetLabels() map[string]string
	GetInnerLabels() map[string]string
	GetPodLabels() map[string]string
	GetVolumes() ([]corev1.Volume, error)
	GetReplicas() int
	GetContainers() ([]corev1.Container, error)
	GetInitContainers() ([]corev1.Container, error)
	GetOwnerReferences() ([]metav1.OwnerReference, error)
	GetImagePullSecret() ([]corev1.LocalObjectReference, error)
	GetServiceAccountName() string
	GetSelector() (*metav1.LabelSelector, error)
	GetResourceVersion() string
	GetUID() string
	GetPodSpec() (*corev1.PodSpec, error)
	GetData() map[string]interface{}
	//GetSpiffe() string

	// REMOVE
	RemoveLabel(string)
	RemoveAnnotation(string)
	RemovePodStatus()
	RemoveResourceVersion()
}

type IWorkload interface {
	IBasicWorkload

	// Convert
	ToUnstructured() (*unstructured.Unstructured, error)
	ToString() string // Return workload in string representation
	Json() string     // DEPRECATED, use ToString

	// GET
	GetWlid() string // Get ARMO workload ID -> wlid://cluster-<cluster-name>/namespace-<namespace>/<kind>-<name>
	GetJobID() *apis.JobTracking
	GenerateWlid(string) string

	// SET
	SetWlid(string)
	SetInject()
	SetIgnore()
	SetUpdateTime()
	SetJobID(apis.JobTracking)
	SetCompatible()
	SetIncompatible()
	SetReplaceheaders()

	// EXIST
	IsIgnore() bool
	IsInject() bool
	IsAttached() bool
	IsCompatible() bool
	IsIncompatible() bool
	IsReplaceheaders() bool

	// REMOVE
	RemoveWlid()
	RemoveSecretData()
	RemoveInject()
	RemoveIgnore()
	RemoveUpdateTime()
	RemoveJobID()
	RemoveCompatible()
	RemoveArmoMetadata()
	RemoveArmoLabels()
	RemoveArmoAnnotations()
}
