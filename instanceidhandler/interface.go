package instanceidhandler

type IInstanceID interface {
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
