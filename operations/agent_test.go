package operations

import (
	"fmt"
	"github.com/appellative-ai/collective/operations/operationstest"
)

func ExampleNewAgent() {
	agent := newAgent(operationstest.NewNotifier())

	fmt.Printf("test: NewAgent() -> [%v]\n", agent)

	//Output:
	//test: NewAgent() -> [test:resiliency:agent/operations/host]

}
