package invoke

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
	opsAgent messaging.Agent
)

func Message(event string) error {
	switch event {
	case messaging.StartupEvent:
		if opsAgent == nil {
			opsAgent = New()
			opsAgent.Run()
		}
	case messaging.ShutdownEvent:
		if opsAgent != nil {
			opsAgent.Shutdown()
			opsAgent = nil
		}
	case messaging.PauseEvent:
		if opsAgent != nil {
			opsAgent.Message(messaging.Pause)
		}
	case messaging.ResumeEvent:
		if opsAgent != nil {
			opsAgent.Message(messaging.Resume)
		}
	default:
		return errors.New(fmt.Sprintf("operative1.Message() -> [%v] [%v]", "error: invalid event", event))
	}
	return nil
}
