package operations

import (
	"fmt"
)

func ExampleNewAgent() {
	a := newAgent()

	fmt.Printf("test: NewAgent() -> [%v]\n", a)

	//Output:
	//test: NewAgent() -> [core:resiliency:agent/operations/host]

}
