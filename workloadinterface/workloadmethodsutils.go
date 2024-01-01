package workloadinterface

import (
	"fmt"
	"strings"
)

func PodSpec(kind string) []string {
	switch kind {
	case "Pod", "Namespace":
		return []string{"spec"}
	case "CronJob":
		return []string{"spec", "jobTemplate", "spec", "template", "spec"}
	default:
		return []string{"spec", "template", "spec"}
	}
}

func PodMetadata(kind string) []string {
	switch kind {
	case "Pod", "Namespace", "Secret":
		return []string{"metadata"}
	case "CronJob":
		return []string{"spec", "jobTemplate", "spec", "template", "metadata"}
	default:
		return []string{"spec", "template", "metadata"}
	}
}

// InspectWorkload - // DEPRECATED
func InspectWorkload(workload interface{}, scopes ...string) (val interface{}, k bool) {
	return InspectMap(workload, scopes...)
}

func getReplicasetNameFromPod(podName, podTemplateHash string) (string, error) {
	// Example pod name: pod-name-123-456
	// Example podTemplateHash: 123
	// Return value: pod-name-123

	// Extract the pod base name without the template hash.
	i := strings.Index(podName, podTemplateHash)
	if i == -1 || len(podName) <= i-1 {
		return "", fmt.Errorf("failed to get replicaset name from pod name: %s", podName)
	}
	basePodName := podName[:i-1]

	// Construct the replicaset name by appending the template hash.
	replicasetName := fmt.Sprintf("%s-%s", basePodName, podTemplateHash)

	if replicasetName == "" || replicasetName == podName {
		return "", fmt.Errorf("failed to get replicaset name from pod name: %s", podName)
	}
	return replicasetName, nil
}
