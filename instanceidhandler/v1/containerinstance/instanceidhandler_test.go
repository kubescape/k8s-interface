package containerinstance

import (
	_ "embed"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/names"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	//go:embed testdata/deployment.json
	deployment string
	//go:embed testdata/jobPod.json
	jobPod string
	//go:embed testdata/mockPod.json
	mockPod string
)

func checkAllsFunctions(t *testing.T, object, apiversion, namespace, kind, name, containerName, formattedString, expectedHash string, expectedLabels map[string]string) error {

	podWorkload, err := workloadinterface.NewWorkload([]byte(object))
	require.NoError(t, err)
	ownerReferences, err := podWorkload.GetOwnerReferences()
	require.NoError(t, err)
	var ownerReference *metav1.OwnerReference
	if len(ownerReferences) > 0 {
		ownerReference = &ownerReferences[0]
	}
	containers, err := podWorkload.GetContainers()
	require.NoError(t, err)
	podWorkloadInstanceID, err := ListInstanceIDs(ownerReference, containers, container, podWorkload.GetApiVersion(), podWorkload.GetNamespace(), podWorkload.GetKind(), podWorkload.GetName(), "")
	require.NoError(t, err)

	assert.Equal(t, 1, len(podWorkloadInstanceID))

	assert.Equal(t, podWorkloadInstanceID[0].ApiVersion, apiversion)
	assert.Equal(t, podWorkloadInstanceID[0].Namespace, namespace)
	assert.Equal(t, podWorkloadInstanceID[0].Kind, kind)
	assert.Equal(t, podWorkloadInstanceID[0].Name, name)
	assert.Equal(t, podWorkloadInstanceID[0].InstanceType, container)
	assert.Equal(t, podWorkloadInstanceID[0].ContainerName, containerName)
	assert.Equal(t, podWorkloadInstanceID[0].GetStringFormatted(), formattedString)
	assert.Equal(t, podWorkloadInstanceID[0].GetHashed(), expectedHash)

	assert.Equal(t, podWorkloadInstanceID[0].GetLabels(), expectedLabels)
	return nil
}

func TestInstanceID(t *testing.T) {
	expectedLabels := map[string]string{
		helpers.ApiGroupMetadataKey:      "apps",
		helpers.ApiVersionMetadataKey:    "v1",
		helpers.NamespaceMetadataKey:     "default",
		helpers.KindMetadataKey:          "ReplicaSet",
		helpers.NameMetadataKey:          "nginx-84f5585d68",
		helpers.ContainerNameMetadataKey: "nginx",
	}

	err := checkAllsFunctions(t, deployment, "apps/v1", "default", "ReplicaSet", "nginx-84f5585d68", "nginx", "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/containerName-nginx", "57366ade3da2e7ba01f8b78251cb57bd70840939f4f207da91cb092b30c06feb", expectedLabels)
	assert.NoError(t, err)

	expectedLabels = map[string]string{
		helpers.ApiGroupMetadataKey:      "batch",
		helpers.ApiVersionMetadataKey:    "v1",
		helpers.NamespaceMetadataKey:     "default",
		helpers.KindMetadataKey:          "Job",
		helpers.NameMetadataKey:          "nginx-job",
		helpers.ContainerNameMetadataKey: "nginx-job",
	}
	err = checkAllsFunctions(t, jobPod, "batch/v1", "default", "Job", "nginx-job", "nginx-job", "apiVersion-batch/v1/namespace-default/kind-Job/name-nginx-job/containerName-nginx-job", "1fdef304b3383588f0e8a267914746de2bf03e1652908d57232cd543a87541c5", expectedLabels)
	assert.NoError(t, err)

	expectedLabels = map[string]string{
		helpers.ApiGroupMetadataKey:      "",
		helpers.ApiVersionMetadataKey:    "v1",
		helpers.NamespaceMetadataKey:     "default",
		helpers.KindMetadataKey:          "Pod",
		helpers.NameMetadataKey:          "nginx",
		helpers.ContainerNameMetadataKey: "nginx",
	}
	err = checkAllsFunctions(t, mockPod, "v1", "default", "Pod", "nginx", "nginx", "apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerName-nginx", "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf", expectedLabels)
	assert.NoError(t, err)
}

func TestInstanceIDToDisplayName(t *testing.T) {
	tt := []struct {
		name    string
		input   *InstanceID
		want    string
		wantErr error
	}{
		{
			name: "valid instanceID produces matching display name",
			input: &InstanceID{
				ApiVersion:    "v1",
				Namespace:     "default",
				Kind:          "Pod",
				Name:          "reverse-proxy",
				ContainerName: "nginx",
				InstanceType:  container,
			},
			want:    "pod-reverse-proxy-nginx-2f07-68bd",
			wantErr: nil,
		},
		{
			name: "valid instanceID produces matching display name",
			input: &InstanceID{
				ApiVersion:    "v1",
				Namespace:     "default",
				Kind:          "Service",
				Name:          "webapp",
				ContainerName: "leader",
				InstanceType:  container,
			},
			want:    "service-webapp-leader-cca3-8ea7",
			wantErr: nil,
		},
		{
			name: "valid instanceID without container name produces matching display name",
			input: &InstanceID{
				ApiVersion: "v1",
				Namespace:  "default",
				Kind:       "Service",
				Name:       "webapp",
			},
			want:    "service-webapp",
			wantErr: nil,
		},
		{
			name: "invalid instanceID produces matching error",
			input: &InstanceID{
				ApiVersion:    "v1",
				Namespace:     "default",
				Kind:          "Service",
				Name:          "web/app",
				ContainerName: "leader",
				InstanceType:  container,
			},
			want:    "",
			wantErr: names.ErrInvalidSlug,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.GetSlug(false)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, tc.wantErr, err)
		})
	}
}
