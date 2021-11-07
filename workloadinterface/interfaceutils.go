package workloadinterface

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
