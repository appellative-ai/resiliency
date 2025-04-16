package endpoint

import (
	"github.com/behavioral-ai/resiliency/operations"
	"net/http"
	"strings"
)

const (
	eventKey = "event"
)

type ops struct{}

func newOperationsEndpoint() *ops {
	return new(ops)
}

func (o *ops) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, OperationsPattern) {
		w.WriteHeader(http.StatusBadRequest)
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
