package workloadinterface

import (
	"encoding/json"
	"testing"
)

var (
	role        = `{"apiVersion": "rbac.authorization.k8s.io/v1","kind": "Role","metadata": {"creationTimestamp": "2021-06-13T13:17:24Z","managedFields": [{"apiVersion": "rbac.authorization.k8s.io/v1","fieldsType": "FieldsV1","fieldsV1": {"f:rules": {}},"manager": "kubectl-edit","operation": "Update","time": "2021-06-13T13:22:29Z"}],"name": "pod-reader","namespace": "default","resourceVersion": "40233","uid": "cea4a847-2f05-4a94-bf3f-a8d1907e60e0"},"rules": [{"apiGroups": [""],"resources": ["pods","secrets"],"verbs": ["get"]}]}`
	rolebinding = `{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"RoleBinding","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"RoleBinding\",\"metadata\":{\"annotations\":{},\"name\":\"read-pods\",\"namespace\":\"default\"},\"roleRef\":{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"Role\",\"name\":\"pod-reader\"},\"subjects\":[{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"User\",\"name\":\"jane\"}]}\n"},"creationTimestamp":"2021-11-11T11:50:38Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:roleRef":{"f:apiGroup":{},"f:kind":{},"f:name":{}},"f:subjects":{}},"manager":"kubectl-client-side-apply","operation":"Update","time":"2021-11-11T11:50:38Z"}],"name":"read-pods","namespace":"default","resourceVersion":"650451","uid":"6038eca8-b13e-4557-bc92-8800a11197d3"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"Role","name":"pod-reader"},"subjects":[{"apiGroup":"rbac.authorization.k8s.io","kind":"User","name":"jane"}]}`
)

func TestRegoResponseVectorObjectFromBytes(t *testing.T) {
	// w := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"3"},"creationTimestamp":"2021-06-21T04:52:05Z","generation":3,"name":"emailservice","namespace":"default"},"spec":{"progressDeadlineSeconds":600,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"app":"emailservice"}},"strategy":{"rollingUpdate":{"maxSurge":"25%","maxUnavailable":"25%"},"type":"RollingUpdate"},"template":{"metadata":{"annotations":{"armo.last-update":"21-06-2021 06:40:42","armo.wlid":"wlid://cluster-david-demo/namespace-default/deployment-emailservice"},"creationTimestamp":null,"labels":{"app":"emailservice","armo.attach":"true"}},"spec":{"containers":[{"env":[{"name":"PORT","value":"8080"},{"name":"DISABLE_PROFILER","value":"1"}],"image":"gcr.io/google-samples/microservices-demo/emailservice:v0.2.3","imagePullPolicy":"IfNotPresent","livenessProbe":{"exec":{"command":["/bin/grpc_health_probe","-addr=:8080"]},"failureThreshold":3,"periodSeconds":5,"successThreshold":1,"timeoutSeconds":1},"name":"server","ports":[{"containerPort":8080,"protocol":"TCP"}],"readinessProbe":{"exec":{"command":["/bin/grpc_health_probe","-addr=:8080"]},"failureThreshold":3,"periodSeconds":5,"successThreshold":1,"timeoutSeconds":1},"resources":{"limits":{"cpu":"200m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"64Mi"}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File"}],"dnsPolicy":"ClusterFirst","restartPolicy":"Always","schedulerName":"default-scheduler","securityContext":{},"serviceAccount":"default","serviceAccountName":"default","terminationGracePeriodSeconds":5}}}}`
	relatedObjects := []IMetadata{}
	relatedObject, _ := NewWorkload([]byte(role))
	relatedObject2, _ := NewWorkload([]byte(rolebinding))
	relatedObjects = append(relatedObjects, relatedObject)
	relatedObjects = append(relatedObjects, relatedObject2)
	obj := `{"name":"Jane","namespace":"","kind":"User","apiVersion":"","relatedObjects":[{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"RoleBinding","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"RoleBinding\",\"metadata\":{\"annotations\":{},\"name\":\"read-pods\",\"namespace\":\"default\"},\"roleRef\":{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"Role\",\"name\":\"pod-reader\"},\"subjects\":[{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"User\",\"name\":\"jane\"}]}\n"},"creationTimestamp":"2021-11-11T11:50:38Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:roleRef":{"f:apiGroup":{},"f:kind":{},"f:name":{}},"f:subjects":{}},"manager":"kubectl-client-side-apply","operation":"Update","time":"2021-11-11T11:50:38Z"}],"name":"read-pods","namespace":"default","resourceVersion":"650451","uid":"6038eca8-b13e-4557-bc92-8800a11197d3"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"Role","name":"pod-reader"},"subjects":[{"apiGroup":"rbac.authorization.k8s.io","kind":"User","name":"jane"}]},{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"Role","metadata":{"creationTimestamp":"2021-06-13T13:17:24Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:rules":{}},"manager":"kubectl-edit","operation":"Update","time":"2021-06-13T13:22:29Z"}],"name":"pod-reader","namespace":"default","resourceVersion":"40233","uid":"cea4a847-2f05-4a94-bf3f-a8d1907e60e0"},"rules":[{"apiGroups":[""],"resources":["pods","secrets"],"verbs":["get"]}]}]}`
	respVector, err := NewRegoResponseVectorObjectFromBytes([]byte(obj), relatedObjects)
	if err != nil {
		t.Errorf(err.Error())
	}
	if respVector.GetApiVersion() != "" {
		t.Errorf("error getting apiVersion, got: '%s', should be: '%s'", respVector.GetApiVersion(), "")
	}
	if respVector.GetName() != "Jane" {
		t.Errorf("error getting name, got: '%s', should be: '%s'", respVector.GetName(), "Jane")
	}
	if respVector.GetKind() != "User" {
		t.Errorf("error getting kind, got: '%s', should be: '%s'", respVector.GetKind(), "User")
	}
	if respVector.GetNamespace() != "" {
		t.Errorf("error getting namespace, got: '%s', should be: '%s'", respVector.GetNamespace(), "")
	}
	id := respVector.GetID()
	if id != "//User/Jane/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods" {
		t.Errorf("error getting kind, got: '%s', should be: '%s'", id, "//User/Jane/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods")
	}
}

