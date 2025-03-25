package routing

import (
	"errors"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/resiliency/common"
	"github.com/behavioral-ai/resiliency/request"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/routing"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	hostName string
	timeout  time.Duration

	exchange httpx.Exchange
	handler  messaging.Agent
}

// New - create a new cache agent
func New(handler messaging.Agent) request.Agent {
	return newAgent(handler)
}

func newAgent(handler messaging.Agent) *agentT {
	a := new(agentT)

	a.exchange = httpx.Do
	a.handler = handler
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

func (a *agentT) Timeout() time.Duration { return a.timeout }
func (a *agentT) Do() httpx.Exchange     { return a.exchange }

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

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if a.hostName == "" {
			err = errors.New("host configuration is empty")
			status := messaging.NewStatusError(messaging.StatusInvalidArgument, err, a.Uri())
			a.handler.Message(eventing.NewNotifyMessage(status))
			return serverErrorResponse, err
		}
		resp, err = a.do(r, uri.BuildURL(a.hostName, r.URL.Path, r.URL.Query()))
		if next != nil && resp.StatusCode == http.StatusOK {
			resp, err = next(r)
		}
		return
	}
}
