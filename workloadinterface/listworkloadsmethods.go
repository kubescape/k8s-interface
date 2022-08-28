package workloadinterface

import (
	"encoding/json"
)

const TypeListWorkloads ObjectType = "List"

type ListWorkloads struct {
	listWorkloads map[string]interface{}
}

// NewListWorkloads construct a NewListWorkloadsObj from []byte. If the byte does not match the object, will return nil and err
func NewListWorkloads(bObject []byte) (*ListWorkloads, error) {
	list := make(map[string]interface{})
	if bObject != nil {
		if err := json.Unmarshal(bObject, &list); err != nil {
			return nil, err
		}
	}
	return NewListWorkloadsObj(list), nil
}

// NewListWorkloadsObj construct a NewListWorkloadsObj from map[string]interface{}. If the map does not match the object, will return nil
func NewListWorkloadsObj(object map[string]interface{}) *ListWorkloads {
	return &ListWorkloads{
		listWorkloads: object,
	}
}

// Irrelevant. A list can contain workloads from several namespaces
func (lw *ListWorkloads) GetNamespace() string {
	return ""
}

// Irrelevant. A list object has no name
func (lw *ListWorkloads) GetName() string {
	return ""
}

func (lw *ListWorkloads) GetKind() string {
	if v, ok := InspectMap(lw.listWorkloads, "kind"); ok {
		return v.(string)
	}
	return ""
}

func (lw *ListWorkloads) GetApiVersion() string {
	if v, ok := InspectMap(lw.listWorkloads, "apiVersion"); ok {
		return v.(string)
	}
	return ""
}

// Irrelevant for list obj
func (lw *ListWorkloads) GetWorkload() map[string]interface{} {
	return nil
}

func (lw *ListWorkloads) GetObject() map[string]interface{} {
	return lw.listWorkloads
}

// Irrelevant for list obj
func (lw *ListWorkloads) GetID() string {
	return ""
}

func (lw *ListWorkloads) GetObjectType() ObjectType {
	return TypeListWorkloads
}

func (lw *ListWorkloads) GetItems() []IMetadata {
	baseObjs := []IMetadata{}
	if i, ok := InspectMap(lw.GetObject(), "items"); ok && i != nil {
		if items, ok := i.([]interface{}); ok && items != nil {
			for item := range items {
				if m, ok := items[item].(map[string]interface{}); ok && m != nil {
					if o := NewBaseObject(m); o != nil {
						baseObjs = append(baseObjs, o)
					}
				}
			}
		}
	}
	return baseObjs
}

// Irrelevant for list obj
func (lw *ListWorkloads) SetNamespace(namespace string) {
}

// Irrelevant for list obj
func (lw *ListWorkloads) SetName(name string) {
}

func (lw *ListWorkloads) SetKind(kind string) {
	lw.listWorkloads["kind"] = kind
}

// Irrelevant for list obj
func (lw *ListWorkloads) SetWorkload(listWorkloads map[string]interface{}) {
}

func (lw *ListWorkloads) SetObject(listWorkloads map[string]interface{}) {
	lw.listWorkloads = listWorkloads
}

func (lw *ListWorkloads) SetApiVersion(apiVersion string) {
	lw.listWorkloads["apiVersion"] = apiVersion
}

func IsTypeListWorkloads(object map[string]interface{}) bool {
	if object == nil {
		return false
	}

	if v, ok := InspectMap(object, "kind"); ok {
		return v.(string) == string(TypeListWorkloads)
	}
	return false
}
