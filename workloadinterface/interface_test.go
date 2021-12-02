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
func TestIMetadata(t *testing.T) {
	var a IMetadata
	a = NewRegoResponseVectorObject(nil)
	a = NewWorkloadObj(nil)
	fmt.Printf("%v", a.GetName())
}
