package instanceidhandler

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"

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
		ApiGroupMetadataKey:      "apps",
		ApiVersionMetadataKey:    "v1",
		NamespaceMetadataKey:     "default",
		KindMetadataKey:          "ReplicaSet",
		NameMetadataKey:          "nginx-84f5585d68",
		ContainerNameMetadataKey: "nginx",
	}

	err = checkAllsFunctions(deployment, "apps/v1", "default", "ReplicaSet", "nginx-84f5585d68", "nginx", "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/containerName-nginx", "57366ade3da2e7ba01f8b78251cb57bd70840939f4f207da91cb092b30c06feb", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		ApiGroupMetadataKey:      "batch",
		ApiVersionMetadataKey:    "v1",
		NamespaceMetadataKey:     "default",
		KindMetadataKey:          "Job",
		NameMetadataKey:          "nginx-job",
		ContainerNameMetadataKey: "nginx-job",
	}
	err = checkAllsFunctions(jobPod, "batch/v1", "default", "Job", "nginx", "nginx-job", "apiVersion-batch/v1/namespace-default/kind-Job/name-nginx-job/containerName-nginx-job", "1fdef304b3383588f0e8a267914746de2bf03e1652908d57232cd543a87541c5", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		ApiGroupMetadataKey:      "",
		ApiVersionMetadataKey:    "v1",
		NamespaceMetadataKey:     "default",
		KindMetadataKey:          "Pod",
		NameMetadataKey:          "nginx",
		ContainerNameMetadataKey: "nginx",
	}
	err = checkAllsFunctions(mockPod, "v1", "default", "Pod", "nginx", "nginx", "apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerName-nginx", "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf", expectedLabels)
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
				apiVersion:    "v1",
				namespace:     "default",
				kind:          "Pod",
				name:          "reverse-proxy",
				containerName: "nginx",
			},
			want:    "default-Pod-reverse-proxy-2f07-68bd",
			wantErr: nil,
		},
		{
			name: "valid instanceID produces matching display name",
			input: &InstanceID{
				apiVersion:    "v1",
				namespace:     "default",
				kind:          "Service",
				name:          "webapp",
				containerName: "leader",
			},
			want:    "default-Service-webapp-cca3-8ea7",
			wantErr: nil,
		},
		{
			name: "invalid instanceID produces matching error",
			input: &InstanceID{
				apiVersion:    "v1",
				namespace:     "default",
				kind:          "Service",
				name:          "web/app",
				containerName: "leader",
			},
			want:    "",
			wantErr: names.ErrInvalidFriendlyName,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.GetFriendlyName()

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, tc.wantErr, err)
		})
	}
}
