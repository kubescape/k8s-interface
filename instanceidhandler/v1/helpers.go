package instanceidhandler

import (
	"fmt"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateInstanceID(instanceID instanceidhandler.IInstanceID) error {
	if instanceID.GetAPIVersion() == "" {
		return fmt.Errorf("invalid instanceID: apiVersion cannot be empty")
	}
	if instanceID.GetNamespace() == "" {
		return fmt.Errorf("invalid instanceID: namespace cannot be empty")
	}
	if instanceID.GetKind() == "" {
		return fmt.Errorf("invalid instanceID: kind cannot be empty")
	}
	if instanceID.GetName() == "" {
		return fmt.Errorf("invalid instanceID: name cannot be empty")
	}
	if instanceID.GetContainerName() == "" {
		return fmt.Errorf("invalid instanceID: containerName cannot be empty")
	}
	return nil
}

func listInstanceIDs(ownerReferences []metav1.OwnerReference, containers []core1.Container, apiVersion, namespace, kind, name string) ([]instanceidhandler.IInstanceID, error) {

	instanceIDs := make([]instanceidhandler.IInstanceID, 0)
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
