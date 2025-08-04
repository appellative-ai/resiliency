package operations

import (
	"fmt"
)

func ExampleNewAgent() {
	agent := newAgent()

	fmt.Printf("test: NewAgent() -> [%v]\n", agent)

	//Output:
	//test: NewAgent() -> [test:resiliency:agent/operations/host]

}
