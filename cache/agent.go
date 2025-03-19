package cache

import (
	"github.com/behavioral-ai/collective/event"
	http2 "github.com/behavioral-ai/core/http"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/cache"
)

// TODO : need host name
type agentT struct {
	running bool
	url     string

	handler  messaging.Agent
	emissary *messaging.Channel
}

// New - create a new cache agent
func New() http2.Agent {
	return newAgent(nil)
}

func newAgent(handler messaging.Agent) *agentT {
	a := new(agentT)
	if handler == nil {
		a.handler = event.Agent
	} else {
		a.handler = handler
	}
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
	if m.Event() == messaging.ConfigEvent {
		if url, ok := m.Body.(string); ok {
			a.url = url
		}
	}
	if !a.running {
		return
	}
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) Run() {
	if a.running {
		return
	}
	//go emissaryAttend(a)
	a.running = true
}

func (a *agentT) Exchange(r *http.Request) (*http.Response, error) {
	return nil, nil
}

func (a *agentT) dispatch(channel any, event1 string) {
	a.handler.Message(event.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) finalize() {
	a.emissary.Close()
}
