package instanceidhandler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/k8sinterface"
)

// metadata keys
const (
	metadataPrefix           = "kubescape.io"
	ApiGroupMetadataKey      = metadataPrefix + "/workload-api-group"
	ApiVersionMetadataKey    = metadataPrefix + "/workload-api-version"
	ContainerNameMetadataKey = metadataPrefix + "/workload-container-name"
	ImageNameMetadataKey     = metadataPrefix + "/image-name"
	ImageTagMetadataKey      = metadataPrefix + "/image-tag"
	InstanceIDMetadataKey    = metadataPrefix + "/instance-id"
	KindMetadataKey          = metadataPrefix + "/workload-kind"
	NameMetadataKey          = metadataPrefix + "/workload-name"
	NamespaceMetadataKey     = metadataPrefix + "/workload-namespace"
	StatusMetadataKey        = metadataPrefix + "/status"
	WlidMetadataKey          = metadataPrefix + "/wlid"
	ContextMetadataKey       = metadataPrefix + "/context"
)

// metadata values
const (
	ContextMetadataKeyFiltered    = "filtered"
	ContextMetadataKeyNonFiltered = "non-filtered"
)

// string format: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
const (
	stringFormatSeparator = "/"
	prefixApiVersion      = "apiVersion-"
	prefixNamespace       = "namespace-"
	prefixKind            = "kind-"
	prefixName            = "name-"
	prefixContainer       = "containerName-"
	stringFormat          = prefixApiVersion + "%s" + stringFormatSeparator + prefixNamespace + "%s" + stringFormatSeparator + prefixKind + "%s" + stringFormatSeparator + prefixName + "%s" + stringFormatSeparator + prefixContainer + "%s"
)

// ensure that InstanceID implements IInstanceID
var _ instanceidhandler.IInstanceID = &InstanceID{}

type InstanceID struct {
	apiVersion    string
	namespace     string
	kind          string
	name          string
	containerName string
}

func (id *InstanceID) GetAPIVersion() string {
	return id.apiVersion
}

func (id *InstanceID) GetNamespace() string {
	return id.namespace
}

func (id *InstanceID) GetKind() string {
	return id.kind
}

func (id *InstanceID) GetName() string {
	return id.name
}

func (id *InstanceID) GetContainerName() string {
	return id.containerName
}

func (id *InstanceID) SetAPIVersion(apiVersion string) {
	id.apiVersion = apiVersion
}

func (id *InstanceID) SetNamespace(namespace string) {
	id.namespace = namespace
}

func (id *InstanceID) SetKind(kind string) {
	id.kind = kind
}

func (id *InstanceID) SetName(name string) {
	id.name = name
}

func (id *InstanceID) SetContainerName(containerName string) {
	id.containerName = containerName
}

func (id *InstanceID) GetStringFormatted() string {
	return fmt.Sprintf(stringFormat, id.GetAPIVersion(), id.GetNamespace(), id.GetKind(), id.GetName(), id.GetContainerName())
}

func (id *InstanceID) GetHashed() string {
	hash := sha256.Sum256([]byte(id.GetStringFormatted()))
	str := hex.EncodeToString(hash[:])
	return str
}

func (id *InstanceID) GetLabels() map[string]string {
	group, version := k8sinterface.SplitApiVersion(id.GetAPIVersion())
	return map[string]string{
		ApiGroupMetadataKey:      group,
		ApiVersionMetadataKey:    version,
		NamespaceMetadataKey:     id.GetNamespace(),
		KindMetadataKey:          id.GetKind(),
		NameMetadataKey:          id.GetName(),
		ContainerNameMetadataKey: id.GetContainerName(),
	}
}
