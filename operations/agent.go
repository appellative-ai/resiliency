package operations

import (
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
	_ "github.com/behavioral-ai/resiliency/link"
	_ "github.com/behavioral-ai/traffic/module"
)

const (
	NamespaceName = "test:resiliency:agent/operations/host"
)

var (
	opsAgent *agentT
)

func init() {
	// Register access.Agent as it is in core and does not have access to the repository
	err := exchange.Register(access2.Agent)
	if err != nil {
		fmt.Printf("repository register error: %v", err)
	}
	//repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
	//	return newAgent(operations.Serve)
	//})
	opsAgent = newAgent(operations.Serve)
	exchange.Register(opsAgent)
}

type agentT struct {
	running bool
	service *operations.Service
	ex      *messaging.Exchange
}

func newAgent(service *operations.Service) *agentT {
	a := new(agentT)
	a.service = service
	a.ex = messaging.NewExchange()
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
	if !a.running {
		if m.Name == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Name == messaging.StartupEvent {
			a.run()
			a.running = true
			a.ex.Broadcast(m)
			return
		}
		return
	}
	list := m.To()
	// No recipient, or only the case officer recipient
	if len(list) == 0 || len(list) == 1 && list[0] == NamespaceName {
		switch m.Channel() {
		case messaging.ChannelEmissary:
			//a.emissary.C <- m
		case messaging.ChannelControl:
			//a.emissary.C <- m
		default:
			fmt.Printf("limiter - invalid channel %v\n", m)
		}
		return
	}
	a.ex.Broadcast(m)
}

// Run - run the agent
func (a *agentT) run() {
}

func (a *agentT) shutdown() {
}

func (a *agentT) registerCaseOfficer(agent messaging.Agent) {
	if agent != nil {
		a.ex.Register(agent)
	}
}

func (a *agentT) configure(m *messaging.Message) {
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}
