package workloadinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ========================================= IS =========================================

var mockDeployment = `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
var mockService = `{"apiVersion":"v1","kind":"Service","metadata":{"creationTimestamp":"2021-12-06T14:01:16Z","labels":{"app":"armo-vuln-scan","app.kubernetes.io\/managed-by":"Helm"},"name":"armo-vuln-scan","resourceVersion":"351796","uid":"12bd4f9f-3ec6-4113-8ec6-0b8a1c772deb"},"spec":{"clusterIP":"10.107.7.78","clusterIPs":["10.107.7.78"],"internalTrafficPolicy":"Cluster","ipFamilies":["IPv4"],"ipFamilyPolicy":"SingleStack","ports":[{"port":8080,"protocol":"TCP","targetPort":8080}],"selector":{"app":"armo-vuln-scan"},"sessionAffinity":"None","type":"ClusterIP"},"status":{"loadBalancer":{}}}`

func TestLabels(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Errorf(err.Error())
	}
	if workload.GetKind() != "Deployment" {
		t.Errorf("wrong kind")
	}
	if workload.GetNamespace() != "default" {
		t.Errorf("wrong namespace")
	}
	if workload.GetName() != "demoservice-server" {
		t.Errorf("wrong name")
	}
}

func TestGetReplicas(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Errorf(err.Error())
	}
	SetInMap(workload.GetObject(), []string{"spec"}, "replicas", 3)
	assert.Equal(t, 3, workload.GetReplicas())
}
func TestGetSelector(t *testing.T) {
	workload, err := NewWorkload([]byte(mockDeployment))
	if err != nil {
		t.Errorf(err.Error())
	}
	l, err := workload.GetSelector()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(l.MatchLabels))
	assert.Equal(t, 0, len(l.MatchExpressions))
}
func TestSetNamespace(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"demoservice-server"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"demoservice-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	workload.SetNamespace("default")
	if workload.GetNamespace() != "default" {
		t.Errorf("wrong namespace")
	}
}
func TestSetLabels(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	workload.SetLabel("bla", "daa")
	v, ok := workload.GetLabel("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find label")
	}
	workload.RemoveLabel("bla")
	v2, ok2 := workload.GetLabel("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("label not deleted")
	}
}

func TestSetAnnotations(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	workload.SetAnnotation("bla", "daa")
	v, ok := workload.GetAnnotation("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find annotation")
	}
	workload.RemoveAnnotation("bla")
	v2, ok2 := workload.GetAnnotation("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("annotation not deleted")
	}
}
func TestSetPodLabels(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	workload.SetPodLabel("bla", "daa")
	v, ok := workload.GetPodLabel("bla")
	if !ok || v != "daa" {
		t.Errorf("expect to find label")
	}
	workload.RemovePodLabel("bla")
	v2, ok2 := workload.GetPodLabel("bla")
	if ok2 || v2 == "daa" {
		t.Errorf("label not deleted")
	}
}

func TestGetResourceVersion(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	if workload.GetResourceVersion() != "1016043" {
		t.Errorf("wrong resourceVersion")
	}

}
func TestGetUID(t *testing.T) {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, err := NewWorkload([]byte(w))
	if err != nil {
		t.Errorf(err.Error())
	}
	if workload.GetUID() != "e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e" {
		t.Errorf("wrong UID")
	}

}
func TestGetID(t *testing.T) {
	dWorkload, _ := NewWorkload([]byte(mockDeployment))
	assert.Equal(t, "apps/v1/default/Deployment/demoservice-server", dWorkload.GetID())

	sWorkload, _ := NewWorkload([]byte(mockService))
	assert.Equal(t, "/v1//Service/armo-vuln-scan", sWorkload.GetID())

}
