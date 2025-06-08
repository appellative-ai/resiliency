package caseofficer

import (
	"errors"
	"fmt"
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

func buildAgent(handler messaging.Agent, cfg map[string]string, role string) (messaging.Agent, error) {
	name, ok := cfg[NameKey]
	if !ok || name == "" {
		return nil, errors.New(fmt.Sprintf("agent name not found or is empty for role %v", role))
	}
	agent := repository.Constructor(name)
	if agent == nil {
		return nil, errors.New(fmt.Sprintf("invalid agent name %v for role %v", name, role))
	}
	// TODO: wait for reply?
	m := messaging.NewMapMessage(cfg)
	agent.Message(m)
	m = messaging.NewHandlerMessage(handler)
	agent.Message(m)
	return agent, nil
}
