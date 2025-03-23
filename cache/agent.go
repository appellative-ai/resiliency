package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"io"
	"net/http"
	"net/url"
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

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if a.hostName == "" {
			err = errors.New("cache host name is empty and not configured")
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewReader([]byte(err.Error())))}, err
		}
		url := newUrl(a.hostName, r.URL)
		h := httpx.Copy(r.Header)
		resp, err = a.get(url, h)
		if resp.StatusCode == http.StatusOK {
			resp.Header.Add(access.XCached, "true")
			return resp, nil
		}
		if next != nil {
			resp, err = next(r)
			if resp.StatusCode == http.StatusOK {
				go a.put(url, h, resp)
			}
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

func (a *agentT) get(url string, h http.Header) (*http.Response, error) {
	start := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

func (a *agentT) put(url string, h http.Header, appResp *http.Response) {
	start := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
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

func newUrl(hostName string, url *url.URL) string {
	return fmt.Sprintf("https://%v%v", hostName, url.String())
}
