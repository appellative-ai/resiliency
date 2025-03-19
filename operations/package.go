package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func Exchange(r *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

var (
	Agent messaging.Agent
)

func Message(event string) error {
	switch event {
	case messaging.StartupEvent:
		if Agent == nil {
			Agent = New()
			Agent.Run()
		}
	case messaging.ShutdownEvent:
		if Agent != nil {
			Agent.Message(messaging.ShutdownMessage)
			Agent = nil
		}
	case messaging.PauseEvent:
		if Agent != nil {
			Agent.Message(messaging.PauseMessage)
		}
	case messaging.ResumeEvent:
		if Agent != nil {
			Agent.Message(messaging.ResumeMessage)
		}
	default:
		return errors.New(fmt.Sprintf("operations.Message() -> [%v] [%v]", "error: invalid event", event))
	}
	return nil
}
