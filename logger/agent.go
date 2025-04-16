package operations

import (
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/logger"
	Route         = "host"
)

var (
	Agent = New()
)

type agentT struct {
	operators []access.Operator
}

// New - create a new operations agent
func New() messaging.Agent {
	return newAgent()
}

func newAgent() *agentT {
	a := new(agentT)
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

// Link - chainable exchange
func (a *agentT) Link(next rest.Exchange) rest.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		start := time.Now().UTC()
		resp, err = next(r)
		access.Log(a.operators, access.IngressTraffic, start, time.Since(start), Route, r, resp, newThreshold(resp))
		return
	}
}

func newThreshold(resp *http.Response) access.Threshold {
	limit := resp.Header.Get(access.XRateLimit)
	resp.Header.Del(access.XRateLimit)
	timeout := resp.Header.Get(access.XTimeout)
	resp.Header.Del(access.XTimeout)
	redirect := resp.Header.Get(access.XRedirect)
	resp.Header.Del(access.XRedirect)
	return access.Threshold{Timeout: timeout, RateLimit: limit, Redirect: redirect}
}

// TODO : need configuration for operators
func (a *agentT) configure(m *messaging.Message) {
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
