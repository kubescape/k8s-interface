package instanceidhandler

import (
	"fmt"

	"github.com/kubescape/k8s-interface/workloadinterface"
	core1 "k8s.io/api/core/v1"
)

// GenerateInstanceID generates instance ID from workload
func GenerateInstanceID(w workloadinterface.IWorkload) ([]*InstanceID, error) {
	if w.GetKind() != "Pod" {
		return nil, fmt.Errorf("CreateInstanceID: workload kind must be Pod for create instance ID")
	}

	ownerReferences, err := w.GetOwnerReferences()
	if err != nil {
		return nil, err
	}

	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	return listInstanceIDs(ownerReferences, containers, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName())
}

// GenerateInstanceIDFromPod generates instance ID from pod
func GenerateInstanceIDFromPod(pod *core1.Pod) ([]*InstanceID, error) {
	return listInstanceIDs(pod.GetOwnerReferences(), pod.Spec.Containers, pod.APIVersion, pod.GetNamespace(), pod.Kind, pod.GetName())
}

func GenerateInstanceIDFromString(input string) (*InstanceID, error) {
	var apiVersion, namespace, kind, name, containerName string
	input += "\n" // add new line to the end of the string so that Sscanf can read it

	_, err := fmt.Sscanf(input, StringFormat, &apiVersion, &namespace, &kind, &name, &containerName)
	if err != nil {
		return nil, err
	}

	instanceID := &InstanceID{
		apiVersion:    apiVersion,
		namespace:     namespace,
		kind:          kind,
		name:          name,
		containerName: containerName,
	}

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	return instanceID, nil
}
