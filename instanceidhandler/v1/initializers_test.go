package instanceidhandler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core1 "k8s.io/api/core/v1"
)

var (
	//go:embed testdata/cronjob.json
	cronjob string
	//go:embed testdata/cronjob1.json
	cronjob1 string
	//go:embed testdata/cronjob2.json
	cronjob2 string
	//go:embed testdata/cronjob3.json
	cronjob3 string
	//go:embed testdata/deployment.json
	deployment string
	//go:embed testdata/jobPod.json
	jobPod string
	//go:embed testdata/mockPod.json
	mockPod string
)

func TestGenerateInstanceID(t *testing.T) {
	tests := []struct {
		name      string
		sWorkload string
		want      []instanceidhandler.IInstanceID
		wantSlug  string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "deployment",
			sWorkload: deployment,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "nginx",
					InstanceType:  Container,
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "bla",
					InstanceType:  InitContainer,
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "abc",
					InstanceType:  EphemeralContainer,
				},
			},
			wantSlug: "replicaset-nginx-84f5585d68",
			wantErr:  assert.NoError,
		},
		{
			name:      "job",
			sWorkload: jobPod,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "batch/v1",
					Namespace:     "default",
					Kind:          "Job",
					Name:          "nginx-job",
					ContainerName: "nginx-job",
					InstanceType:  Container,
				},
			},
			wantSlug: "job-nginx-job",
			wantErr:  assert.NoError,
		},
		{
			name:      "cronjob",
			sWorkload: cronjob,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "batch/v1",
					Namespace:     "kubescape",
					Kind:          "Job",
					Name:          "kubevuln-scheduler-28677846",
					ContainerName: "kubevuln-scheduler",
					InstanceType:  Container,
					AlternateName: "kubevuln-scheduler-b449cf78f",
				},
			},
			wantSlug: "job-kubevuln-scheduler-b449cf78f",
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp, err := workloadinterface.NewWorkload([]byte(tt.sWorkload))
			require.NoError(t, err)
			got, err := GenerateInstanceID(wp)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateInstanceID - %s", tt.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateInstanceID - %s", tt.name)
			instanceID := got[0].(*containerinstance.InstanceID)
			// replicate filtered SBOM and application profile slug generation
			slug, err := instanceID.GetSlug(true)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantSlug, slug, "GenerateInstanceID - %s", tt.name)
		})
	}
}

// TestSameSlug ensures generated slug stays the same for subsequent cronjob runs
func TestSameSlug(t *testing.T) {
	for _, s := range []string{cronjob1, cronjob2, cronjob3} {
		wp, err := workloadinterface.NewWorkload([]byte(s))
		require.NoError(t, err)
		ins, err := GenerateInstanceID(wp)
		require.NoError(t, err)
		slug, err := ins[0].(*containerinstance.InstanceID).GetSlug(true)
		require.NoError(t, err)
		assert.Equal(t, "job-hello-65866c76d6", slug)
		slugFull, err := ins[0].(*containerinstance.InstanceID).GetSlug(false)
		require.NoError(t, err)
		assert.Equal(t, "job-hello-65866c76d6-hello-d9a7-58fb", slugFull)
	}
}

// Test_InitInstanceID tests the instance id initialization
func TestInitInstanceID(t *testing.T) {
	wp, err := workloadinterface.NewWorkload([]byte(mockPod))
	require.NoError(t, err)
	insFromWorkload, err := GenerateInstanceID(wp)
	require.NoError(t, err)

	p := &core1.Pod{}
	err = json.Unmarshal([]byte(mockPod), p)
	require.NoError(t, err)
	insFromPod, err := GenerateInstanceIDFromPod(p)
	require.NoError(t, err)

	assert.NotEqual(t, 0, len(insFromWorkload))
	assert.Equal(t, len(insFromWorkload), len(insFromPod))

	for i := range insFromWorkload {
		compare(t, insFromWorkload[i].(*containerinstance.InstanceID), insFromPod[i].(*containerinstance.InstanceID))
	}

	insFromString, err := GenerateInstanceIDFromString("apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerName-nginx") //insFromWorkload[0].GetStringFormatted())
	require.NoError(t, err)
	compare(t, insFromWorkload[0].(*containerinstance.InstanceID), insFromString.(*containerinstance.InstanceID))
}

func compare(t *testing.T, a, b *containerinstance.InstanceID) {
	assert.Equal(t, a.GetHashed(), b.GetHashed())
	assert.Equal(t, a.GetStringFormatted(), b.GetStringFormatted())

	assert.Equal(t, a.ApiVersion, b.ApiVersion)
	assert.Equal(t, a.Namespace, b.Namespace)
	assert.Equal(t, a.Kind, b.Kind)
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.ContainerName, b.ContainerName)
}
