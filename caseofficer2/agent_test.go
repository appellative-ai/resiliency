package caseofficer2

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
)

func ExampleNewPrimaryAgent() {
	a := newAgent("core:common:agent/caseofficer/request/http/test", operationstest.NewService())

	fmt.Printf("test: NewPrimaryAgent() -> [%v]\n", a.Name())

	//Output:
	//test: NewPrimaryAgent() -> [core:common:agent/caseofficer/request/http/test]

}
