package instanceidhandler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	//go:embed testdata/cronjob.json
	cronjob    string
	cronjobObj batchv1.CronJob
	_          = json.Unmarshal([]byte(cronjob), &cronjobObj)
	//go:embed testdata/cronjob1.json
	cronjob1    string
	cronjob1Obj batchv1.CronJob
	_           = json.Unmarshal([]byte(cronjob1), &cronjob1Obj)
	//go:embed testdata/cronjob2.json
	cronjob2    string
	cronjob2Obj batchv1.CronJob
	_           = json.Unmarshal([]byte(cronjob2), &cronjob2Obj)
	//go:embed testdata/cronjob3.json
	cronjob3    string
	cronjob3Obj batchv1.CronJob
	_           = json.Unmarshal([]byte(cronjob3), &cronjob3Obj)
	//go:embed testdata/daemonset.json
	daemonset    string
	daemonsetObj corev1.Pod
	_            = json.Unmarshal([]byte(daemonset), &daemonsetObj)
	//go:embed testdata/deployment.json
	deployment    string
	deploymentObj appsv1.Deployment
	_             = json.Unmarshal([]byte(deployment), &deploymentObj)
	//go:embed testdata/jobPod.json
	jobPod    string
	jobPodObj batchv1.Job
	_         = json.Unmarshal([]byte(jobPod), &jobPodObj)
	//go:embed testdata/mockPod.json
	mockPod    string
	mockPodObj corev1.Pod
	_          = json.Unmarshal([]byte(mockPod), &mockPodObj)
	//go:embed testdata/statefulset.json
	statefulset    string
	statefulsetObj appsv1.StatefulSet
	_              = json.Unmarshal([]byte(statefulset), &statefulsetObj)
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
			name:      "daemonset",
			sWorkload: daemonset,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "kubescape",
					Kind:          "DaemonSet",
					Name:          "node-agent",
					ContainerName: "node-agent",
					InstanceType:  Container,
					AlternateName: "node-agent-f9dd7596f",
				},
			},
			wantSlug: "daemonset-node-agent-f9dd7596f",
			wantErr:  assert.NoError,
		},
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
					TemplateHash:  "84f5585d68",
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "bla",
					InstanceType:  InitContainer,
					TemplateHash:  "84f5585d68",
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "abc",
					InstanceType:  EphemeralContainer,
					TemplateHash:  "84f5585d68",
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
					AlternateName: "kubevuln-scheduler-5f99858564",
					TemplateHash:  "5f99858564",
				},
			},
			wantSlug: "job-kubevuln-scheduler-5f99858564",
			wantErr:  assert.NoError,
		},
		{
			name:      "statefulset",
			sWorkload: statefulset,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "StatefulSet",
					Name:          "web",
					ContainerName: "nginx",
					InstanceType:  Container,
					AlternateName: "web-7757fc6447",
				},
			},
			wantSlug: "statefulset-web-7757fc6447",
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp, err := workloadinterface.NewWorkload([]byte(tt.sWorkload))
			require.NoError(t, err)
			got, err := GenerateInstanceID(wp, nil)
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

func TestGenerateInstanceIDFromRuntime(t *testing.T) {
	tests := []struct {
		name      string
		sWorkload string
		want      []instanceidhandler.IInstanceID
		wantSlug  string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "daemonset",
			sWorkload: daemonset,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "kubescape",
					Kind:          "DaemonSet",
					Name:          "node-agent",
					ContainerName: "node-agent",
					InstanceType:  Container,
					AlternateName: "node-agent-f9dd7596f",
				},
			},
			wantSlug: "daemonset-node-agent-f9dd7596f",
			wantErr:  assert.NoError,
		},
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
					TemplateHash:  "84f5585d68",
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "bla",
					InstanceType:  InitContainer,
					TemplateHash:  "84f5585d68",
				},
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "nginx-84f5585d68",
					ContainerName: "abc",
					InstanceType:  EphemeralContainer,
					TemplateHash:  "84f5585d68",
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
					AlternateName: "kubevuln-scheduler-5f99858564",
					TemplateHash:  "5f99858564",
				},
			},
			wantSlug: "job-kubevuln-scheduler-5f99858564",
			wantErr:  assert.NoError,
		},
		{
			name:      "statefulset",
			sWorkload: statefulset,
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "StatefulSet",
					Name:          "web",
					ContainerName: "nginx",
					InstanceType:  Container,
					AlternateName: "web-7757fc6447",
				},
			},
			wantSlug: "statefulset-web-7757fc6447",
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj corev1.Pod
			err := json.Unmarshal([]byte(tt.sWorkload), &obj)
			require.NoError(t, err)
			got, err := GenerateInstanceIDFromRuntimeObj(&obj, nil)
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
		ins, err := GenerateInstanceID(wp, nil)
		require.NoError(t, err)
		slug, err := ins[0].(*containerinstance.InstanceID).GetSlug(true)
		require.NoError(t, err)
		assert.Equal(t, "job-hello-77bdd46fc5", slug)
		slugFull, err := ins[0].(*containerinstance.InstanceID).GetSlug(false)
		require.NoError(t, err)
		assert.Equal(t, "job-hello-77bdd46fc5-hello-6f4e-3f75", slugFull)
	}
}

// Test_InitInstanceID tests the instance id initialization
func TestInitInstanceID(t *testing.T) {
	wp, err := workloadinterface.NewWorkload([]byte(mockPod))
	require.NoError(t, err)
	insFromWorkload, err := GenerateInstanceID(wp, nil)
	require.NoError(t, err)

	p := &corev1.Pod{}
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

func TestDeepHashObject_ExcludedFields(t *testing.T) {
	pod1 := &corev1.Pod{}
	err := json.Unmarshal([]byte(cronjob), pod1)
	require.NoError(t, err)

	pod2 := &corev1.Pod{}
	err = json.Unmarshal([]byte(cronjob), pod2)
	require.NoError(t, err)

	// Simulate pod scheduled to a different node with different dynamic/injected env vars
	pod1.Spec.NodeName = "node-alpha"
	pod1.Spec.Hostname = "host-alpha"

	pod2.Spec.NodeName = "node-beta"
	pod2.Spec.Hostname = "host-beta"
	pod2.Spec.Containers[0].Env[0].Value = "different-uuid-value"
	pod2.Spec.Containers[0].Env[1].Value = "9999999999"

	// Also verify initContainers with transient Datadog env vars
	initEnv1 := []corev1.EnvVar{{Name: "DD_INSTRUMENTATION_INSTALL_ID", Value: "init-id-1"}}
	initEnv2 := []corev1.EnvVar{{Name: "DD_INSTRUMENTATION_INSTALL_ID", Value: "init-id-2"}}
	pod1.Spec.InitContainers = []corev1.Container{{Name: "init-c", Env: initEnv1}}
	pod2.Spec.InitContainers = []corev1.Container{{Name: "init-c", Env: initEnv2}}

	hasher1 := fnv.New32a()
	DeepHashObject(hasher1, &pod1.Spec, nil)
	hash1 := hasher1.Sum32()

	hasher2 := fnv.New32a()
	DeepHashObject(hasher2, &pod2.Spec, nil)
	hash2 := hasher2.Sum32()

	assert.Equal(t, hash1, hash2, "hashes must match even when nodeName or Datadog inject env vars differ")
}
