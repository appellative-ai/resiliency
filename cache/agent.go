package cache

import (
	"bytes"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/resiliency/common"
	"github.com/behavioral-ai/resiliency/request"
	"io"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/cache"
)

var (
	okResponse          = httpx.NewResponse(http.StatusOK, nil, nil)
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
	maxDuration         = time.Second * 15
)

type agentT struct {
	running  bool
	hostName string
	timeout  time.Duration

	exchange httpx.Exchange
	ticker   *messaging.Ticker
	emissary *messaging.Channel
	handler  messaging.Agent
}

// New - create a new cache agent
func New(handler messaging.Agent) messaging.Agent {
	return newAgent(handler)
}

func newAgent(handler messaging.Agent) *agentT {
	a := new(agentT)

	a.exchange = httpx.Do
	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
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
	go emissaryAttend(a, content.Resolver, nil)
	a.running = true
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool                { return true }
func (a *agentT) Timeout() time.Duration   { return a.timeout }
func (a *agentT) Exchange() httpx.Exchange { return a.exchange }

// Link - chainable exchange
func (a *agentT) Link(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		var (
			url    string
			status *messaging.Status
		)
		if a.enabled(r) {
			url = uri.BuildURL(a.hostName, r.URL.Path, r.URL.Query())
			h := make(http.Header)
			h.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
			resp, status = request.Do(a, http.MethodGet, url, h, nil)
			if resp.StatusCode == http.StatusOK {
				return resp, nil
			}
			if status.Err != nil {
				a.handler.Message(eventing.NewNotifyMessage(status.WithAgent(a.Uri())))
			}
		}
		if next == nil {
			return httpx.NewResponse(http.StatusNotFound, nil, nil), nil
		}
		resp, err = next(r)
		if a.enabled(r) && resp.StatusCode == http.StatusOK {
			var buf []byte
			buf, err = io.ReadAll(resp.Body)
			if err != nil {
				status = messaging.NewStatusError(messaging.StatusIOError, err, a.Uri())
				a.handler.Message(eventing.NewNotifyMessage(status))
				return serverErrorResponse, err
			}
			resp.ContentLength = int64(len(buf))
			resp.Body = io.NopCloser(bytes.NewReader(buf))
			go func() {
				h := httpx.CloneHeader(resp.Header)
				h.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
				_, status = request.Do(a, http.MethodPut, url, h, io.NopCloser(bytes.NewReader(buf)))
				if status.Err != nil {
					a.handler.Message(eventing.NewNotifyMessage(status.WithAgent(a.Uri())))
				}
			}()
		}
		return
	}
}

func (a *agentT) configure(m *messaging.Message) {
	var ok bool

	if a.hostName, ok = common.CacheHostName(a, m); !ok {
		return
	}
	if a.timeout, ok = common.Timeout(a, m); !ok {
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}

func (a *agentT) enabled(r *http.Request) bool {
	return a.hostName != "" && r.Method == http.MethodGet
}

func (a *agentT) dispatch(channel any, event1 string) {
	a.handler.Message(eventing.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) emissaryShutdown() {
	a.emissary.Close()
	a.ticker.Stop()
}