func TestNewRegoResponseVectorObject(t *testing.T) {
	var b IMetadata
	relatedObjects := []IMetadata{}
	relatedObject, _ := NewWorkload([]byte(role))
	relatedObject2, _ := NewWorkload([]byte(rolebinding))
	relatedObjects = append(relatedObjects, relatedObject)
	relatedObjects = append(relatedObjects, relatedObject2)
	subject := map[string]interface{}{"name": "user@example.com", "kind": "User", "namespace": "default", "group": "rbac.authorization.k8s.io"}
	b = NewRegoResponseVectorObject(subject, relatedObjects)
	if b.GetApiVersion() != "" {
		t.Errorf("error getting apiVersion, got: '%s', should be: '%s'", b.GetApiVersion(), "rbac.authorization.k8s.io")
	}
	if b.GetName() != "user@example.com" {
		t.Errorf("error getting name, got: '%s', should be: '%s'", b.GetName(), "user@example.com")
	}
	if b.GetKind() != "User" {
		t.Errorf("error getting kind, got: '%s', should be: '%s'", b.GetKind(), "User")
	}
	if b.GetNamespace() != "default" {
		t.Errorf("error getting namespace, got: '%s', should be: '%s'", b.GetNamespace(), "default")
	}
	id := b.GetID()
	if id != "/default/User/user@example.com/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods" {
		t.Errorf("error getting kind, got: '%s', should be: '%s'", id, "/default/User/user@example.com/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods")
	}
}

func TestSetters(t *testing.T) {
	// obj := `{"name":"Jane","namespace":"","kind":"User","apiVersion":"","relatedObjects":[{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"RoleBinding","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"RoleBinding\",\"metadata\":{\"annotations\":{},\"name\":\"read-pods\",\"namespace\":\"default\"},\"roleRef\":{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"Role\",\"name\":\"pod-reader\"},\"subjects\":[{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"User\",\"name\":\"jane\"}]}\n"},"creationTimestamp":"2021-11-11T11:50:38Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:roleRef":{"f:apiGroup":{},"f:kind":{},"f:name":{}},"f:subjects":{}},"manager":"kubectl-client-side-apply","operation":"Update","time":"2021-11-11T11:50:38Z"}],"name":"read-pods","namespace":"default","resourceVersion":"650451","uid":"6038eca8-b13e-4557-bc92-8800a11197d3"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"Role","name":"pod-reader"},"subjects":[{"apiGroup":"rbac.authorization.k8s.io","kind":"User","name":"jane"}]},{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"Role","metadata":{"creationTimestamp":"2021-06-13T13:17:24Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:rules":{}},"manager":"kubectl-edit","operation":"Update","time":"2021-06-13T13:22:29Z"}],"name":"pod-reader","namespace":"default","resourceVersion":"40233","uid":"cea4a847-2f05-4a94-bf3f-a8d1907e60e0"},"rules":[{"apiGroups":[""],"resources":["pods","secrets"],"verbs":["get"]}]}]}`
	obj := `{"name":"Jane","namespace":"","kind":"User","apiVersion":""}`
	relatedObjects := []IMetadata{}
	relatedObject, _ := NewWorkload([]byte(role))
	relatedObject2, _ := NewWorkload([]byte(rolebinding))
	relatedObjects = append(relatedObjects, relatedObject)
	relatedObjects = append(relatedObjects, relatedObject2)
	respVector, err := NewRegoResponseVectorObjectFromBytes([]byte(obj), relatedObjects)
	if err != nil {
		t.Errorf(err.Error())
	}
	if respVector.GetApiVersion() != "" {
		t.Errorf("error getting apiVersion, got: '%s', should be: '%s'", respVector.GetApiVersion(), "")
	}
	respVector.SetName("Yossi")
	if respVector.GetName() != "Yossi" {
		t.Errorf("error getting name, got: '%s', should be: '%s'", respVector.GetName(), "Yossi")
	}
	respVector.SetKind("Group")
	if respVector.GetKind() != "Group" {
		t.Errorf("error getting kind, got: '%s', should be: '%s'", respVector.GetKind(), "Group")
	}
	respVector.SetNamespace("default")
	if respVector.GetNamespace() != "default" {
		t.Errorf("error getting namespace, got: '%s', should be: '%s'", respVector.GetNamespace(), "default")
	}
	id := respVector.GetID()
	if id != "/default/Group/Yossi/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods" {
		t.Errorf("error getting id, got: '%s', should be: '%s'", id, "/default/Group/Yossi/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods")
	}
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(obj), &m)
	if err != nil {
		t.Errorf("error unmarshaling, %s", err.Error())
	}
	respVector2 := NewRegoResponseVectorObject(nil, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	respVector2.SetObject(m)
	if respVector2.GetID() == "" {
		t.Errorf("error setting object")
	}
}
