package workloadinterface

type ObjectType string

const (
	TypeWorkloadObject           ObjectType = "workload"
	TypeRegoResponseVectorObject ObjectType = "regoResponse"
	TypeUnknown                  ObjectType = "unknown"
)

// Returns the currect object that supports the IMetadata interface
func NewObject(object map[string]interface{}) IMetadata {
	if object == nil {
		return nil
	}
	switch GetObjectType(object) {
	case TypeWorkloadObject:
		return NewWorkloadObj(object)
	case TypeRegoResponseVectorObject:
		return NewRegoResponseVectorObject(object)
	}
	return nil
}

func GetObjectType(object map[string]interface{}) ObjectType {
	if IsTypeWorkload(object) {
		return TypeWorkloadObject
	}
	if IsTypeRegoResponseVector(object) {
		return TypeRegoResponseVectorObject
	}
	return TypeUnknown
}

// InspectMap -
func InspectMap(mapobject interface{}, scopes ...string) (val interface{}, k bool) {

	val, k = nil, false
	if len(scopes) == 0 {
		if mapobject != nil {
			return mapobject, true
		}
		return nil, false
	}
	if data, ok := mapobject.(map[string]interface{}); ok {
		val, k = InspectMap(data[scopes[0]], scopes[1:]...)
	}
	return val, k

}
