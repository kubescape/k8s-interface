package workloadinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestInspectMap(t *testing.T) {
	mapexample := map[string]interface{}{
		"name": "a",
		"relatedObjects": struct {
			role        string
			rolebinding string
		}{
			"role-a", "rolebinding-b",
		},
		"namespace": "c",
	}
	if name, ok := InspectMap(mapexample, "name"); ok {
		if name != "a" {
			t.Errorf("error in inspectMap, name %s != 'a'", name)
		}
	} else {
		t.Errorf("error in inspectMap, couldn't find name: '%s", name)
	}
}

// Update a container using the "SetInMap" function
func TestSetInMap(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"demoservice-server"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"demoservice-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}}}`
	workload, _ := NewWorkload([]byte(w))
	c, _ := workload.GetContainers()
	// obj := workload.GetObject()
	SetInMap(workload.GetObject(), PodSpec(workload.GetKind()), "containers", append(c, corev1.Container{Name: "bla"}))
	wc, _ := workload.GetContainers()
	assert.Equal(t, 2, len(wc))
	assert.Equal(t, "demoservice", wc[0].Name)
	assert.Equal(t, "bla", wc[1].Name)
}

func TestRemoveFromMap(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"demoservice-server"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"demoservice-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}}}`
	workload, _ := NewWorkload([]byte(w))
	RemoveFromMap(workload.GetObject(), append(PodSpec(workload.GetKind()), "containers")...)
	wc, _ := workload.GetContainers()
	assert.Equal(t, 0, len(wc))
}
