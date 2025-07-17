package operations

import (
	"fmt"
	"github.com/appellative-ai/agency/caseofficer"
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/collective/operations"

	"github.com/appellative-ai/core/messaging"

	_ "github.com/appellative-ai/traffic/module"
)

const (
	NamespaceName      = "test:resiliency:agent/operations/host"
	caseOfficerNameFmt = "core:common:agent/caseofficer/request/http/%v"
)

type Agent interface {
	messaging.Agent
	Operative(mame string) messaging.Agent
}

var (
	opsAgent *agentT
)

func init() {
	//repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
	//	return newAgent(operations.Serve)
	//})
	opsAgent = newAgent(operations.Notifier)
	exchange.Register(opsAgent)
}

type agentT struct {
	running  bool
	notifier *operations.Notification
	ex       *messaging.Exchange
}

func newAgent(notifier *operations.Notification) *agentT {
	a := new(agentT)
	a.notifier = notifier
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

func (a *agentT) registerCaseOfficer(name string) caseofficer.Agent {
	agent := caseofficer.NewAgent(fmt.Sprintf(caseOfficerNameFmt, name))
	a.ex.Register(agent)
	return agent
}

func (a *agentT) Operative(name string) messaging.Agent {
	return a.ex.Get(name)
}

func (a *agentT) trace() {
	list := a.ex.List()
	for _, v := range list {
		fmt.Printf("trace: %v -> %v\n", a.Name(), v)
	}
}
