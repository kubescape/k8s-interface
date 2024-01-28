package containerinstance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/names"
)

// string format: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
const (
	prefixContainer = "containerName-"
	stringFormat    = helpers.PrefixApiVersion + "%s" + helpers.StringFormatSeparator + helpers.PrefixNamespace + "%s" + helpers.StringFormatSeparator + helpers.PrefixKind + "%s" + helpers.StringFormatSeparator + helpers.PrefixName + "%s" + helpers.StringFormatSeparator + prefixContainer + "%s"
)

const InstanceType helpers.InstanceType = "container"

// ensure that InstanceID implements IInstanceID
var _ instanceidhandler.IInstanceID = &InstanceID{}

type InstanceID struct {
	apiVersion    string
	namespace     string
	kind          string
	name          string
	containerName string
}

func (id *InstanceID) GetInstanceType() helpers.InstanceType {
	return InstanceType
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
		helpers.ApiGroupMetadataKey:      group,
		helpers.ApiVersionMetadataKey:    version,
		helpers.NamespaceMetadataKey:     id.GetNamespace(),
		helpers.KindMetadataKey:          id.GetKind(),
		helpers.NameMetadataKey:          id.GetName(),
		helpers.ContainerNameMetadataKey: id.GetContainerName(),
	}
}

func (id *InstanceID) GetSlug() (string, error) {
	return names.InstanceIDToSlug(id.GetName(), id.GetKind(), id.GetContainerName(), id.GetHashed())
}
