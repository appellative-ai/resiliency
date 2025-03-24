package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/common"
	"github.com/behavioral-ai/resiliency/limit"
	"github.com/behavioral-ai/resiliency/redirect"
	"github.com/behavioral-ai/resiliency/routing"
)

var (
	Agent messaging.Agent
)

func Initialize(notifier eventing.NotifyFunc) {
	Agent = New(notifier)
}

// Configure - configure all agents
func Configure(m *messaging.Message) {
	if m.Event() == messaging.ConfigEvent && m.ContentType() == messaging.ContentTypeMap {
		o, ok := common.SetOrigin(Agent, m)
		if ok {
			access.SetOrigin(access.Origin{
				Region:     o.Region,
				Zone:       o.Zone,
				SubZone:    o.SubZone,
				Host:       o.Host,
				Route:      "",
				InstanceId: o.InstanceId,
			})
		}
		limit.Agent.Message(m)
		redirect.Agent.Message(m)
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
