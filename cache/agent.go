package cache

import (
	"bytes"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
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

//err = errors.New("cache host name is empty and not configured")
//return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewReader([]byte(err.Error())))}, err

func (a *agentT) cacheEnabled() bool {
	return a.hostName != ""
}

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		uri := common.NewUrl(a.hostName, r.URL)
		h := httpx.Copy(r.Header)

		if a.cacheEnabled() {
			resp, err = a.do(uri, h, http.MethodGet, nil)
			if resp.StatusCode == http.StatusOK {
				resp.Header.Add(access.XCached, "true")
				return resp, nil
			}
		}
		if next != nil {
			resp, err = next(r)
			if resp.StatusCode == http.StatusOK && a.cacheEnabled() {
				go a.do(uri, h, http.MethodPut, resp.Body)
			}
		}
		return
	}
}

func (a *agentT) do(url string, h http.Header, method string, body io.ReadCloser) (*http.Response, error) {
	ctx, cancel := common.NewContext(a.timeout)
	if cancel != nil {
		defer cancel()
	}
	start := time.Now().UTC()
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewReader([]byte(err.Error())))}, err
	}
	req.Header = h
	resp, _ := httpx.Do(req)
	reasonCode := ""
	if resp.StatusCode == http.StatusGatewayTimeout {
		reasonCode = access.ControllerTimeout
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: a.timeout, Code: reasonCode})
	return resp, nil
}

/*
func (a *agentT) put(url string, h http.Header, appResp *http.Response) {
	ctx,cancel := common.NewContext(a.timeout)
	if cancel != nil {
		defer cancel()
	}
	start := time.Now().UTC()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, appResp.Body)
	if err != nil {
		return
	}
	req.Header = h
	resp, _ := httpx.Do(req)
	reasonCode := ""
	if resp.StatusCode == http.StatusGatewayTimeout {
		reasonCode = access.ControllerTimeout
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: a.timeout, Code: reasonCode})
	return
}


*/
