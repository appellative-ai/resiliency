package caseofficer

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/namespace"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
)

func buildNetwork(a messaging.Agent, netCfg map[string]map[string]string) (chain []any, errs []error) {
	if a == nil {
		return nil, []error{errors.New("agent is nil")}
	}
	if len(netCfg) == 0 {
		return nil, []error{errors.New("network configuration is nil or empty")}
	}
	var router bool
	var roles = []string{LoggingRole, AuthorizationRole, CacheRole, RateLimiterRole, RoutingRole}

	for _, role := range roles {
		agentCfg, ok := netCfg[role]
		if !ok {
			continue
		}
		link, err := buildLink(role, agentCfg, a)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		chain = append(chain, link)
		if role == RoutingRole {
			router = true
		}
	}
	if len(errs) > 0 {
		return
	}
	if !router {
		errs = append(errs, errors.New("no routing agent was configured"))
	}
	return
}

func buildLink(role string, cfg map[string]string, officer messaging.Agent) (any, error) {
	name, ok := cfg[NameKey]
	if !ok || name == "" {
		return nil, errors.New(fmt.Sprintf("agent or exchange name not found or is empty for role: %v", role))
	}
	switch namespace.Kind(name) {
	case namespace.Link:
		// Since this is only code and no state, the same link can be used in all networks
		link := repository.ExchangeLink(name)
		if link == nil {
			return nil, errors.New(fmt.Sprintf("exchange link is nil for name: %v and role: %v", name, role))
		}
		return link, nil
	case namespace.AgentKind:
		// Construct a new agent as each agent has state, and a new instance is required for each network
		agent := repository.NewAgent(name)
		if agent == nil {
			return nil, errors.New(fmt.Sprintf("agent is nil for name: %v and role: %v", name, role))
		}
		// TODO: wait for reply?
		agent.Message(messaging.NewMapMessage(cfg))
		agent.Message(messaging.NewAgentMessage(officer))
		return agent, nil
	default:
	}
	return nil, errors.New(fmt.Sprintf("invalid Namespace kind: %v and role: %v", namespace.Kind(name), role))
}
