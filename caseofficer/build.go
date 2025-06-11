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
	NameKey           = "name"
)

func configureOperative(officer messaging.Agent, cfg map[string]string, role string) error {
	name, ok := cfg[NameKey]
	if !ok || name == "" {
		return errors.New(fmt.Sprintf("agent or exchange name not found or is empty for role: %v", role))
	}
	switch namespace.Kind(name) {
	case namespace.Link:
		l := repository.ExchangeLink(name)
		if l == nil {
			return errors.New(fmt.Sprintf("exchange link is nil for name: %v and role: %v", name, role))
		}
		return nil
	case namespace.AgentKind:
		agent := repository.NewAgent(name)
		if agent == nil {
			return errors.New(fmt.Sprintf("agent is nil for name: %v and role: %v", name, role))
		}
		// TODO: wait for reply?
		agent.Message(messaging.NewMapMessage(cfg))
		agent.Message(messaging.NewHandlerMessage(officer))
		return nil
	default:
	}
	return errors.New(fmt.Sprintf("invalid Namespace kind: %v and role: %v", namespace.Kind(name), role))
}
