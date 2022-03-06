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
func TestWorkloadWithIMetadata(t *testing.T) {
	var a IMetadata
	a = NewWorkloadObj(nil)
	fmt.Printf("%v", a.GetName())
}

func TestBasicObject(t *testing.T) {
	var a IMetadata
	a = NewBaseObject(nil)
	fmt.Printf("%v", a.GetName())
}
