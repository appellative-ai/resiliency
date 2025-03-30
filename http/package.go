package http

import (
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/operations"
	"github.com/behavioral-ai/resiliency/routing"
	"github.com/behavioral-ai/traffic/analytics"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
	"net/http"
	"strings"
)

// http://localhost:8080/resiliency?event=startup

const (
	operationsResource = "operations"
	eventKey           = "event"
)

var (
	chain = httpx.BuildChain(host.AccessLogLink, host.AuthorizationLink, redirect.Agent,
		analytics.Agent, cache.Agent, limiter.Agent, routing.Agent)
)

// Exchange - HTTP exchange function
func Exchange(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/"+operationsResource) {
		opsRequest(w, r)
		return
	}
	host.Exchange(w, r, chain)
}

func opsRequest(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: no query args"))
		return
	}
	event := values.Get(eventKey)
	if event == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: event query key not found"))
		return
	}
	err := operations.Message("event:" + event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
