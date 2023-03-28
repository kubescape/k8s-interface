package instanceidhandler

import (
	"fmt"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/k8sinterface"
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

	if len(containers) == 0 {
		return nil, fmt.Errorf("failed to validate instance ID: missing containers")
	}

	instanceIDs := make([]instanceidhandler.IInstanceID, 0)

	parentApiVersion, parentKind, parentName := apiVersion, kind, name

	if len(ownerReferences) != 0 && !ignoreOwnerReference(ownerReferences[0].Kind) {
		parentApiVersion = ownerReferences[0].APIVersion
		parentKind = ownerReferences[0].Kind
		parentName = ownerReferences[0].Name
	}

	for i := range containers {
		instanceID := &InstanceID{
			apiVersion:    parentApiVersion,
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

// ignoreOwnerReference returns true if the owner reference is a node or a unknown resource (CRD)
func ignoreOwnerReference(ownerKind string) bool {
	if ownerKind == "Node" {
		return true
	}
	if _, e := k8sinterface.GetGroupVersionResource(ownerKind); e != nil {
		return true
	}
	return false
}
