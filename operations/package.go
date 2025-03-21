package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/limit"
	"github.com/behavioral-ai/resiliency/redirect"
	"github.com/behavioral-ai/resiliency/routing"
)

var (
	Agent messaging.Agent
)

// Agent configuration

func ConfigureLimitAgent(m *messaging.Message) {
	message(limit.Agent, m)
}

func ConfigureRoutingAgent(m *messaging.Message) {
	message(routing.Agent, m)
}

func ConfigureCacheAgent(m *messaging.Message) {
	message(cache.Agent, m)
}

func ConfigureRedirectAgent(m *messaging.Message) {
	message(redirect.Agent, m)
}

func message(agent messaging.Agent, m *messaging.Message) {
	if m.Event() == messaging.ConfigEvent && m.ContentType() == messaging.ContentTypeMap {
		agent.Message(m)
	}
}

// Message - operations agent messaging
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
