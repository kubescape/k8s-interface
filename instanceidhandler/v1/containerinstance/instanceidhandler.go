package containerinstance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/names"
)

var (
	anyGroup          = "(.+)"
	anyNoSlash        = "([^/]+)"
	RegexFormatted    = regexp.MustCompile(strings.Join([]string{helpers.PrefixApiVersion + anyGroup, helpers.PrefixNamespace + anyNoSlash, helpers.PrefixKind + anyNoSlash, helpers.PrefixName + anyNoSlash, anyNoSlash + "Name-" + anyNoSlash}, helpers.StringFormatSeparator))
	StringFormatted   = strings.Join([]string{helpers.PrefixApiVersion + "%s", helpers.PrefixNamespace + "%s", helpers.PrefixKind + "%s", helpers.PrefixName + "%s", "%sName-%s"}, helpers.StringFormatSeparator)
	RegexNoContainer  = regexp.MustCompile(strings.Join([]string{helpers.PrefixApiVersion + anyGroup, helpers.PrefixNamespace + anyNoSlash, helpers.PrefixKind + anyNoSlash, helpers.PrefixName + anyNoSlash}, helpers.StringFormatSeparator))
	StringNoContainer = strings.Join([]string{helpers.PrefixApiVersion + "%s", helpers.PrefixNamespace + "%s", helpers.PrefixKind + "%s", helpers.PrefixName + "%s"}, helpers.StringFormatSeparator)
)

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
	TemplateHash  string
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
	if id.AlternateName != "" {
		return fmt.Sprintf(StringFormatted, id.ApiVersion, id.Namespace, id.Kind, id.AlternateName, id.InstanceType, id.ContainerName)
	}
	return fmt.Sprintf(StringFormatted, id.ApiVersion, id.Namespace, id.Kind, id.Name, id.InstanceType, id.ContainerName)
}

func (id *InstanceID) GetStringNoContainer() string {
	if id.AlternateName != "" {
		return fmt.Sprintf(StringNoContainer, id.ApiVersion, id.Namespace, id.Kind, id.AlternateName)
	}
	return fmt.Sprintf(StringNoContainer, id.ApiVersion, id.Namespace, id.Kind, id.Name)
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
		helpers.TemplateHashKey:          id.TemplateHash,
	}
}

func (id *InstanceID) GetOneTimeSlug(noContainer bool) (string, error) {
	u := uuid.New()
	hexSuffix := hex.EncodeToString(u[:]) // u is [16]byte, u[:] is []byte
	suffix := "-" + hexSuffix
	return id.getSlugWithSuffix(noContainer, suffix)
}

func (id *InstanceID) GetSlug(noContainer bool) (string, error) {
	return id.getSlugWithSuffix(noContainer, "")
}

func (id *InstanceID) getSlugWithSuffix(noContainer bool, suffix string) (string, error) {
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
	slug, err := names.InstanceIDToSlugWithMaxLength(name, kind, containerName, hashedID, names.MaxDNSSubdomainLength-len(suffix))
	if err != nil {
		return "", err
	}
	return slug + suffix, nil
}

func (id *InstanceID) GetTemplateHash() string {
	return id.TemplateHash
}
