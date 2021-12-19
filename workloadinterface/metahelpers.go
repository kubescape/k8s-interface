package workloadinterface

import "encoding/json"

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

// MapToBytes convert map[string]interface{} to []byte while ignoring errors. will return nil if failed to convert
func MapToBytes(m map[string]interface{}) []byte {
	if len(m) > 0 {
		if b, e := json.Marshal(m); e == nil {
			return b
		}
	}
	return nil
}

// BytesToMap convert []byte to map[string]interface{} while ignoring errors. will return nil if failed to convert
func BytesToMap(b []byte) map[string]interface{} {
	if len(b) > 0 {
		m := map[string]interface{}{}
		if err := json.Unmarshal(b, &m); err == nil {
			return m
		}
	}
	return nil
}
