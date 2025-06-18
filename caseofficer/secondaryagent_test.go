package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
)

func ExampleNewSecondaryAgent() {
	a := newSecondaryAgent(operationstest.NewService())

	fmt.Printf("test: NewSecondaryAgent() -> [%v]\n", a.Name())

	//Output:
	//test: NewSecondaryAgent() -> [test:resiliency:agent/caseOfficer/service/traffic/ingress/secondary]

}
