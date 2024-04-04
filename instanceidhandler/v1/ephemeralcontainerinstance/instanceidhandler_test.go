package ephemeralcontainerinstance

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/kubescape/k8s-interface/names"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/service.json
	service string
	//go:embed testdata/deployment.json
	deployment string
	//go:embed testdata/jobPod.json
	jobPod string
	//go:embed testdata/mockPod.json
	mockPod string
)

func checkAllsFunctions(object string, apiversion, namespace, kind, name, containerName, formattedString, expectedHash string, expectedLabels map[string]string) error {

	podWorkload, err := workloadinterface.NewWorkload([]byte(object))
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	podWorkloadInstanceID, err := GenerateInstanceID(podWorkload)
	if err != nil {
		return fmt.Errorf("TestCreate: GenerateInstanceID - pod instance ID should be created")
	}
	if len(podWorkloadInstanceID) != 1 {
		return fmt.Errorf("TestCreate: should return only one ")
	}

	expected := apiversion
	if podWorkloadInstanceID[0].GetAPIVersion() != expected {
		return fmt.Errorf("TestCreate: GetAPIVersion - wrong , get %s, expected %s", podWorkloadInstanceID[0].GetAPIVersion(), expected)
	}
	expected = namespace
	if podWorkloadInstanceID[0].GetNamespace() != expected {
		return fmt.Errorf("TestCreate: GetNamespace - wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	expected = kind
	if podWorkloadInstanceID[0].GetKind() != expected {
		return fmt.Errorf("TestCreate: GetKind - wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetKind(), expected)
	}
	expected = name
	if !strings.Contains(podWorkloadInstanceID[0].GetName(), expected) {
		return fmt.Errorf("TestCreate: GetName - wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetName(), expected)
	}
	expected = containerName
	if !strings.Contains(podWorkloadInstanceID[0].GetContainerName(), expected) {
		return fmt.Errorf("TestCreate: GetContainerName - wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetContainerName(), expected)
	}
	expected = formattedString
	format := podWorkloadInstanceID[0].GetStringFormatted()
	if format != expected {
		return fmt.Errorf("TestCreate: GetStringFormatted - fail to format Instance ID in string format, get %s, expected %s", podWorkloadInstanceID[0].GetStringFormatted(), expected)
	}
	expected = expectedHash
	if podWorkloadInstanceID[0].GetHashed() != expected {
		return fmt.Errorf("TestCreate: GetHashed - GetHashed, get %s, expected %s", podWorkloadInstanceID[0].GetHashed(), expected)
	}

	labels := podWorkloadInstanceID[0].GetLabels()
	eq := reflect.DeepEqual(labels, expectedLabels)
	if !eq {
		return fmt.Errorf("TestCreate: GetLabels - GetLabels failed, get %s, expected %s", podWorkloadInstanceID[0].GetLabels(), expectedLabels)
	}

	expected = "123"
	podWorkloadInstanceID[0].SetAPIVersion(expected)
	if podWorkloadInstanceID[0].GetAPIVersion() != expected {
		return fmt.Errorf("TestCreate: SetAPIVersion - wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetNamespace(expected)
	if podWorkloadInstanceID[0].GetNamespace() != expected {
		return fmt.Errorf("TestCreate: SetNamespace - wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetKind(expected)
	if podWorkloadInstanceID[0].GetKind() != expected {
		return fmt.Errorf("TestCreate: SetKind - wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetKind(), expected)
	}
	podWorkloadInstanceID[0].SetName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetName(), expected) {
		return fmt.Errorf("TestCreate: SetName - wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetName(), expected)
	}
	podWorkloadInstanceID[0].SetContainerName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetContainerName(), expected) {
		return fmt.Errorf("TestCreate: SetContainerName - wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetContainerName(), expected)
	}
	return nil
}

func TestInstanceID(t *testing.T) {
	serviceWorkload, err := workloadinterface.NewWorkload([]byte(service))
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = GenerateInstanceID(serviceWorkload)
	if err == nil {
		t.Errorf("can't create instance ID from service")
	}
	expectedLabels := map[string]string{
		helpers.ApiGroupMetadataKey:               "apps",
		helpers.ApiVersionMetadataKey:             "v1",
		helpers.NamespaceMetadataKey:              "default",
		helpers.KindMetadataKey:                   "ReplicaSet",
		helpers.NameMetadataKey:                   "nginx-84f5585d68",
		helpers.EphemeralContainerNameMetadataKey: "nginx",
	}

	err = checkAllsFunctions(deployment, "apps/v1", "default", "ReplicaSet", "nginx-84f5585d68", "nginx", "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/ephemeralContainerName-nginx", "71dd1e22d7246a38b30dc1cb974fe2bf56a6dd0a59299c371b85d1684d3d0f6d", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		helpers.ApiGroupMetadataKey:               "batch",
		helpers.ApiVersionMetadataKey:             "v1",
		helpers.NamespaceMetadataKey:              "default",
		helpers.KindMetadataKey:                   "Job",
		helpers.NameMetadataKey:                   "nginx-job",
		helpers.EphemeralContainerNameMetadataKey: "nginx-job",
	}
	err = checkAllsFunctions(jobPod, "batch/v1", "default", "Job", "nginx", "nginx-job", "apiVersion-batch/v1/namespace-default/kind-Job/name-nginx-job/ephemeralContainerName-nginx-job", "9f451da6492f37b69a61ed95bbb7ef1f4da5a4c191fda69c5fe87ee0b95bc191", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		helpers.ApiGroupMetadataKey:               "",
		helpers.ApiVersionMetadataKey:             "v1",
		helpers.NamespaceMetadataKey:              "default",
		helpers.KindMetadataKey:                   "Pod",
		helpers.NameMetadataKey:                   "nginx",
		helpers.EphemeralContainerNameMetadataKey: "nginx",
	}
	err = checkAllsFunctions(mockPod, "v1", "default", "Pod", "nginx", "nginx", "apiVersion-v1/namespace-default/kind-Pod/name-nginx/ephemeralContainerName-nginx", "458ef20641136a35dcdbbbc25f90824e9f263faccfaa0da4b0da8c766514fc15", expectedLabels)
	if err != nil {
		t.Error(err)
	}

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
				apiVersion:             "v1",
				namespace:              "default",
				kind:                   "Pod",
				name:                   "reverse-proxy",
				ephemeralContainerName: "nginx",
			},
			want:    "pod-reverse-proxy-nginx-36fe-11b1",
			wantErr: nil,
		},
		{
			name: "valid instanceID produces matching display name",
			input: &InstanceID{
				apiVersion:             "v1",
				namespace:              "default",
				kind:                   "Service",
				name:                   "webapp",
				ephemeralContainerName: "leader",
			},
			want:    "service-webapp-leader-d90e-d2d8",
			wantErr: nil,
		},
		{
			name: "valid instanceID without container name produces matching display name",
			input: &InstanceID{
				apiVersion: "v1",
				namespace:  "default",
				kind:       "Service",
				name:       "webapp",
			},
			want:    "service-webapp",
			wantErr: nil,
		},
		{
			name: "invalid instanceID produces matching error",
			input: &InstanceID{
				apiVersion:             "v1",
				namespace:              "default",
				kind:                   "Service",
				name:                   "web/app",
				ephemeralContainerName: "leader",
			},
			want:    "",
			wantErr: names.ErrInvalidSlug,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.GetSlug()

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, tc.wantErr, err)
		})
	}
}
