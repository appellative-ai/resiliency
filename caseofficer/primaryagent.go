package caseofficer

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	link "github.com/behavioral-ai/resiliency/link"
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

func (a *primaryAgentT) Startup(net map[string]map[string]string) (e *rest.Endpoint, errs []error) {
	if net == nil {
		return nil, []error{errors.New("error: configuration nil")}
	}
	var links []any
	var router bool

	if cfg, ok := net[LoggingRole]; ok {
		if _, ok1 := cfg[NameKey]; ok1 {
			links = append(links, link.Logger)
		}
	}
	if cfg, ok := net[AuthorizationRole]; ok {
		if _, ok1 := cfg[NameKey]; ok1 {
			links = append(links, link.Authorization)
		}
	}
	if cfg, ok := net[RedirectRole]; ok {
		if name, ok1 := cfg[NameKey]; ok1 {
			agent := repository.Agent(name)
			if agent == nil {
				errs = append(errs, errors.New(fmt.Sprintf("invalid agent name %v for role %v", name, RedirectRole)))
			} else {
				links = append(links, agent)
			}
		}
	}
	if cfg, ok := net[CacheRole]; ok {
		if name, ok1 := cfg[NameKey]; ok1 {
			agent := repository.Agent(name)
			if agent == nil {
				errs = append(errs, errors.New(fmt.Sprintf("invalid agent name %v for role %v", name, CacheRole)))
			} else {
				links = append(links, agent)
			}
		}
	}
	if cfg, ok := net[RateLimiterRole]; ok {
		if name, ok1 := cfg[NameKey]; ok1 {
			agent := repository.Agent(name)
			if agent == nil {
				errs = append(errs, errors.New(fmt.Sprintf("invalid agent name %v for role %v", name, RateLimiterRole)))
			} else {
				links = append(links, agent)
			}
		}
	}
	if cfg, ok := net[RoutingRole]; ok {
		if name, ok1 := cfg[NameKey]; ok1 {
			agent := repository.Agent(name)
			if agent == nil {
				errs = append(errs, errors.New(fmt.Sprintf("invalid agent name %v for role %v", name, RoutingRole)))
			} else {
				router = true
				links = append(links, agent)
			}
		}
	}
	if len(errs) > 0 {
		return
	}
	if !router {
		return nil, []error{errors.New("error: no routing agent was configured")}
	}
	return host.NewEndpoint(links), nil
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
