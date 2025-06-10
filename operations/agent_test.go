package operations

import "fmt"

func ExampleNewAgent() {
	agent := newAgent()
	//status := messaging.NewStatusError(http.StatusTeapot, errors.New("error"), agent.Uri())
	//status.WithMessage("notify message")
	//status.WithRequestId("123-request-id")
	//agent.Message(eventing.NewNotifyMessage(status))

	fmt.Printf("test: NewAgent() -> [%v]\n", agent)

	//Output:
	//test: NewAgent() -> [test:resiliency:agent/operations/host]

}
