package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceNameSecondary = "test:resiliency:agent/caseOfficer/service/traffic/ingress/secondary"
)

// TODO : need host name
type secondaryAgentT struct {
	running bool

	ex       *messaging.Exchange
	emissary *messaging.Channel
	service  *operations.Service
}

// NewSecondaryAgent - create a new agent
func NewSecondaryAgent(service *operations.Service) Agent {
	return newSecondaryAgent(service)
}

func newSecondaryAgent(service *operations.Service) *secondaryAgentT {
	a := new(secondaryAgentT)
	a.service = service

	a.ex = messaging.NewExchange()
	a.emissary = messaging.NewEmissaryChannel()
	return a
}

// Name - agent identifier
func (a *secondaryAgentT) Name() string { return NamespaceNameSecondary }

func (a *secondaryAgentT) Trace() {
	for _, v := range a.ex.List() {
		fmt.Printf("trace: %v -> %v\n", NamespaceNameSecondary, v)
	}
}

// Message - message the agent
func (a *secondaryAgentT) Message(m *messaging.Message) {
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
	if len(list) == 0 || len(list) == 1 && list[0] == NamespaceNamePrimary {
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

func (a *secondaryAgentT) BuildNetwork(net map[string]map[string]string) (chain []any, errs []error) {
	return buildNetwork(a, net)
}

func (a *secondaryAgentT) Operative(name string) messaging.Agent {
	return a.ex.Get(name)
}

// Run - run the agent
func (a *secondaryAgentT) run() {
	go secondaryEmissaryAttend(a)
}

func (a *secondaryAgentT) shutdown() {
	a.emissary.Close()
}

func (a *secondaryAgentT) configure(m *messaging.Message) {
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
