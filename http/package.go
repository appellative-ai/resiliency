package http

import (
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/resiliency/cache"
	"github.com/behavioral-ai/resiliency/limit"
	"github.com/behavioral-ai/resiliency/operations"
	"github.com/behavioral-ai/resiliency/redirect"
	"github.com/behavioral-ai/resiliency/routing"
	"net/http"
	"strings"
)

// http://localhost:8080/resiliency?event=startup

const (
	operationsResource = "operations"
	eventKey           = "event"
)

var (
	pipeline = httpx.NewExchangePipeline(redirect.Agent.Exchange, cache.Agent.Exchange, limit.Agent.Exchange, routing.Agent.Exchange)
)

// Exchange - HTTP exchange function
func Exchange(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/"+operationsResource) {
		opsRequest(w, r)
		return
	}
	resp, err := pipeline.Run(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		httpx.WriteResponse(w, resp.Header, resp.StatusCode, resp.Body, r.Header)
	}
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
