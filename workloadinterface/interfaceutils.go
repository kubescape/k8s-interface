package workloadinterface

const (
	TypeUnknown ObjectType = "unknown"
)

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

func SetInMap(workload map[string]interface{}, scope []string, key string, val interface{}) {
	for i := range scope {
		if _, ok := workload[scope[i]]; !ok {
			workload[scope[i]] = make(map[string]interface{})
		}
		workload, _ = workload[scope[i]].(map[string]interface{})
	}

	workload[key] = val
}

func RemoveFromMap(workload map[string]interface{}, scope ...string) {
	for i := 0; i < len(scope)-1; i++ {
		if _, ok := workload[scope[i]]; ok {
			workload, _ = workload[scope[i]].(map[string]interface{})
		}
	}
	delete(workload, scope[len(scope)-1])
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
