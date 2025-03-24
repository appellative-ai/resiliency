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
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/operations"
)

// TODO : need host name
type agentT struct {
	running bool

	notifier event.NotifyFunc
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
	switch m.Event() {
	case messaging.StartupEvent:
		a.agents.Broadcast(m)
	case messaging.ShutdownEvent:
		a.agents.Broadcast(m)
	case messaging.PauseEvent:
		a.agents.Broadcast(m)
	case messaging.ResumeEvent:
		a.agents.Broadcast(m)
	case event.NotifyEvent:
		a.notifier(event.NotifyContent(m))
	case event.DispatchEvent:
		/*
			if m.ContentType() == ContentTypeDispatch {
				a.dispatch(DispatchContent(m))
				return
			}
		*/
	case event.ActivityEvent:
		/*
			if m.ContentType() == ContentTypeActivity {
				a.addActivity(ActivityContent(m))
				return
			}
		*/
	}
}

func (a *agentT) dispatch(channel any, event1 string) {
	//a.handler.Message(event.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) shutdown() {
	a.agents.Shutdown()
}
