package workloadinterface

func ListMetaToMap(meta []IMetadata) []map[string]interface{} {
	resourceMap := []map[string]interface{}{}
	for i := range meta {
		resourceMap = append(resourceMap, meta[i].GetObject())
	}
	return resourceMap
}

func ListMetaIDs(meta []IMetadata) []string {
	ids := []string{}
	for i := range meta {
		ids = append(ids, meta[i].GetID())
	}
	return ids
}
