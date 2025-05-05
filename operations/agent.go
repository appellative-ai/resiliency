package operations

import (
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceName = "resiliency:agent/host/operations"
)

// TODO : need host name
type agentT struct {
	//running bool
	//notifier   eventing.NotifyFunc
	//activity   eventing.ActivityFunc
	//dispatcher messaging.Dispatcher
}

// New - create a new operations agent
func New() messaging.Agent {
	return newAgent()
}

func newAgent() *agentT {
	a := new(agentT)
	//a.notifier = eventing.OutputNotify
	//a.activity = eventing.OutputActivity
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
	if m.Event() == messaging.ConfigEvent {
		a.configure(m)
		return
	}
	switch m.Event() {
	case messaging.StartupEvent:
		host.Broadcast(m)
	case messaging.ShutdownEvent:
		host.Broadcast(m)
	case messaging.PauseEvent:
		host.Broadcast(m)
	case messaging.ResumeEvent:
		host.Broadcast(m)
	}
}

func (a *agentT) configure(m *messaging.Message) {
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
