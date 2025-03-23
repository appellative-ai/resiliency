package routing

import (
	"errors"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/routing"
)

type agentT struct {
	hostName string
	timeout  time.Duration

	handler messaging.Agent
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

func (a *agentT) configure(m *messaging.Message) {
	var ok bool

	if a.hostName, ok = common.AppHostName(a, m); !ok {
		return
	}
	if a.timeout, ok = common.Timeout(a, m); !ok {
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
