package initcontainerinstance

import (
	"fmt"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/workloadinterface"

	core1 "k8s.io/api/core/v1"
)

// GenerateInstanceID generates instance ID from workload
func GenerateInstanceID(w workloadinterface.IWorkload) ([]InstanceID, error) {
	if w.GetKind() != "Pod" {
		return nil, fmt.Errorf("CreateInstanceID: workload kind must be Pod for create instance ID")
	}

	ownerReferences, err := w.GetOwnerReferences()
	if err != nil {
		return nil, err
	}

	initContainers, err := w.GetInitContainers()
	if err != nil {
		return nil, err
	}

	return listInstanceIDs(ownerReferences, initContainers, w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName())
}

// GenerateInstanceIDFromPod generates instance ID from pod
func GenerateInstanceIDFromPod(pod *core1.Pod) ([]InstanceID, error) {
	return listInstanceIDs(pod.GetOwnerReferences(), pod.Spec.InitContainers, pod.APIVersion, pod.GetNamespace(), pod.Kind, pod.GetName())
}

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
func GenerateInstanceIDFromString(input string) (*InstanceID, error) {

	instanceID := &InstanceID{}

	// Split the input string by the field separator "/"
	fields := strings.Split(input, helpers.StringFormatSeparator)
	if len(fields) != 5 && len(fields) != 6 {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	i := 0
	instanceID.apiVersion = strings.TrimPrefix(fields[0], helpers.PrefixApiVersion)

	// if the apiVersion has a group, e.g. apps/v1
	if len(fields) == 6 {
		instanceID.apiVersion += helpers.StringFormatSeparator + fields[1]
		i += 1
	}

	instanceID.namespace = strings.TrimPrefix(fields[1+i], helpers.PrefixNamespace)
	instanceID.kind = strings.TrimPrefix(fields[2+i], helpers.PrefixKind)
	instanceID.name = strings.TrimPrefix(fields[3+i], helpers.PrefixName)
	instanceID.initContainerName = strings.TrimPrefix(fields[4+i], prefixInitContainer)

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	// Check if the input string is valid
	if instanceID.GetStringFormatted() != input {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	return instanceID, nil
}
