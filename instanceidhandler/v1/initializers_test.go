package instanceidhandler

import (
	"encoding/json"
	"testing"

	core1 "k8s.io/api/core/v1"

	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
)

// Test_InitInstanceID tests the instance id initialization
func Test_InitInstanceID(t *testing.T) {
	wp, err := workloadinterface.NewWorkload([]byte(mockPod))
	if err != nil {
		t.Fatalf(err.Error())
	}
	insFromWorkload, err := GenerateInstanceID(wp)
	if err != nil {
		t.Fatalf("can't create instance ID from pod")
	}

	p := &core1.Pod{}
	if err := json.Unmarshal([]byte(mockPod), p); err != nil {
		t.Fatalf(err.Error())
	}
	insFromPod, err := GenerateInstanceIDFromPod(p)
	if err != nil {
		t.Fatalf("can't create instance ID from pod")
	}

	assert.NotEqual(t, 0, len(insFromWorkload))
	assert.Equal(t, len(insFromWorkload), len(insFromPod))

	for i := range insFromWorkload {
		compare(t, insFromWorkload[i], insFromPod[i])
	}

	insFromString, err := GenerateInstanceIDFromString(insFromWorkload[0].GetStringFormatted())
	if err != nil {
		t.Fatalf("can't create instance ID from string: %s, error: %s", insFromWorkload[0].GetStringFormatted(), err.Error())
	}
	compare(t, insFromWorkload[0], insFromString)

}

func compare(t *testing.T, a, b *InstanceID) {
	assert.Equal(t, a.GetIDHashed(), b.GetIDHashed())
	assert.Equal(t, a.GetStringFormatted(), b.GetStringFormatted())

	assert.Equal(t, a.apiVersion, b.apiVersion)
	assert.Equal(t, a.namespace, b.namespace)
	assert.Equal(t, a.kind, b.kind)
	assert.Equal(t, a.name, b.name)
	assert.Equal(t, a.containerName, b.containerName)
}
