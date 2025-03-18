package http

import (
	"github.com/behavioral-ai/operations/module"
	"net/http"
)

// http://localhost:8080/resiliency?event=startup

const (
	resiliencyResource = "resiliency"
	eventKey           = "event"
)

// Exchange - HTTP exchange function
func Exchange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/"+resiliencyResource {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: invalid path"))
		return
	}
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
	err := module.AgentMessage("event:" + event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
