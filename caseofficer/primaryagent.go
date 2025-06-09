package caseofficer

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceNamePrimary = "test:resiliency:agent/caseOfficer/service/traffic/ingress/primary"
	NetworkNamePrimary   = "test:resiliency:network/service/traffic/ingress/primary"
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

// NewPrimary - create a new agent
func NewPrimary(cfg map[string]map[string]string, service *operations.Service) Agent {
	return newAgent(cfg, service)
}

func newAgent(cfg map[string]map[string]string, service *operations.Service) *primaryAgentT {
	a := new(primaryAgentT)
	a.cfg = cfg

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
	switch m.Channel() {
	case messaging.ChannelEmissary:
		a.emissary.C <- m
	case messaging.ChannelControl:
		a.emissary.C <- m
	default:
		fmt.Printf("limiter - invalid channel %v\n", m)
	}
}

func (a *primaryAgentT) BuildNetwork(net map[string]map[string]string) (links []any, errs []error) {
	if net == nil {
		return nil, []error{errors.New("error: configuration nil")}
	}
	var router bool
	var roles = []string{LoggingRole, AuthorizationRole, RedirectRole, CacheRole, RateLimiterRole, RoutingRole}

	for _, role := range roles {
		cfg, ok := net[role]
		if !ok {
			continue
		}
		link, err := buildLink(a, cfg, role)
		if err != nil {
			errs = append(errs, err)
		} else {
			links = append(links, link)
		}
	}
	if len(errs) > 0 {
		return
	}
	if !router {
		return nil, []error{errors.New("error: no routing agent was configured")}
	}
	return links, nil
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
	case messaging.ContentTypeHandler:
		h, status := messaging.HandlerContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		a.handler = h
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}
