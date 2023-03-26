package workloadinterface

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const (
	mockPodFromDeployment = `{"apiVersion":"v1","kind":"Pod","metadata":{"creationTimestamp":"2023-03-26T09:10:32Z","generateName":"nginx-84f5585d68-","labels":{"app":"nginx","pod-template-hash":"84f5585d68"},"name":"nginx-84f5585d68-42mxd","namespace":"default","ownerReferences":[{"apiVersion":"apps/v1","blockOwnerDeletion":true,"controller":true,"kind":"ReplicaSet","name":"nginx-84f5585d68","uid":"1747930a-61ba-4ee1-b1a9-838af6cb094d"}],"resourceVersion":"8051","uid":"fa771236-d117-453c-b5dc-af14127cc648"},"spec":{"containers":[{"env":[{"name":"DEMO_GREETING","value":"Hellofromtheenvironment"}],"image":"nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d","imagePullPolicy":"IfNotPresent","name":"nginx","resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","volumeMounts":[{"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount","name":"kube-api-access-g42z4","readOnly":true}]}],"dnsPolicy":"ClusterFirst","enableServiceLinks":true,"nodeName":"minikube","preemptionPolicy":"PreemptLowerPriority","priority":0,"restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"serviceAccount":"default","serviceAccountName":"default","terminationGracePeriodSeconds":30,"tolerations":[{"effect":"NoExecute","key":"node.kubernetes.io/not-ready","operator":"Exists","tolerationSeconds":300},{"effect":"NoExecute","key":"node.kubernetes.io/unreachable","operator":"Exists","tolerationSeconds":300}],"volumes":[{"name":"kube-api-access-g42z4","projected":{"defaultMode":420,"sources":[{"serviceAccountToken":{"expirationSeconds":3607,"path":"token"}},{"configMap":{"items":[{"key":"ca.crt","path":"ca.crt"}],"name":"kube-root-ca.crt"}},{"downwardAPI":{"items":[{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"},"path":"namespace"}]}}]}}]},"status":{"conditions":[{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T09:10:32Z","status":"True","type":"Initialized"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T09:10:35Z","status":"True","type":"Ready"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T09:10:35Z","status":"True","type":"ContainersReady"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T09:10:32Z","status":"True","type":"PodScheduled"}],"containerStatuses":[{"containerID":"docker://f41bb9e0b603517517603bab82762e33a95fb5f8fad62c9152089875e08e0b38","image":"sha256:295c7be079025306c4f1d65997fcf7adb411c88f139ad1d34b537164aa060369","imageID":"docker-pullable://nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d","lastState":{},"name":"nginx","ready":true,"restartCount":0,"started":true,"state":{"running":{"startedAt":"2023-03-26T09:10:35Z"}}}],"hostIP":"192.168.49.2","phase":"Running","podIP":"10.244.0.37","podIPs":[{"ip":"10.244.0.37"}],"qosClass":"BestEffort","startTime":"2023-03-26T09:10:32Z"}}`
	mockPodFromJob        = `{"apiVersion":"v1","kind":"Pod","metadata":{"creationTimestamp":"2023-03-26T12:17:45Z","finalizers":["batch.kubernetes.io/job-tracking"],"generateName":"nginx-job-","labels":{"controller-uid":"1ca79890-1432-48fd-8e04-bd6189c194b7","job-name":"nginx-job"},"name":"nginx-job-88b9z","namespace":"default","ownerReferences":[{"apiVersion":"batch/v1","blockOwnerDeletion":true,"controller":true,"kind":"Job","name":"nginx-job","uid":"1ca79890-1432-48fd-8e04-bd6189c194b7"}],"resourceVersion":"17001","uid":"324de883-e0cc-4c8e-b0fb-1ed6aa070c88"},"spec":{"containers":[{"image":"nginx","imagePullPolicy":"Always","name":"nginx-job","resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","volumeMounts":[{"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount","name":"kube-api-access-8b7q6","readOnly":true}]}],"dnsPolicy":"ClusterFirst","enableServiceLinks":true,"nodeName":"minikube","preemptionPolicy":"PreemptLowerPriority","priority":0,"restartPolicy":"Never","schedulerName":"default-scheduler","securityContext":{},"serviceAccount":"default","serviceAccountName":"default","terminationGracePeriodSeconds":30,"tolerations":[{"effect":"NoExecute","key":"node.kubernetes.io/not-ready","operator":"Exists","tolerationSeconds":300},{"effect":"NoExecute","key":"node.kubernetes.io/unreachable","operator":"Exists","tolerationSeconds":300}],"volumes":[{"name":"kube-api-access-8b7q6","projected":{"defaultMode":420,"sources":[{"serviceAccountToken":{"expirationSeconds":3607,"path":"token"}},{"configMap":{"items":[{"key":"ca.crt","path":"ca.crt"}],"name":"kube-root-ca.crt"}},{"downwardAPI":{"items":[{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"},"path":"namespace"}]}}]}}]},"status":{"conditions":[{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T12:17:45Z","status":"True","type":"Initialized"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T12:18:05Z","status":"True","type":"Ready"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T12:18:05Z","status":"True","type":"ContainersReady"},{"lastProbeTime":null,"lastTransitionTime":"2023-03-26T12:17:45Z","status":"True","type":"PodScheduled"}],"containerStatuses":[{"containerID":"docker://3b15c1346e9054ba3da3f87159cc22ede0f0cce8e309f28022055bb3860b500b","image":"nginx:latest","imageID":"docker-pullable://nginx@sha256:f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c","lastState":{},"name":"nginx-job","ready":true,"restartCount":0,"started":true,"state":{"running":{"startedAt":"2023-03-26T12:18:03Z"}}}],"hostIP":"192.168.49.2","phase":"Running","podIP":"10.244.0.38","podIPs":[{"ip":"10.244.0.38"}],"qosClass":"BestEffort","startTime":"2023-03-26T12:17:45Z"}}`
)

