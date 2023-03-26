package instanceidhandler

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kubescape/k8s-interface/workloadinterface"
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
		return fmt.Errorf("TestCreate: pod instance ID should be created")
	}
	if len(podWorkloadInstanceID) != 1 {
		return fmt.Errorf("TestCreate: should return only one ")
	}

	expected := apiversion
	if podWorkloadInstanceID[0].GetAPIVersion() != expected {
		return fmt.Errorf("TestCreate: wrong , get %s, expected %s", podWorkloadInstanceID[0].GetAPIVersion(), expected)
	}
	expected = namespace
	if podWorkloadInstanceID[0].GetNamespace() != expected {
		return fmt.Errorf("TestCreate: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	expected = kind
	if podWorkloadInstanceID[0].GetKind() != expected {
		return fmt.Errorf("TestCreate: wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetKind(), expected)
	}
	expected = name
	if !strings.Contains(podWorkloadInstanceID[0].GetName(), expected) {
		return fmt.Errorf("TestCreate: wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetName(), expected)
	}
	expected = containerName
	if !strings.Contains(podWorkloadInstanceID[0].GetContainerName(), expected) {
		return fmt.Errorf("TestCreate: wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetContainerName(), expected)
	}
	expected = formattedString
	format := podWorkloadInstanceID[0].GetStringFormatted()
	if format != expected {
		return fmt.Errorf("TestCreate: fail to format Instance ID in string format, get %s, expected %s", podWorkloadInstanceID[0].GetStringFormatted(), expected)
	}
	expected = expectedHash
	if podWorkloadInstanceID[0].GetIDHashed() != expected {
		return fmt.Errorf("TestCreate: GetHashed, get %s, expected %s", podWorkloadInstanceID[0].GetIDHashed(), expected)
	}

	labels := podWorkloadInstanceID[0].GetLabels()
	eq := reflect.DeepEqual(labels, expectedLabels)
	if !eq {
		return fmt.Errorf("TestCreate: GetLabels failed, get %s, expected %s", podWorkloadInstanceID[0].GetLabels(), expectedLabels)
	}

	expected = "123"
	podWorkloadInstanceID[0].SetAPIVersion(expected)
	if podWorkloadInstanceID[0].GetAPIVersion() != expected {
		return fmt.Errorf("TestCreate: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetNamespace(expected)
	if podWorkloadInstanceID[0].GetNamespace() != expected {
		return fmt.Errorf("TestCreate: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetKind(expected)
	if podWorkloadInstanceID[0].GetKind() != expected {
		return fmt.Errorf("TestCreate: wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetKind(), expected)
	}
	podWorkloadInstanceID[0].SetName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetName(), expected) {
		return fmt.Errorf("TestCreate: wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetName(), expected)
	}
	podWorkloadInstanceID[0].SetContainerName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetContainerName(), expected) {
		return fmt.Errorf("TestCreate: wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetContainerName(), expected)
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
		LabelFormatKeyApiGroup:      "",
		LabelFormatKeyApiVersion:    "v1",
		LabelFormatKeyNamespace:     "default",
		LabelFormatKeyKind:          "ReplicaSet",
		LabelFormatKeyName:          "nginx-84f5585d68",
		LabelFormatKeyContainerName: "nginx",
	}

	err = checkAllsFunctions(deployment, "v1", "default", "ReplicaSet", "nginx-84f5585d68", "nginx", "apiVersion-v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/containerName-nginx", "1e1d6a960736b854844e98664e87f7bc6e43c84c04db55a952afe31e2b805689", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		LabelFormatKeyApiGroup:      "",
		LabelFormatKeyApiVersion:    "v1",
		LabelFormatKeyNamespace:     "default",
		LabelFormatKeyKind:          "Job",
		LabelFormatKeyName:          "nginx-job",
		LabelFormatKeyContainerName: "nginx-job",
	}
	err = checkAllsFunctions(jobPod, "v1", "default", "Job", "nginx", "nginx-job", "apiVersion-v1/namespace-default/kind-Job/name-nginx-job/containerName-nginx-job", "031d32a8c548dccfee4d3694890d36a44d4c8a6a5a4f689d0341ba9930e2e3ee", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		LabelFormatKeyApiGroup:      "",
		LabelFormatKeyApiVersion:    "v1",
		LabelFormatKeyNamespace:     "default",
		LabelFormatKeyKind:          "Pod",
		LabelFormatKeyName:          "nginx",
		LabelFormatKeyContainerName: "nginx",
	}
	err = checkAllsFunctions(mockPod, "v1", "default", "Pod", "nginx", "nginx", "apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerName-nginx", "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf", expectedLabels)
	if err != nil {
		t.Error(err)
	}

}
