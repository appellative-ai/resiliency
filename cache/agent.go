package cache

import (
	"bufio"
	"bytes"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
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
	exchange httpx.Exchange
	handler  messaging.Agent
}

// New - create a new cache agent
func New(handler messaging.Agent) httpx.Agent {
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
			url string
		)

		if a.enabled(r) {
			url = uri.BuildURL(a.hostName, "", r.URL.Path, r.URL.Query())
			h := httpx.CloneHeader(r.Header)
			h.Add(iox.AcceptEncoding, "gzip")
			resp, err = a.do(url, h, http.MethodGet, nil)
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
				go func() {
					h := httpx.CloneHeader(r.Header)
					go a.do(url, h, http.MethodPut, io.NopCloser(bufio.NewReader(buf)))
				}()
			}
		} else {
			resp = common.OkResponse
			err = nil
		}
		return
	}
}
