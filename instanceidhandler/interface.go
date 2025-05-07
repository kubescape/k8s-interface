package instanceidhandler

import "github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"

type IInstanceID interface {
	GetInstanceType() helpers.InstanceType
	GetName() string
	GetContainerName() string
	GetStringFormatted() string
	GetStringNoContainer() string
	GetHashed() string
	GetLabels() map[string]string
	// GetSlug returns a string with a human-friendly and Kubernetes-compatible name
	GetSlug(noContainer bool) (string, error)
	GetTemplateHash() string
}
