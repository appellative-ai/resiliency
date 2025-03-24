package operations

import (
	"errors"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func ExampleNewAgent() {
	agent := newAgent(nil, nil)
	status := messaging.NewStatusError(http.StatusTeapot, errors.New("error"), agent.Uri())
	status.WithMessage("notify message")
	status.WithRequestId("123-request-id")
	agent.Message(eventing.NewNotifyMessage(status))

	//Output:
	//notify-> 2025-03-24T16:40:48.869Z [resiliency:agent/behavioral-ai/resiliency/operations] [core:messaging.status] [123-request-id] [I'm A Teapot] [error]

}
