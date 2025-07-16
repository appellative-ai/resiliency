package operations

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
)

func ExampleNewAgent() {
	agent := newAgent(operationstest.NewNService())

	fmt.Printf("test: NewAgent() -> [%v]\n", agent)

	//Output:
	//test: NewAgent() -> [test:resiliency:agent/operations/host]

}
