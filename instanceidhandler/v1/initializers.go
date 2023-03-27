package instanceidhandler

import (
	"fmt"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/workloadinterface"

	core1 "k8s.io/api/core/v1"
)

// GenerateInstanceID generates instance ID from workload
func GenerateInstanceID(w workloadinterface.IWorkload) ([]instanceidhandler.IInstanceID, error) {
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
func GenerateInstanceIDFromPod(pod *core1.Pod) ([]instanceidhandler.IInstanceID, error) {
	return listInstanceIDs(pod.GetOwnerReferences(), pod.Spec.Containers, pod.APIVersion, pod.GetNamespace(), pod.Kind, pod.GetName())
}

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
func GenerateInstanceIDFromString(input string) (instanceidhandler.IInstanceID, error) {

	instanceID := &InstanceID{}

	// Split the input string by the field separator "/"
	fields := strings.Split(input, stringFormatSeparator)
	if len(fields) != 5 && len(fields) != 6 {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	i := 0
	instanceID.apiVersion = strings.TrimPrefix(fields[0], prefixApiVersion)

	// if the apiVersion has a group, e.g. apps/v1
	if len(fields) == 6 {
		instanceID.apiVersion += stringFormatSeparator + fields[1]
		i += 1
	}

	instanceID.namespace = strings.TrimPrefix(fields[1+i], prefixNamespace)
	instanceID.kind = strings.TrimPrefix(fields[2+i], prefixKind)
	instanceID.name = strings.TrimPrefix(fields[3+i], prefixName)
	instanceID.containerName = strings.TrimPrefix(fields[4+i], prefixContainer)

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	// Check if the input string is valid
	if instanceID.GetStringFormatted() != input {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	return instanceID, nil
}
