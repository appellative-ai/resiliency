package http

import (
	"github.com/behavioral-ai/resiliency/operations"
	"net/http"
)

// http://localhost:8080/resiliency?event=startup

const (
	operationsResource = "operations"
	eventKey           = "event"
)

// Exchange - HTTP exchange function
func Exchange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/"+operationsResource {
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
	err := operations.Message("event:" + event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
