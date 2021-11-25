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

// ToUnique removes the resource duplication based on resource ID
func ToUnique(resources []IMetadata) {
	uniqueRuleResponses := map[string]bool{}

	lenK8sResources := len(resources)
	for i := 0; i < lenK8sResources; i++ {
		resourceID := resources[i].GetID()
		if found := uniqueRuleResponses[resourceID]; found {
			// resource found -> remove from slice
			resources = removeMetadataFromSliceSlice(resources, i)
			lenK8sResources -= 1
			i -= 1
			continue
		} else {
			uniqueRuleResponses[resourceID] = true
		}
	}
}

func removeMetadataFromSliceSlice(resources []IMetadata, i int) []IMetadata {
	if i != len(resources)-1 {
		resources[i] = resources[len(resources)-1]
	}

	return resources[:len(resources)-1]
}
