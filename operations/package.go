package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/traffic/analytics"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
)

var (
	Agent = New(nil)
)

func init() {
	// intermediary agents
	cache.Initialize(Agent)
	routing.Initialize(Agent)

	// traffic agents
	analytics.Initialize(Agent)
	limiter.Initialize(Agent)
	redirect.Initialize(Agent)
}

func Initialize(notifier eventing.NotifyFunc) {
	if notifier != nil {
		Agent.Message(newConfigNotifier(notifier))
	}
}

// Configure - configure all agents
func Configure(m *messaging.Message) {
	if m.Event() == messaging.ConfigEvent && m.ContentType() == messaging.ContentTypeMap {
		o, ok := newOriginFromMessage(Agent, m)
		if ok {
			access.SetOrigin(o)
		}
		limiter.Agent.Message(m)
		redirect.Agent.Message(m)
		analytics.Agent.Message(m)
		routing.Agent.Message(m)
		cache.Agent.Message(m)
	}
}

// Message - operations agent messaging
func Message(event string) error {
	switch event {
	case messaging.StartupEvent:
		if Agent != nil {
			Agent.Message(messaging.StartupMessage)
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
