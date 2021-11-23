package workloadinterface

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	role        = `{"apiVersion": "rbac.authorization.k8s.io/v1", "kind": "Role", "metadata": {"creationTimestamp": "2021-06-13T13:17:24Z","managedFields": [{"apiVersion": "rbac.authorization.k8s.io/v1","fieldsType": "FieldsV1","fieldsV1": {"f:rules": {}},"manager": "kubectl-edit","operation": "Update","time": "2021-06-13T13:22:29Z"}],"name": "pod-reader","namespace": "default","resourceVersion": "40233","uid": "cea4a847-2f05-4a94-bf3f-a8d1907e60e0"},"rules": [{"apiGroups": [""],"resources": ["pods","secrets"],"verbs": ["get"]}]}`
	rolebinding = `{"apiVersion":"rbac.authorization.k8s.io/v1", "kind":"RoleBinding", "metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"rbac.authorization.k8s.io/v1\",\"kind\":\"RoleBinding\",\"metadata\":{\"annotations\":{},\"name\":\"read-pods\",\"namespace\":\"default\"},\"roleRef\":{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"Role\",\"name\":\"pod-reader\"},\"subjects\":[{\"apiGroup\":\"rbac.authorization.k8s.io\",\"kind\":\"User\",\"name\":\"jane\"}]}\n"},"creationTimestamp":"2021-11-11T11:50:38Z","managedFields":[{"apiVersion":"rbac.authorization.k8s.io/v1","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:roleRef":{"f:apiGroup":{},"f:kind":{},"f:name":{}},"f:subjects":{}},"manager":"kubectl-client-side-apply","operation":"Update","time":"2021-11-11T11:50:38Z"}],"name":"read-pods","namespace":"default","resourceVersion":"650451","uid":"6038eca8-b13e-4557-bc92-8800a11197d3"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"Role","name":"pod-reader"},"subjects":[{"apiGroup":"rbac.authorization.k8s.io","kind":"User","name":"jane"}]}`
)

func getMock(r string) map[string]interface{} {
	relatedObject, err := NewWorkload([]byte(r))
	if err != nil {
		panic(err)
	}
	return relatedObject.GetObject()
}
func assertObjectFields(t *testing.T, b IMetadata) {
	assert.Equal(t, "", b.GetApiVersion())
	assert.Equal(t, "user@example.com", b.GetName())
	assert.Equal(t, "User", b.GetKind())
	assert.Equal(t, "default", b.GetNamespace())
	assert.Equal(t, "/default/User/user@example.com/rbac.authorization.k8s.io/v1/default/Role/pod-reader/rbac.authorization.k8s.io/v1/default/RoleBinding/read-pods", b.GetID())
}
func TestNewRegoResponseVectorObject(t *testing.T) {
	relatedObjects := []map[string]interface{}{}
	relatedObject := getMock(role)
	relatedObject2 := getMock(rolebinding)
	relatedObjects = append(relatedObjects, relatedObject)
	relatedObjects = append(relatedObjects, relatedObject2)
	subject := map[string]interface{}{"name": "user@example.com", "kind": "User", "namespace": "default", "group": "rbac.authorization.k8s.io", RelatedObjectsKey: relatedObjects}
	assert.True(t, IsTypeRegoResponseVector(subject))

	obj := NewRegoResponseVectorObject(subject)
	assert.Equal(t, 2, len(obj.GetRelatedObjects()))
	assertObjectFields(t, obj)

	respVector, err := NewRegoResponseVectorObjectFromBytes([]byte(obj.ToString()))
	assert.NoError(t, err)
	assertObjectFields(t, respVector)
}

func TestSetGetObject(t *testing.T) {
	obj := `{"name":"Jane","namespace":"","kind":"User","apiVersion":""}`
	relatedObjects := []map[string]interface{}{}
	relatedObject := getMock(role)
	relatedObject2 := getMock(rolebinding)
	relatedObjects = append(relatedObjects, relatedObject)
	relatedObjects = append(relatedObjects, relatedObject2)
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(obj), &m)
	if err != nil {
		t.Errorf("error unmarshaling, %s", err.Error())
	}
	respVector2 := NewRegoResponseVectorObject(nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	respVector2.SetObject(m)
	respVector2.SetRelatedObjects(relatedObjects)
	if respVector2.GetID() == "" {
		t.Errorf("error setting object")
	}
	object := respVector2.GetObject()
	if len(object) == 0 {
		t.Errorf("error getting object")
	}
	if len(object["relatedObjects"].([]map[string]interface{})) == 0 {
		t.Errorf("error getting object")
	}
}
