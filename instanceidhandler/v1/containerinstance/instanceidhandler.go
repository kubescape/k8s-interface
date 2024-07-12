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

// ensure that InstanceID implements IInstanceID
var _ instanceidhandler.IInstanceID = (*InstanceID)(nil)

type InstanceID struct {
	ApiVersion    string
	Namespace     string
	Kind          string
	Name          string
	AlternateName string
	ContainerName string
	InstanceType  string
}

func (id *InstanceID) GetInstanceType() helpers.InstanceType {
	return helpers.InstanceType(id.InstanceType)
}

func (id *InstanceID) GetName() string {
	return id.Name
}

func (id *InstanceID) GetContainerName() string {
	return id.ContainerName
}

func (id *InstanceID) GetStringFormatted() string {
	stringFormat := helpers.PrefixApiVersion + "%s" + helpers.StringFormatSeparator + helpers.PrefixNamespace + "%s" + helpers.StringFormatSeparator + helpers.PrefixKind + "%s" + helpers.StringFormatSeparator + helpers.PrefixName + "%s" + helpers.StringFormatSeparator + "%sName-%s"
	if id.AlternateName != "" {
		return fmt.Sprintf(stringFormat, id.ApiVersion, id.Namespace, id.Kind, id.AlternateName, id.InstanceType, id.ContainerName)
	}
	return fmt.Sprintf(stringFormat, id.ApiVersion, id.Namespace, id.Kind, id.Name, id.InstanceType, id.ContainerName)
}

func (id *InstanceID) GetHashed() string {
	hash := sha256.Sum256([]byte(id.GetStringFormatted()))
	str := hex.EncodeToString(hash[:])
	return str
}

func (id *InstanceID) GetLabels() map[string]string {
	group, version := k8sinterface.SplitApiVersion(id.ApiVersion)
	return map[string]string{
		helpers.ApiGroupMetadataKey:      group,
		helpers.ApiVersionMetadataKey:    version,
		helpers.NamespaceMetadataKey:     id.Namespace,
		helpers.KindMetadataKey:          id.Kind,
		helpers.NameMetadataKey:          id.Name,
		helpers.ContainerNameMetadataKey: id.ContainerName,
	}
}

func (id *InstanceID) GetSlug(noContainer bool) (string, error) {
	name := id.Name
	kind := id.Kind
	containerName := id.ContainerName
	hashedID := id.GetHashed()
	// use alternate name if present
	if len(id.AlternateName) > 0 {
		name = id.AlternateName
	}
	// eventually remove the container name
	if noContainer {
		containerName = ""
	}
	return names.InstanceIDToSlug(name, kind, containerName, hashedID)
}
