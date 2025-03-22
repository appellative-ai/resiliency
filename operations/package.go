package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/common"
	"github.com/behavioral-ai/resiliency/limit"
	"github.com/behavioral-ai/resiliency/redirect"
	"github.com/behavioral-ai/resiliency/routing"
)

const (
	TimeoutKey    = "timeout"
	CacheHostKey  = cache.CacheHostKey
	AppHostKey    = routing.AppHostKey
	RegionKey     = common.RegionKey
	ZoneKey       = common.ZoneKey
	SubZoneKey    = common.SubZoneKey
	HostKey       = common.HostKey
	InstanceIdKey = common.InstanceIdKey
)

var (
	Agent messaging.Agent
)

// Agent configuration

// Configure - configure all agents
func Configure(m *messaging.Message) {
	if m.Event() == messaging.ConfigEvent && m.ContentType() == messaging.ContentTypeMap {
		access.SetOrigin(configure(m))
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
		if Agent == nil {
			Agent = New()
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
