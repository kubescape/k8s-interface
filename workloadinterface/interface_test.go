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
