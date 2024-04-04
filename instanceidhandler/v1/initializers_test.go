package instanceidhandler

import (
	_ "embed"
	"testing"

	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/deployment.json
	deployment string
)

// Test_InitInstanceID tests the instance id initialization
func TestInitInstanceID(t *testing.T) {
	wp, err := workloadinterface.NewWorkload([]byte(deployment))
	if err != nil {
		t.Fatalf(err.Error())
	}
	insFromWorkload, err := GenerateInstanceID(wp)
	if err != nil {
		t.Fatalf("can't create instance ID from pod")
	}

	assert.Equal(t, 3, len(insFromWorkload))

	assert.Equal(t, "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/containerName-nginx", insFromWorkload[0].GetStringFormatted())
	assert.Equal(t, "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/initContainerName-bla", insFromWorkload[1].GetStringFormatted())
	assert.Equal(t, "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/ephemeralContainerName-abc", insFromWorkload[2].GetStringFormatted())

	s0, _ := insFromWorkload[0].GetSlug()
	assert.Equal(t, "replicaset-nginx-84f5585d68-nginx-5736-6feb", s0)
	s1, _ := insFromWorkload[1].GetSlug()
	assert.Equal(t, "replicaset-nginx-84f5585d68-bla-c2db-8092", s1)
	s2, _ := insFromWorkload[2].GetSlug()
	assert.Equal(t, "replicaset-nginx-84f5585d68-abc-725c-ef45", s2)

}
