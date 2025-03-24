package cache

import (
	"bufio"
	"bytes"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"io"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/cache"
)

type agentT struct {
	hostName string
	timeout  time.Duration

	handler messaging.Agent
}

// New - create a new cache agent
func New(handler messaging.Agent) httpx.Agent {
	return newAgent(handler, "", 0)
}

func newAgent(handler messaging.Agent, hostName string, timeout time.Duration) *agentT {
	a := new(agentT)
	a.hostName = hostName
	a.timeout = timeout

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

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		var (
			uri string
		)

		if a.enabled(r) {
			uri = common.NewUrl(a.hostName, r.URL)
			h := httpx.CloneHeader(r.Header)
			h.Add(iox.AcceptEncoding, "gzip")
			resp, err = a.do(uri, h, http.MethodGet, nil)
			if resp.StatusCode == http.StatusOK {
				resp.Header.Add(access.XCached, "true")
				return resp, nil
			}
		}
		if next != nil {
			resp, err = next(r)
			if a.enabled(r) && resp.StatusCode == http.StatusOK {
				buf := &bytes.Buffer{}
				reader := io.TeeReader(resp.Body, buf)
				resp.Body = io.NopCloser(reader)
				h := httpx.CloneHeader(r.Header)
				go a.do(uri, h, http.MethodPut, io.NopCloser(bufio.NewReader(buf)))
			}
		}
		return
	}
}
