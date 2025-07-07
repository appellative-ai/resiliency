package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/messaging"
)

const (
//NamespaceNamePrimary = "test:resiliency:agent/caseOfficer/service/traffic/ingress/primary"
//NetworkNamePrimary   = "test:resiliency:network/service/traffic/ingress/primary"

)

type Agent interface {
	messaging.Agent
	BuildNetwork(m map[string]map[string]string, roles []string) ([]any, []error)
	Operative(mame string) messaging.Agent
	Trace()
}

// TODO : need host name
type agentT struct {
	running bool
	name    string

	ex       *messaging.Exchange
	emissary *messaging.Channel
	service  *operations.Service
}

// NewAgent - create a new agent
func NewAgent(name string) Agent {
	return newAgent(name, operations.Serve)
}

func newAgent(name string, service *operations.Service) *agentT {
	a := new(agentT)
	a.name = name
	a.service = service

	a.ex = messaging.NewExchange()
	a.emissary = messaging.NewEmissaryChannel()
	return a
}

// Name - agent identifier
func (a *agentT) Name() string { return a.name }

func (a *agentT) Trace() {
	list := a.ex.List()
	for _, v := range list {
		fmt.Printf("trace: %v -> %v\n", a.Name(), v)
	}
}

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if !a.running {
		if m.Name == messaging.ConfigEvent {
			if m.IsRecipient(a.name) {
				messaging.UpdateAgent(a.name, func(agent messaging.Agent) {
					err := a.ex.Register(agent)
					if err != nil {
						messaging.Reply(m, messaging.NewStatus(messaging.StatusInvalidContent, err.Error()), a.Name())
					}
				}, m)
			} else {
				a.ex.Message(m)
			}
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
	if m.Name == messaging.ShutdownEvent {
		a.running = false
	}
	// System events
	switch m.Name {
	case messaging.ShutdownEvent, messaging.PauseEvent, messaging.ResumeEvent:
		a.emissary.C <- m
		a.ex.Broadcast(m)
		return
	}
	list := m.To()
	// No recipient, or only the case officer recipient
	if len(list) == 0 || len(list) == 1 && list[0] == a.name {
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
	/*
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

	*/
	// Send to appropriate agent
	a.ex.Message(m)
}

func (a *agentT) BuildNetwork(net map[string]map[string]string, roles []string) (chain []any, errs []error) {
	return buildNetwork(a, net, roles)
}

func (a *agentT) Operative(name string) messaging.Agent {
	return a.ex.Get(name)
}

// Run - run the agent
func (a *agentT) run() {
	go emissaryAttend(a)
}

func (a *agentT) shutdown() {
	a.emissary.Close()
}

/*
func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeAgent:
		agent, status := messaging.AgentContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		err := a.ex.Register(agent)
		if err != nil {
			messaging.Reply(m, messaging.NewStatus(messaging.StatusInvalidContent, err.Error()), a.Name())
		}
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}


*/
