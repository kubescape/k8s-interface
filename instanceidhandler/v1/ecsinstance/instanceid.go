package ecsinstance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/names"
)

var (
	anyGroup             = "(.+)"
	anyNoSlash           = "([^/]+)"
	RegexFormattedECS    = regexp.MustCompile(strings.Join([]string{helpers.PrefixApiVersion + anyGroup, helpers.PrefixNamespace + anyNoSlash, helpers.PrefixKind + anyNoSlash, helpers.PrefixName + anyNoSlash, anyNoSlash + "Name-" + anyNoSlash}, helpers.StringFormatSeparator))
	StringFormattedECS   = strings.Join([]string{helpers.PrefixApiVersion + "%s", helpers.PrefixNamespace + "%s", helpers.PrefixKind + "%s", helpers.PrefixName + "%s", "%sName-%s"}, helpers.StringFormatSeparator)
	RegexNoContainerECS  = regexp.MustCompile(strings.Join([]string{helpers.PrefixApiVersion + anyGroup, helpers.PrefixNamespace + anyNoSlash, helpers.PrefixKind + anyNoSlash, helpers.PrefixName + anyNoSlash}, helpers.StringFormatSeparator))
	StringNoContainerECS = strings.Join([]string{helpers.PrefixApiVersion + "%s", helpers.PrefixNamespace + "%s", helpers.PrefixKind + "%s", helpers.PrefixName + "%s"}, helpers.StringFormatSeparator)
	ECSDefaultApiVersion = "ecs/v1"
	ECSDefaultAPIDate    = "2014-11-13"
)

const (
	TaskDefinitionKind = "TaskDefinition"
	ServiceKind        = "Service"
	TaskKind           = "Task"
)

var _ instanceidhandler.IInstanceID = (*ECSInstanceID)(nil)

type ECSInstanceID struct {
	ApiVersion             string
	ClusterName            string
	Kind                   string
	Name                   string
	AlternateName          string
	ContainerName          string
	InstanceType           string
	TemplateHash           string
	Region                 string
	AccountID              string
	TaskDefinition         string
	TaskDefinitionRevision string
}

func (id *ECSInstanceID) GetInstanceType() helpers.InstanceType {
	return helpers.InstanceType(id.InstanceType)
}

func (id *ECSInstanceID) GetName() string {
	return id.Name
}

func (id *ECSInstanceID) GetContainerName() string {
	return id.ContainerName
}

func (id *ECSInstanceID) GetStringFormatted() string {
	if id.AlternateName != "" {
		return fmt.Sprintf(StringFormattedECS, id.ApiVersion, id.ClusterName, id.Kind, id.AlternateName, id.InstanceType, id.ContainerName)
	}
	return fmt.Sprintf(StringFormattedECS, id.ApiVersion, id.ClusterName, id.Kind, id.Name, id.InstanceType, id.ContainerName)
}

func (id *ECSInstanceID) GetStringNoContainer() string {
	if id.AlternateName != "" {
		return fmt.Sprintf(StringNoContainerECS, id.ApiVersion, id.ClusterName, id.Kind, id.AlternateName)
	}
	return fmt.Sprintf(StringNoContainerECS, id.ApiVersion, id.ClusterName, id.Kind, id.Name)
}

func (id *ECSInstanceID) GetHashed() string {
	hash := sha256.Sum256([]byte(id.GetStringFormatted()))
	str := hex.EncodeToString(hash[:])
	return str
}

func (id *ECSInstanceID) GetLabels() map[string]string {
	labels := map[string]string{
		helpers.ApiGroupMetadataKey:      "ecs",
		helpers.ApiVersionMetadataKey:    ECSDefaultAPIDate,
		helpers.NamespaceMetadataKey:     id.ClusterName,
		helpers.KindMetadataKey:          id.Kind,
		helpers.NameMetadataKey:          id.Name,
		helpers.ContainerNameMetadataKey: id.ContainerName,
		helpers.TemplateHashKey:          id.TemplateHash,
	}
	if id.Region != "" {
		labels["kubescape.io/aws-region"] = id.Region
	}
	if id.AccountID != "" {
		labels["kubescape.io/aws-account-id"] = id.AccountID
	}
	if id.TaskDefinition != "" {
		labels["kubescape.io/ecs-task-definition"] = id.TaskDefinition
	}
	if id.TaskDefinitionRevision != "" {
		labels["kubescape.io/ecs-task-definition-revision"] = id.TaskDefinitionRevision
	}
	return labels
}

func (id *ECSInstanceID) GetOneTimeSlug(noContainer bool) (string, error) {
	u := uuid.New()
	hexSuffix := hex.EncodeToString(u[:])
	suffix := "-" + hexSuffix
	return id.getSlugWithSuffix(noContainer, suffix)
}

func (id *ECSInstanceID) GetSlug(noContainer bool) (string, error) {
	return id.getSlugWithSuffix(noContainer, "")
}

func (id *ECSInstanceID) getSlugWithSuffix(noContainer bool, suffix string) (string, error) {
	name := id.Name
	kind := strings.ToLower(id.Kind)
	containerName := id.ContainerName
	hashedID := id.GetHashed()

	if len(id.AlternateName) > 0 {
		name = id.AlternateName
	}
	if noContainer {
		containerName = ""
	}

	slug, err := names.InstanceIDToSlugWithMaxLength(name, kind, containerName, hashedID, names.MaxDNSSubdomainLength-len(suffix))
	if err != nil {
		return "", err
	}
	return slug + suffix, nil
}

func (id *ECSInstanceID) GetTemplateHash() string {
	return id.TemplateHash
}

func NewTaskDefinitionInstanceID(clusterName, taskDefinition, revision, containerName string) *ECSInstanceID {
	templateHash := fmt.Sprintf("%s:%s", taskDefinition, revision)
	return &ECSInstanceID{
		ApiVersion:             ECSDefaultApiVersion,
		ClusterName:            clusterName,
		Kind:                   TaskDefinitionKind,
		Name:                   taskDefinition,
		ContainerName:          containerName,
		InstanceType:           "container",
		TemplateHash:           templateHash,
		TaskDefinition:         taskDefinition,
		TaskDefinitionRevision: revision,
	}
}

func NewServiceInstanceID(clusterName, serviceName, containerName, taskDefinition, revision string) *ECSInstanceID {
	templateHash := fmt.Sprintf("%s:%s", taskDefinition, revision)
	return &ECSInstanceID{
		ApiVersion:             ECSDefaultApiVersion,
		ClusterName:            clusterName,
		Kind:                   ServiceKind,
		Name:                   serviceName,
		ContainerName:          containerName,
		InstanceType:           "container",
		TemplateHash:           templateHash,
		TaskDefinition:         taskDefinition,
		TaskDefinitionRevision: revision,
	}
}

func NewTaskInstanceID(clusterName, taskArn, containerName, taskDefinition, revision string) *ECSInstanceID {
	templateHash := fmt.Sprintf("%s:%s", taskDefinition, revision)
	return &ECSInstanceID{
		ApiVersion:             ECSDefaultApiVersion,
		ClusterName:            clusterName,
		Kind:                   TaskKind,
		Name:                   taskArn,
		AlternateName:          extractTaskIDFromArn(taskArn),
		ContainerName:          containerName,
		InstanceType:           "container",
		TemplateHash:           templateHash,
		TaskDefinition:         taskDefinition,
		TaskDefinitionRevision: revision,
	}
}

func extractTaskIDFromArn(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx != -1 {
		return arn[idx+1:]
	}
	return arn
}

const container = "container"
