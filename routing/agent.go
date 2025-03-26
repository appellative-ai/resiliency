package routing

import (
	"errors"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
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
func New(handler messaging.Agent) messaging.Agent {
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

// Timeout - implementation for Requester interface
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

// Link - implementation for httpx.Chainable interface
func (a *agentT) Link(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if a.hostName == "" {
			status := messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("host configuration is empty"), a.Uri())
			a.handler.Message(eventing.NewNotifyMessage(status))
			return serverErrorResponse, status.Err
		}
		var status *messaging.Status

		h := httpx.CloneHeader(r.Header)
		if r.Method == http.MethodGet && h.Get(iox.AcceptEncoding) == "" {
			h.Add(iox.AcceptEncoding, iox.GzipEncoding)
		}
		resp, status = request.Do(a, r.Method, uri.BuildURL(a.hostName, r.URL.Path, r.URL.Query()), h, r.Body)
		if status.Err != nil {
			status.WithAgent(a.Uri())
			a.handler.Message(eventing.NewNotifyMessage(status))
		}
		return resp, status.Err
	}
}
