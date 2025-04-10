package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
)

var (
	Agent = New()
)

func Initialize(notifier eventing.NotifyFunc, activity eventing.ActivityFunc) {
	if notifier != nil {
		eventing.Handler.Message(newConfigNotifier(notifier))
	}
	if activity != nil {
		//eventing.Handler.Message(newConfigActivity(activity))
	}
}

// Configure - configure all agents
func Configure(m *messaging.Message) {
	if m.Event() == messaging.ConfigEvent && m.ContentType() == messaging.ContentTypeMap {
		o, ok := newOriginFromMessage(Agent, m)
		if ok {
			access.SetOrigin(o)
		}
		exchange.Broadcast(m)
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
