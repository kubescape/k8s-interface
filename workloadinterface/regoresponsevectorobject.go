package workloadinterface

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// This object is the structure of externalObject responses from regos
// that are ru as part of kubescape/ posture scan
// i.e any object that isn't an exact k8s object
// such as - subjects.

// expected fields:
// name string, namespace string, kind string, apiVersion string (can be empty str if not relevant)
// relatedObjects []IMetadata - includes related objects that need to be shown together with failed object
// e.g subjects will have in relatedObjects - role + rolebinding

type RegoResponseVectorObject struct {
	Object         map[string]interface{} `json:"object,omitempty"`
	RelatedObjects []IMetadata            `json:"relatedObjects"`
}

func NewRegoResponseVectorObject(object map[string]interface{}, relatedObjects []IMetadata) *RegoResponseVectorObject {
	return &RegoResponseVectorObject{
		Object:         object,
		RelatedObjects: relatedObjects,
	}
}

func NewRegoResponseVectorObjectFromBytes(object []byte, relatedObjects []IMetadata) (*RegoResponseVectorObject, error) {
	obj := make(map[string]interface{})
	if object != nil {
		if err := json.Unmarshal(object, &obj); err != nil {
			return nil, err
		}
	}
	delete(obj, "relatedObjects")
	return &RegoResponseVectorObject{
		Object:         obj,
		RelatedObjects: relatedObjects,
	}, nil
}

// =================== Set ================================
func (obj *RegoResponseVectorObject) SetNamespace(namespace string) {
	obj.Object["namespace"] = namespace
}

func (obj *RegoResponseVectorObject) SetName(name string) {
	obj.Object["name"] = name
}

func (obj *RegoResponseVectorObject) SetKind(kind string) {
	obj.Object["kind"] = kind
}

func (obj *RegoResponseVectorObject) SetWorkload(object map[string]interface{}) { // DEPRECATED
	obj.SetObject(object)
}

func (obj *RegoResponseVectorObject) SetObject(object map[string]interface{}) {
	obj.Object = object
}

func (obj *RegoResponseVectorObject) SetRelatedObjects(relatedObjects []IMetadata) {
	obj.RelatedObjects = relatedObjects
}

// =================== Get ================================
func (obj *RegoResponseVectorObject) GetApiVersion() string {
	if v, ok := InspectMap(obj.Object, "apiVersion"); ok {
		return v.(string)
	}
	return ""
}

func (obj *RegoResponseVectorObject) GetNamespace() string {
	if v, ok := InspectMap(obj.Object, "namespace"); ok {
		return v.(string)
	}
	return ""
}

func (obj *RegoResponseVectorObject) GetName() string {
	if v, ok := InspectMap(obj.Object, "name"); ok {
		return v.(string)
	}
	return ""
}

func (obj *RegoResponseVectorObject) GetKind() string {
	if v, ok := InspectMap(obj.Object, "kind"); ok {
		return v.(string)
	}
	return ""
}

func (obj *RegoResponseVectorObject) GetWorkload() map[string]interface{} { // DEPRECATED
	return obj.GetObject()
}

func (obj *RegoResponseVectorObject) GetObject() map[string]interface{} {
	var object map[string]interface{}
	if temp, err := json.Marshal(obj.Object); err == nil {
		json.Unmarshal(temp, &object)
	}
	object["relatedObjects"] = []map[string]interface{}{}
	for _, relatedobj := range obj.RelatedObjects {
		ro := relatedobj.GetObject()
		object["relatedObjects"] = append(object["relatedObjects"].([]map[string]interface{}), ro)
	}
	return object
}

func (obj *RegoResponseVectorObject) GetRelatedObjects() []IMetadata {
	return obj.RelatedObjects
}

func (obj *RegoResponseVectorObject) GetID() string {
	relatedObjectsIDs := []string{}
	for _, o := range obj.RelatedObjects {
		if o != nil {
			relatedObjectsIDs = append(relatedObjectsIDs, o.GetID())
		}
	}
	relatedObjectsIDs = append(relatedObjectsIDs, fmt.Sprintf("%s/%s/%s/%s", obj.GetApiVersion(), obj.GetNamespace(), obj.GetKind(), obj.GetName()))
	sort.Strings(relatedObjectsIDs)
	return strings.Join(relatedObjectsIDs, "/")
}
