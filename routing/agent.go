package routing

import (
	"errors"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/routing"
)

// TODO : need host name
type agentT struct {
	running  bool
	hostName string
	timeout  time.Duration

	handler  messaging.Agent
	emissary *messaging.Channel
}

// New - create a new cache agent
func New() httpx.Agent {
	return newAgent(nil, "", 0)
}

func newAgent(handler messaging.Agent, hostName string, timeout time.Duration) *agentT {
	a := new(agentT)
	a.hostName = hostName
	a.timeout = timeout
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
		a.configure(m)
		return
	}
	if m.Event() == messaging.StartupEvent {
		a.run()
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
	//go emissaryAttend(a)
	a.running = true
}

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		if a.hostName == "" {
			return &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("configuration is empty and not configured")
		}
		// TODO: create request and send to the application

		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		return
	}
}

func (a *agentT) dispatch(channel any, event1 string) {
	a.handler.Message(event.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) finalize() {
	a.emissary.Close()
}

func (a *agentT) configure(m *messaging.Message) {
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Uri())
		return
	}
	a.hostName = cfg[AppHostKey]
	if a.hostName == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, AppHostKey), a.Uri())
		return
	}
	s := cfg[TimeoutKey]
	if s != "" {
		var err error
		a.timeout, err = fmtx.ParseDuration(s)
		if err != nil {
			messaging.Reply(m, messaging.ConfigContentStatusError(a, TimeoutKey), a.Uri())
			return
		}
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
