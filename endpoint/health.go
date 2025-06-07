package endpoint

import (
	"net/http"
)

type health struct{}

func newHealthEndpoint() *health {
	return new(health)
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}
