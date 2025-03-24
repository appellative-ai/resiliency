package operations

import (
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/limit"
	"github.com/behavioral-ai/resiliency/redirect"
	"github.com/behavioral-ai/resiliency/routing"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/agency/operations"
)

// TODO : need host name
type agentT struct {
	running bool

	notifier event.NotifyFunc
	emissary *messaging.Channel
	agents   *messaging.Exchange
}

// New - create a new operations agent
func New(notifier event.NotifyFunc) messaging.Agent {
	return newAgent(notifier)
}

func newAgent(notifier event.NotifyFunc) *agentT {
	a := new(agentT)
	if notifier == nil {
		a.notifier = event.OutputNotify
	} else {
		a.notifier = notifier
	}
	a.agents = messaging.NewExchange()
	a.agents.RegisterMailbox(cache.Agent)
	a.agents.RegisterMailbox(limit.Agent)
	a.agents.RegisterMailbox(redirect.Agent)
	a.agents.RegisterMailbox(routing.Agent)
	a.emissary = messaging.NewEmissaryChannel()

	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if m.Event() == messaging.StartupEvent && !a.running {
		a.run()
		a.agents.Broadcast(m)
		return
	}
	if !a.running {
		return
	}
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) run() {
	if a.running {
		return
	}
	go emissaryAttend(a)
	a.running = true
}

func (a *agentT) dispatch(channel any, event1 string) {
	//a.handler.Message(event.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) shutdown() {
	a.emissary.Close()
	a.agents.Shutdown()
}
