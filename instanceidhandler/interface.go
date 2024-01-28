package instanceidhandler

import "github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"

type IInstanceID interface {
	// GetInstanceType returns the type of the instance ID
	GetInstanceType() helpers.InstanceType

	GetAPIVersion() string
	GetNamespace() string
	GetKind() string
	GetName() string
	GetContainerName() string
	SetAPIVersion(string)
	SetNamespace(string)
	SetKind(string)
	SetName(string)
	SetContainerName(string)
	GetStringFormatted() string
	GetHashed() string
	GetLabels() map[string]string
	// GetSlug returns a string with a human-friendly and Kubernetes-compatible name
	GetSlug() (string, error)
}
