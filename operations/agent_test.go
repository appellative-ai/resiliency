package operations

import (
	"errors"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func ExampleNewAgent() {
	agent := newAgent(nil)
	status := messaging.NewStatusError(http.StatusTeapot, errors.New("error"), agent.Uri())
	status.WithMessage("notify message")
	status.WithRequestId("123-request-id")
	agent.Message(event.NewNotifyMessage(status))

	//Output:
	//fail

}
