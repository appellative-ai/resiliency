package endpoint

import (
	"net/http"
)

var ()

type health struct{}

func NewHealthEndpoint() *health {
	return new(health)
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}
