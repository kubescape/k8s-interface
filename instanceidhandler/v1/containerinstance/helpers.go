package containerinstance

import (
	"fmt"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateInstanceID(instanceID *InstanceID) error {
	if instanceID.ApiVersion == "" {
		return fmt.Errorf("invalid instanceID: apiVersion cannot be empty")
	}
	if instanceID.Namespace == "" {
		return fmt.Errorf("invalid instanceID: namespace cannot be empty")
	}
	if instanceID.Kind == "" {
		return fmt.Errorf("invalid instanceID: kind cannot be empty")
	}
	if instanceID.Name == "" {
		return fmt.Errorf("invalid instanceID: name cannot be empty")
	}
	return nil
}

func ListInstanceIDs(ownerReference *metav1.OwnerReference, containers []core1.Container, instanceType, apiVersion, namespace, kind, name, alternateName, templateHash string) ([]InstanceID, error) {
	instanceIDs := make([]InstanceID, 0)

	if len(containers) == 0 {
		return instanceIDs, nil
	}

	parentApiVersion, parentKind, parentName := apiVersion, kind, name

	if ownerReference != nil && !helpers.IgnoreOwnerReference(ownerReference.Kind) {
		parentApiVersion = ownerReference.APIVersion
		parentKind = ownerReference.Kind
		parentName = ownerReference.Name
	}

	for i := range containers {
		instanceID := InstanceID{
			ApiVersion:    parentApiVersion,
			Namespace:     namespace,
			Kind:          parentKind,
			Name:          parentName,
			AlternateName: alternateName,
			ContainerName: containers[i].Name,
			InstanceType:  instanceType,
			TemplateHash:  templateHash,
		}

		if err := validateInstanceID(&instanceID); err != nil {
			return nil, fmt.Errorf("failed to validate instance ID: %w", err)
		}

		instanceIDs = append(instanceIDs, instanceID)
	}

	return instanceIDs, nil
}