func checkAllInstanceIDsFunctions(object string, apiversion, namespace, kind, name, containerName, formattedString, expectedHash string, expectedLabels map[string]string) error {
	podWorkload, err := NewWorkload([]byte(object))
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	podWorkloadInstanceID, err := podWorkload.CreateInstanceID()
	if err != nil {
		return fmt.Errorf("TestCreateInstanceID: pod instance ID should be created")
	}
	if len(podWorkloadInstanceID) != 1 {
		return fmt.Errorf("TestCreateInstanceID: should return only one instanceID")
	}

	expected := apiversion
	if podWorkloadInstanceID[0].GetInstanceIDAPIVersion() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong instanceID, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDAPIVersion(), expected)
	}
	expected = namespace
	if podWorkloadInstanceID[0].GetInstanceIDNamespace() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDNamespace(), expected)
	}
	expected = kind
	if podWorkloadInstanceID[0].GetInstanceIDKind() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDKind(), expected)
	}
	expected = name
	if !strings.Contains(podWorkloadInstanceID[0].GetInstanceIDName(), expected) {
		return fmt.Errorf("TestCreateInstanceID: wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDName(), expected)
	}
	expected = containerName
	if !strings.Contains(podWorkloadInstanceID[0].GetInstanceIDContainerName(), expected) {
		return fmt.Errorf("TestCreateInstanceID: wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDContainerName(), expected)
	}
	expected = formattedString
	formatInstanceID := podWorkloadInstanceID[0].GetInstanceIDStringFormatted()
	if formatInstanceID != expected {
		return fmt.Errorf("TestCreateInstanceID: fail to format Instance ID in string format, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDStringFormatted(), expected)
	}
	expected = expectedHash
	if podWorkloadInstanceID[0].GetInstanceIDHashed() != expected {
		return fmt.Errorf("TestCreateInstanceID: GetInstanceIDHashed, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDHashed(), expected)
	}

	labels := podWorkloadInstanceID[0].GetInstanceIDLabels()
	eq := reflect.DeepEqual(labels, expectedLabels)
	if !eq {
		return fmt.Errorf("TestCreateInstanceID: GetInstanceIDLabels failed, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDLabels(), expectedLabels)
	}

	expected = "123"
	podWorkloadInstanceID[0].SetInstanceIDAPIVersion(expected)
	if podWorkloadInstanceID[0].GetInstanceIDAPIVersion() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetInstanceIDNamespace(expected)
	if podWorkloadInstanceID[0].GetInstanceIDNamespace() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong namespace, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDNamespace(), expected)
	}
	podWorkloadInstanceID[0].SetInstanceIDKind(expected)
	if podWorkloadInstanceID[0].GetInstanceIDKind() != expected {
		return fmt.Errorf("TestCreateInstanceID: wrong parent kind, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDKind(), expected)
	}
	podWorkloadInstanceID[0].SetInstanceIDName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetInstanceIDName(), expected) {
		return fmt.Errorf("TestCreateInstanceID: wrong parent name, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDName(), expected)
	}
	podWorkloadInstanceID[0].SetInstanceIDContainerName(expected)
	if !strings.Contains(podWorkloadInstanceID[0].GetInstanceIDContainerName(), expected) {
		return fmt.Errorf("TestCreateInstanceID: wrong container name, get %s, expected %s", podWorkloadInstanceID[0].GetInstanceIDContainerName(), expected)
	}
	return nil
}

func TestCreateInstanceID(t *testing.T) {
	serviceWorkload, err := NewWorkload([]byte(mockService))
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = serviceWorkload.CreateInstanceID()
	if err == nil {
		t.Errorf("can't create instance ID from service")
	}
	expectedLabels := map[string]string{
		instanceIDLabelFormatKeyApiGroup:      "",
		instanceIDLabelFormatKeyApiVersion:    "v1",
		instanceIDLabelFormatKeyNamespace:     "default",
		instanceIDLabelFormatKeyKind:          "ReplicaSet",
		instanceIDLabelFormatKeyName:          "nginx-84f5585d68",
		instanceIDLabelFormatKeyContainerName: "nginx",
	}

	err = checkAllInstanceIDsFunctions(mockPodFromDeployment, "v1", "default", "ReplicaSet", "nginx-84f5585d68", "nginx", "apiVersion-v1/namespace-default/kind-ReplicaSet/name-nginx-84f5585d68/containerName-nginx", "1e1d6a960736b854844e98664e87f7bc6e43c84c04db55a952afe31e2b805689", expectedLabels)
	if err != nil {
		t.Error(err)
	}

	expectedLabels = map[string]string{
		instanceIDLabelFormatKeyApiGroup:      "",
		instanceIDLabelFormatKeyApiVersion:    "v1",
		instanceIDLabelFormatKeyNamespace:     "default",
		instanceIDLabelFormatKeyKind:          "Job",
		instanceIDLabelFormatKeyName:          "nginx-job",
		instanceIDLabelFormatKeyContainerName: "nginx-job",
	}
	err = checkAllInstanceIDsFunctions(mockPodFromJob, "v1", "default", "Job", "nginx", "nginx-job", "apiVersion-v1/namespace-default/kind-Job/name-nginx-job/containerName-nginx-job", "031d32a8c548dccfee4d3694890d36a44d4c8a6a5a4f689d0341ba9930e2e3ee", expectedLabels)
	if err != nil {
		t.Error(err)
	}
}

func CurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}
