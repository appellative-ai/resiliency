package operations

import (
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceName = "test:resiliency:agent/operations/host"
)

func init() {
	// Register access.Agent as it is in core and does not have access to the repository
	repository.Register(access2.Agent)
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(operations.Serve)
	})

}

type agentT struct {
	service *operations.Service
}

func newAgent(service *operations.Service) *agentT {
	a := new(agentT)
	a.service = service

	return a
}

// String - identity
func (a *agentT) String() string { return a.Name() }

// Name - agent identifier
func (a *agentT) Name() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if m.Name == messaging.ConfigEvent {
		a.configure(m)
		return
	}
	/*
		switch m.Name {
		case messaging.StartupEvent:
			host.Broadcast(m)
		case messaging.ShutdownEvent:
			host.Broadcast(m)
		case messaging.PauseEvent:
			host.Broadcast(m)
		case messaging.ResumeEvent:
			host.Broadcast(m)
		}

	*/
}

func (a *agentT) configure(m *messaging.Message) {
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}
