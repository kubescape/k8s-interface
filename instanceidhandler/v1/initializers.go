package instanceidhandler

import (
	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/initcontainerinstance"
	"github.com/kubescape/k8s-interface/workloadinterface"

	core1 "k8s.io/api/core/v1"
)

// GenerateInstanceID generates instance ID from workload
func GenerateInstanceID(w workloadinterface.IWorkload) ([]instanceidhandler.IInstanceID, error) {
	c, err := containerinstance.GenerateInstanceID(w)
	if err != nil {
		return nil, err
	}
	l := convertContainersToIInstanceID(c)

	initC, err := initcontainerinstance.GenerateInstanceID(w)
	if err != nil {
		return l, err
	}

	l = append(l, convertInitContainersToIInstanceID(initC)...)

	return l, nil
}

// GenerateInstanceIDFromPod generates instance ID from pod
func GenerateInstanceIDFromPod(pod *core1.Pod) ([]instanceidhandler.IInstanceID, error) {

	c, err := containerinstance.GenerateInstanceIDFromPod(pod)
	if err != nil {
		return nil, err
	}
	l := convertContainersToIInstanceID(c)

	initC, err := initcontainerinstance.GenerateInstanceIDFromPod(pod)
	if err != nil {
		return l, err
	}

	l = append(l, convertInitContainersToIInstanceID(initC)...)

	return l, nil
}

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/...
func GenerateInstanceIDFromString(input string) (instanceidhandler.IInstanceID, error) {
	if instID, err := containerinstance.GenerateInstanceIDFromString(input); err == nil && instID != nil {
		return instID, nil
	}
	return initcontainerinstance.GenerateInstanceIDFromString(input)
}

// convert list containerinstance.InstanceID to instanceidhandler.IInstanceID
func convertContainersToIInstanceID(l []containerinstance.InstanceID) []instanceidhandler.IInstanceID {
	li := []instanceidhandler.IInstanceID{}
	for _, i := range l {
		li = append(li, &i)
	}
	return li
}
func convertInitContainersToIInstanceID(l []initcontainerinstance.InstanceID) []instanceidhandler.IInstanceID {
	li := []instanceidhandler.IInstanceID{}
	for _, i := range l {
		li = append(li, &i)
	}
	return li
}
