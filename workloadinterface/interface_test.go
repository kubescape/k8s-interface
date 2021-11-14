package workloadinterface

import (
	"fmt"
	"testing"
)

func TestIWorkload(t *testing.T) {
	var a IWorkload
	a = NewWorkloadObj(nil)
	a = NewWorkloadMock(nil)
	fmt.Printf("%v", a.GetName())
}

func TestInspectMap(t *testing.T) {
	mapexample := map[string]interface{}{"name": "a",
		"relatedObjects": struct {
			role        string
			rolebinding string
		}{"role-a", "rolebinding-b"},
		"namespace": "c",
	}
	if name, ok := InspectMap(mapexample, "name"); ok {
		if name != "a" {
			t.Errorf("error in inspectMap, name %s != 'a'", name)
		}
	} else {
		t.Errorf("error in inspectMap, couldn't find name: '%s", name)
	}
}
