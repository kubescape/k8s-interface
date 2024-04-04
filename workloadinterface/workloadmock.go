package workloadinterface

import (
	"github.com/armosec/armoapi-go/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const TypeWorkloadObjectMock ObjectType = "workloadMock"

type WorkloadMock struct {
	workload *Workload
}

func NewWorkloadMock(ww interface{}) *WorkloadMock {
	w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1"},"creationTimestamp":"2021-05-03T13:10:32Z","generation":1,"labels":{"app":"demoservice-server","cyberarmor.inject":"true"},"managedFields":[{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{},"f:cyberarmor.inject":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demoservice\"}":{".":{},"f:env":{".":{},"k:{\"name\":\"ARMO_TEST_NAME\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"CAA_ENABLE_CRASH_REPORTER\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"DEMO_FOLDERS\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SERVER_PORT\"}":{".":{},"f:name":{},"f:value":{}},"k:{\"name\":\"SLEEP_DURATION\"}":{".":{},"f:name":{},"f:value":{}}},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8089,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}},"manager":"OpenAPI-Generator","operation":"Update","time":"2021-05-03T13:10:32Z"},{"apiVersion":"apps/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}},"manager":"kube-controller-manager","operation":"Update","time":"2021-05-03T13:52:58Z"}],"name":"demoservice-server","namespace":"default","resourceVersion":"1016043","uid":"e9e8a3e9-6cb4-4301-ace1-2c0cef3bd61e"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"demoservice-server"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"demoservice-server"}},"spec":{"containers":[{"env":[{"name":"SERVER_PORT","value":"8089"},{"name":"SLEEP_DURATION","value":"1"},{"name":"DEMO_FOLDERS","value":"/app"},{"name":"ARMO_TEST_NAME","value":"auto_attach_deployment"},{"name":"CAA_ENABLE_CRASH_REPORTER","value":"1"}],"image":"quay.io/armosec/demoservice:v25","imagePullPolicy":"IfNotPresent","name":"demoservice","ports":[{"containerPort":8089,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"terminationGracePeriodSeconds":30}}},"status":{"availableReplicas":1,"conditions":[{"lastTransitionTime":"2021-05-03T13:10:32Z","lastUpdateTime":"2021-05-03T13:10:37Z","message":"ReplicaSet \"demoservice-server-7d478b6998\" has successfully progressed.","reason":"NewReplicaSetAvailable","status":"True","type":"Progressing"},{"lastTransitionTime":"2021-05-03T13:52:58Z","lastUpdateTime":"2021-05-03T13:52:58Z","message":"Deployment has minimum availability.","reason":"MinimumReplicasAvailable","status":"True","type":"Available"}],"observedGeneration":1,"readyReplicas":1,"replicas":1,"updatedReplicas":1}}`
	workload, _ := NewWorkload([]byte(w))
	return &WorkloadMock{
		workload: workload,
	}
}

func (wm *WorkloadMock) Json() string {
	return wm.workload.Json()
}
func (wm *WorkloadMock) ToString() string {
	return wm.workload.ToString()
}

func (wm *WorkloadMock) DeepCopy(w map[string]interface{}) {
	wm.workload.DeepCopy(w)

}

func (wm *WorkloadMock) ToUnstructured() (*unstructured.Unstructured, error) {
	return wm.workload.ToUnstructured()
}

// ======================================= DELETE ========================================

func (wm *WorkloadMock) RemoveJobID() {
	wm.workload.RemoveJobID()
}

func (wm *WorkloadMock) RemoveArmoAnnotations() {
	wm.workload.RemoveJobID()
}
func (wm *WorkloadMock) RemoveSecretData() {
	wm.workload.RemoveSecretData()
}

func (wm *WorkloadMock) RemovePodStatus() {
	wm.workload.RemovePodStatus()
}

func (wm *WorkloadMock) RemoveResourceVersion() {
	wm.workload.RemoveResourceVersion()
}

func (wm *WorkloadMock) RemoveLabel(key string) {
	wm.workload.RemoveLabel(key)
}

func (wm *WorkloadMock) RemoveAnnotation(key string) {
	wm.workload.RemoveAnnotation(key)
}

func (wm *WorkloadMock) RemovePodAnnotation(key string) {
	wm.workload.RemovePodAnnotation(key)
}

func (wm *WorkloadMock) RemovePodLabel(key string) {
	wm.workload.RemovePodLabel(key)
}

func (wm *WorkloadMock) RemoveMetadata(scope []string, metadata, key string) {
	wm.workload.RemoveMetadata(scope, metadata, key)
}

// ========================================= SET =========================================

func (wm *WorkloadMock) SetWorkload(workload map[string]interface{}) {
	wm.workload.SetWorkload(workload)
}
func (wm *WorkloadMock) SetObject(workload map[string]interface{}) {
	wm.workload.SetObject(workload)
}

func (wm *WorkloadMock) SetApiVersion(apiVersion string) {
	wm.workload.SetApiVersion(apiVersion)
}

func (wm *WorkloadMock) SetKind(kind string) {
	wm.workload.SetKind(kind)
}
func (wm *WorkloadMock) SetJobID(jobTracking apis.JobTracking) {
	wm.workload.SetJobID(jobTracking)
}

func (wm *WorkloadMock) SetNamespace(namespace string) {
	wm.workload.SetNamespace(namespace)
}

func (wm *WorkloadMock) SetName(name string) {
	wm.workload.SetName(name)
}

func (wm *WorkloadMock) SetLabel(key, value string) {
	wm.workload.SetLabel(key, value)
}

func (wm *WorkloadMock) SetPodLabel(key, value string) {
	wm.workload.SetPodLabel(key, value)
}
func (wm *WorkloadMock) SetAnnotation(key, value string) {
	wm.workload.SetAnnotation(key, value)
}
func (wm *WorkloadMock) SetPodAnnotation(key, value string) {
	wm.workload.SetPodAnnotation(key, value)
}

// ========================================= GET =========================================

func (wm *WorkloadMock) GetObjectType() ObjectType {
	return TypeWorkloadObjectMock
}
func (wm *WorkloadMock) GetWorkload() map[string]interface{} {
	return wm.workload.GetWorkload()
}
func (wm *WorkloadMock) GetObject() map[string]interface{} {
	return wm.workload.GetObject()
}
func (wm *WorkloadMock) GetNamespace() string {
	return wm.workload.GetNamespace()

}
func (wm *WorkloadMock) GetID() string {
	return wm.workload.GetID()
}
func (wm *WorkloadMock) GetName() string {
	return wm.workload.GetName()

}

func (wm *WorkloadMock) GetApiVersion() string {
	return wm.workload.GetApiVersion()
}

func (wm *WorkloadMock) GetVersion() string {
	return wm.workload.GetVersion()
}

func (wm *WorkloadMock) GetGroup() string {
	return wm.workload.GetGroup()
}

func (wm *WorkloadMock) GetGenerateName() string {
	return wm.workload.GetGenerateName()
}

func (wm *WorkloadMock) GetReplicas() int {
	return wm.workload.GetReplicas()
}

func (wm *WorkloadMock) GetKind() string {
	return wm.workload.GetKind()
}
func (wm *WorkloadMock) GetSelector() (*metav1.LabelSelector, error) {
	return wm.workload.GetSelector()
}

func (wm *WorkloadMock) GetServiceSelector() map[string]string {
	return wm.workload.GetServiceSelector()
}

func (wm *WorkloadMock) GetAnnotation(annotation string) (string, bool) {
	return wm.workload.GetAnnotation(annotation)
}
func (wm *WorkloadMock) GetLabel(label string) (string, bool) {
	return wm.workload.GetLabel(label)
}

func (wm *WorkloadMock) GetPodLabel(label string) (string, bool) {
	return wm.workload.GetPodLabel(label)
}

func (wm *WorkloadMock) GetLabels() map[string]string {
	return wm.workload.GetLabels()
}

// GetInnerLabels - DEPRECATED
func (wm *WorkloadMock) GetInnerLabels() map[string]string {
	return wm.workload.GetInnerLabels()
}

func (wm *WorkloadMock) GetPodLabels() map[string]string {
	return wm.workload.GetPodLabels()
}

// GetInnerAnnotations - DEPRECATED
func (wm *WorkloadMock) GetInnerAnnotations() map[string]string {
	return wm.workload.GetInnerAnnotations()
}

// GetPodAnnotations
func (wm *WorkloadMock) GetPodAnnotations() map[string]string {
	return wm.workload.GetPodAnnotations()
}

// GetInnerAnnotation DEPRECATED
func (wm *WorkloadMock) GetInnerAnnotation(annotation string) (string, bool) {
	return wm.workload.GetInnerAnnotation(annotation)
}

func (wm *WorkloadMock) GetPodAnnotation(annotation string) (string, bool) {
	return wm.workload.GetPodAnnotation(annotation)
}

func (wm *WorkloadMock) GetAnnotations() map[string]string {
	return wm.workload.GetAnnotations()
}

// GetVolumes -
func (wm *WorkloadMock) GetVolumes() ([]corev1.Volume, error) {
	return wm.workload.GetVolumes()
}

func (wm *WorkloadMock) GetServiceAccountName() string {
	return wm.workload.GetServiceAccountName()
}

func (wm *WorkloadMock) GetPodSpec() (*corev1.PodSpec, error) {
	return wm.workload.GetPodSpec()
}

func (wm *WorkloadMock) GetImagePullSecret() ([]corev1.LocalObjectReference, error) {
	return wm.workload.GetImagePullSecret()
}

// GetContainers -
func (wm *WorkloadMock) GetContainers() ([]corev1.Container, error) {
	return wm.workload.GetContainers()
}

// GetInitContainers -
func (wm *WorkloadMock) GetInitContainers() ([]corev1.Container, error) {
	return wm.workload.GetInitContainers()
}

// GetEphemeralContainers -
func (wm *WorkloadMock) GetEphemeralContainers() ([]corev1.EphemeralContainer, error) {
	return wm.workload.GetEphemeralContainers()
}

// GetOwnerReferences -
func (wm *WorkloadMock) GetOwnerReferences() ([]metav1.OwnerReference, error) {
	return wm.workload.GetOwnerReferences()
}
func (wm *WorkloadMock) GetResourceVersion() string {
	return wm.workload.GetResourceVersion()
}
func (wm *WorkloadMock) GetUID() string {
	return wm.workload.GetUID()
}
func (wm *WorkloadMock) GetWlid() string {
	return wm.workload.GetWlid()
}

func (wm *WorkloadMock) GenerateWlid(clusterName string) string {
	return wm.workload.GenerateWlid(clusterName)
}

func (wm *WorkloadMock) GetJobID() *apis.JobTracking {
	return wm.workload.GetJobID()
}

func (wm *WorkloadMock) GetData() map[string]interface{} {
	return wm.workload.GetData()
}

func (wm *WorkloadMock) GetSecretsOfContainer() (map[string][]string, error) {
	return wm.workload.GetSecretsOfContainer()
}

func (wm *WorkloadMock) GetConfigMapsOfContainer() (map[string][]string, error) {
	return wm.workload.GetConfigMapsOfContainer()
}

func (wm *WorkloadMock) GetSecrets() ([]string, error) {
	return wm.workload.GetSecrets()
}

func (wm *WorkloadMock) GetConfigMaps() ([]string, error) {
	return wm.workload.GetConfigMaps()
}

func (wm *WorkloadMock) GetPodStatus() (*corev1.PodStatus, error) {
	return wm.workload.GetPodStatus()
}
