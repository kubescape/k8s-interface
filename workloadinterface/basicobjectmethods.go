package workloadinterface

import (
	"encoding/json"
	"fmt"
	"strings"
)

const TypeBaseObject ObjectType = "base"

type BaseObject struct {
	base map[string]interface{}
}

func NewBaseObjBytes(b []byte) (*BaseObject, error) {
	base := make(map[string]interface{})
	if b != nil {
		if err := json.Unmarshal(b, &base); err != nil {
			return nil, err
		}
	}
	if !IsBaseObject(base) {
		return nil, fmt.Errorf("invalid map[string]interface{} - expected baseObject compatibility")
	}
	return NewBaseObject(base), nil
}

func NewBaseObject(b map[string]interface{}) *BaseObject {
	return &BaseObject{
		base: b,
	}
}

func (b *BaseObject) GetObjectType() ObjectType {
	return TypeBaseObject
}

func (b *BaseObject) Json() string {
	return b.ToString()
}
func (b *BaseObject) ToString() string {
	if b.GetWorkload() == nil {
		return ""
	}
	bWorkload, err := json.Marshal(b.GetWorkload())
	if err != nil {
		return err.Error()
	}
	return string(bWorkload)
}

// ========================================= GET =========================================
func (b *BaseObject) GetWorkload() map[string]interface{} {
	return b.GetObject()
}
func (b *BaseObject) GetObject() map[string]interface{} {
	return b.base
}
func (b *BaseObject) GetNamespace() string {
	if v, ok := InspectWorkload(b.base, "metadata", "namespace"); ok {
		return v.(string)
	}
	return ""
}
func (b *BaseObject) GetID() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", b.GetGroup(), b.GetVersion(), b.GetNamespace(), b.GetKind(), b.GetName())

	// TODO - return like selfLink - e.g. /apis/apps/v1/namespaces/monitoring/statefulsets/alertmanager-prometheus-
	// return fmt.Sprintf("apps/%s/%s/%s/%s", b.GetApiVersion(), b.GetNamespace(), b.GetKind(), b.GetName())
}
func (b *BaseObject) GetName() string {
	if v, ok := InspectWorkload(b.base, "metadata", "name"); ok {
		return v.(string)
	}
	return ""
}

func (b *BaseObject) GetApiVersion() string {
	if v, ok := InspectWorkload(b.base, "apiVersion"); ok {
		return v.(string)
	}
	return ""
}

func (b *BaseObject) GetVersion() string {
	apiVersion := b.GetApiVersion()
	splitted := strings.Split(apiVersion, "/")
	if len(splitted) == 1 {
		return splitted[0]
	} else if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func (b *BaseObject) GetGroup() string {
	apiVersion := b.GetApiVersion()
	splitted := strings.Split(apiVersion, "/")
	if len(splitted) == 2 {
		return splitted[0]
	}
	return ""
}

func (b *BaseObject) GetKind() string {
	if v, ok := InspectWorkload(b.base, "kind"); ok {
		return v.(string)
	}
	return ""
}

func (b *BaseObject) GetWorkloadData() map[string]interface{} {
	if v, ok := InspectWorkload(b.base, "data"); ok {
		return v.(map[string]interface{})
	}
	return nil
}

// ========================================= SET =========================================

func (b *BaseObject) SetWorkload(workload map[string]interface{}) {
	b.SetObject(workload)
}

func (b *BaseObject) SetObject(workload map[string]interface{}) {
	b.base = workload
}

func (b *BaseObject) SetKind(kind string) {
	b.base["kind"] = kind
}

func (b *BaseObject) SetNamespace(namespace string) {
	SetInMap(b.base, []string{"metadata"}, "namespace", namespace)
}

func (b *BaseObject) SetName(name string) {
	SetInMap(b.base, []string{"metadata"}, "name", name)
}

// ===================== UTILS =======================
func IsBaseObject(b map[string]interface{}) bool {
	if kind, ok := InspectMap(b, "kind"); !ok || kind == "" {
		return false
	}
	if apiVersion, ok := InspectMap(b, "apiVersion"); !ok || apiVersion == "" {
		return false
	}
	if name, ok := InspectMap(b, "metadata", "name"); !ok || name == "" {
		return false
	}
	return true
}
