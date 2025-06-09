package caseofficer

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/namespace"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
)

const (
	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "role"
	RateLimiterRole   = "rate-limiter"
	RoutingRole       = "routing"
	RedirectRole      = "redirect"
	NameKey           = "name"
)

func buildLink(handler messaging.Agent, cfg map[string]string, role string) (any, error) {
	if handler == nil || cfg == nil {
		return nil, errors.New(fmt.Sprintf("agent or configuration map is nil"))
	}
	name, ok := cfg[NameKey]
	if !ok || name == "" {
		return nil, errors.New(fmt.Sprintf("agent or exchange name not found or is empty for role: %v", role))
	}
	switch namespace.Kind(name) {
	case namespace.Link:
		l := repository.ExchangeLink(name)
		if l == nil {
			return nil, errors.New(fmt.Sprintf("exchange link is nil for name: %v and role: %v", name, role))
		}
		return l, nil
	case namespace.AgentKind:
		agent := repository.NewAgent(name)
		if agent == nil {
			return nil, errors.New(fmt.Sprintf("agent is nil for name: %v and role: %v", name, role))
		}
		// TODO: wait for reply?
		m := messaging.NewMapMessage(cfg)
		agent.Message(m)
		m = messaging.NewHandlerMessage(handler)
		agent.Message(m)
		return agent, nil
	default:
	}
	return nil, errors.New(fmt.Sprintf("invalid Namespace kind: %v and role: %v", namespace.Kind(name), role))
}
