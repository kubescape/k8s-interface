package instanceidhandler

import (
	"fmt"

	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateInstanceID(instanceID *InstanceID) error {
	if instanceID.apiVersion == "" {
		return fmt.Errorf("invalid instanceID: apiVersion cannot be empty")
	}
	if instanceID.namespace == "" {
		return fmt.Errorf("invalid instanceID: namespace cannot be empty")
	}
	if instanceID.kind == "" {
		return fmt.Errorf("invalid instanceID: kind cannot be empty")
	}
	if instanceID.name == "" {
		return fmt.Errorf("invalid instanceID: name cannot be empty")
	}
	if instanceID.containerName == "" {
		return fmt.Errorf("invalid instanceID: containerName cannot be empty")
	}
	return nil
}

func listInstanceIDs(ownerReferences []metav1.OwnerReference, containers []core1.Container, apiVersion, namespace, kind, name string) ([]*InstanceID, error) {

	instanceIDs := make([]*InstanceID, 0)
	parentKind, parentName := "", ""

	if len(ownerReferences) == 0 || ownerReferences[0].Kind == "Node" {
		parentKind = kind
		parentName = name
	} else {
		parentKind = ownerReferences[0].Kind
		parentName = ownerReferences[0].Name
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("failed to validate instance ID: missing containers")
	}

	for i := range containers {
		instanceID := &InstanceID{
			apiVersion:    apiVersion,
			namespace:     namespace,
			kind:          parentKind,
			name:          parentName,
			containerName: containers[i].Name,
		}

		if err := validateInstanceID(instanceID); err != nil {
			return nil, fmt.Errorf("failed to validate instance ID: %w", err)
		}

		instanceIDs = append(instanceIDs, instanceID)
	}

	return instanceIDs, nil
}
