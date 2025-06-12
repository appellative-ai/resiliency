package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceNamePrimary = "test:resiliency:agent/caseOfficer/service/traffic/ingress/primary"
	NetworkNamePrimary   = "test:resiliency:network/service/traffic/ingress/primary"

	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "cache"
	RateLimiterRole   = "rate-limiter"
	RoutingRole       = "routing"
	NameKey           = "name"
)

// TODO : need host name
type primaryAgentT struct {
	running bool
	handler messaging.Agent

	cfg map[string]map[string]string

	ex       *messaging.Exchange
	emissary *messaging.Channel
	service  *operations.Service

	//dispatcher messaging.Dispatcher
}

// NewPrimaryAgent - create a new agent
func NewPrimaryAgent(service *operations.Service) Agent {
	return newAgent(service)
}

func newAgent(service *operations.Service) *primaryAgentT {
	a := new(primaryAgentT)
	//a.cfg = cfg

	a.service = service
	a.ex = messaging.NewExchange()
	a.emissary = messaging.NewEmissaryChannel()
	return a
}

// Name - agent identifier
func (a *primaryAgentT) Name() string { return NamespaceNamePrimary }

// Message - message the agent
func (a *primaryAgentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if !a.running {
		if m.Name == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Name == messaging.StartupEvent {
			a.run()
			a.running = true
			return
		}
		return
	}
	if m.Name == messaging.ShutdownEvent {
		a.running = false
	}
	list := m.To()
	if len(list) == 0 {
		// Need to create some sort of error
		return
	}
	// If to is the current case officer, then send to channel
	if list[0] == NetworkNamePrimary {
		switch m.Channel() {
		case messaging.ChannelEmissary:
			a.emissary.C <- m
		case messaging.ChannelControl:
			a.emissary.C <- m
		default:
			fmt.Printf("limiter - invalid channel %v\n", m)
		}
		return
	}
	// Send to appropriate agent
	a.ex.Message(m)
}

func (a *primaryAgentT) BuildNetwork(net map[string]map[string]string) (chain []any, errs []error) {
	return buildNetwork(a, net)
	/*
		if net == nil {
			return nil, []error{errors.New("error: configuration nil")}
		}
		var router bool
		var roles = []string{LoggingRole, AuthorizationRole, CacheRole, RateLimiterRole, RoutingRole}

		for _, role := range roles {
			cfg, ok := net[role]
			if !ok {
				continue
			}
			link, err := buildLink(role, cfg, a)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			chain = append(chain, link)
		}
		if len(errs) > 0 {
			return
		}
		if !router {
			errs = append(errs, errors.New("error: no routing agent was configured"))
		}
		return

	*/
}

// Run - run the agent
func (a *primaryAgentT) run() {
	// TODO: initialize network
	go primaryEmissaryAttend(a)
}

func (a *primaryAgentT) shutdown() {
	a.ex.Broadcast(messaging.ShutdownMessage)
	a.emissary.Close()
}

func (a *primaryAgentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeAgent:
		h, status := messaging.AgentContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		a.handler = h
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}
