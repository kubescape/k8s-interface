package instanceIDhandler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/kubescape/k8s-interface/workloadinterface"
)

type InstanceID struct {
	apiVersion    string
	namespace     string
	kind          string
	name          string
	containerName string
}

const (
	labelPrefix                 = "kubescape.io"
	StringFormat                = "apiVersion-%s/namespace-%s/kind-%s/name-%s/containerName-%s"
	LabelFormatKeyApiGroup      = labelPrefix + "/workload-api-group"
	LabelFormatKeyApiVersion    = labelPrefix + "/workload-api-version"
	LabelFormatKeyNamespace     = labelPrefix + "/workload-namespace"
	LabelFormatKeyKind          = labelPrefix + "/workload-kind"
	LabelFormatKeyName          = labelPrefix + "/workload-name"
	LabelFormatKeyContainerName = labelPrefix + "/workload-container-name"
)

func GenerateInstanceID(w *workloadinterface.Workload) ([]InstanceID, error) {
	instanceIDs := make([]InstanceID, 0)
	parentKind, parentName := "", ""

	if w.GetKind() != "Pod" {
		return nil, fmt.Errorf("CreateInstanceID: workload kind must be Pod for create instance ID")
	}
	ownerReferences, err := w.GetOwnerReferences()
	if err != nil {
		return nil, err
	}
	if len(ownerReferences) == 0 {
		parentKind = w.GetKind()
		parentName = w.GetName()
	} else {
		parentKind = ownerReferences[0].Kind
		parentName = ownerReferences[0].Name
		if parentKind == "Node" {
			parentKind = w.GetKind()
			parentName = w.GetName()
		}
	}

	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	for i := range containers {
		instanceID := InstanceID{
			apiVersion:    w.GetApiVersion(),
			namespace:     w.GetNamespace(),
			kind:          parentKind,
			name:          parentName,
			containerName: containers[i].Name,
		}
		instanceIDs = append(instanceIDs, instanceID)
	}

	return instanceIDs, nil
}

func (id *InstanceID) GetAPIVersion() string {
	if id != nil {
		return id.apiVersion
	}
	return ""
}

func (id *InstanceID) GetNamespace() string {
	if id != nil {
		return id.namespace
	}
	return ""
}

func (id *InstanceID) GetKind() string {
	if id != nil {
		return id.kind
	}
	return ""
}

func (id *InstanceID) GetName() string {
	if id != nil {
		return id.name
	}
	return ""
}

func (id *InstanceID) GetContainerName() string {
	if id != nil {
		return id.containerName
	}
	return ""
}

func (id *InstanceID) SetAPIVersion(apiVersion string) {
	if id != nil {
		id.apiVersion = apiVersion
	}
}

func (id *InstanceID) SetNamespace(namespace string) {
	if id != nil {
		id.namespace = namespace
	}
}

func (id *InstanceID) SetKind(kind string) {
	if id != nil {
		id.kind = kind
	}
}

func (id *InstanceID) SetName(name string) {
	if id != nil {
		id.name = name
	}
}

func (id *InstanceID) SetContainerName(containerName string) {
	if id != nil {
		id.containerName = containerName
	}
}

func (id *InstanceID) GetStringFormatted() string {
	return fmt.Sprintf(StringFormat, id.GetAPIVersion(), id.GetNamespace(), id.GetKind(), id.GetName(), id.GetContainerName())
}

func (id *InstanceID) GetIDHashed() string {
	hash := sha256.Sum256([]byte(id.GetStringFormatted()))
	str := hex.EncodeToString(hash[:])
	return str
}

func (id *InstanceID) GetLabels() map[string]string {
	if id != nil {
		group, version := k8sinterface.SplitApiVersion(id.GetAPIVersion())
		return map[string]string{
			LabelFormatKeyApiGroup:      group,
			LabelFormatKeyApiVersion:    version,
			LabelFormatKeyNamespace:     id.GetNamespace(),
			LabelFormatKeyKind:          id.GetKind(),
			LabelFormatKeyName:          id.GetName(),
			LabelFormatKeyContainerName: id.GetContainerName(),
		}
	}
	return map[string]string{}
}
