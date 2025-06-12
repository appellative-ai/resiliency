package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
)

func ExampleNewPrimaryAgent() {
	a := newPrimaryAgent(operationstest.NewService())

	fmt.Printf("test: NewPrimaryAgent() -> [%v]\n", a.Name())

	//Output:
	//test: NewPrimaryAgent() -> [test:resiliency:agent/caseOfficer/service/traffic/ingress/primary]

}
