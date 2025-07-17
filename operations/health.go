package operations

import (
	"github.com/appellative-ai/core/rest"
	"net/http"
)

type health struct {
	pattern string
}

func newHealthEndpoint(pattern string) rest.Endpoint {
	h := new(health)
	h.pattern = pattern
	return h
}

func (h *health) Pattern() string {
	return h.pattern
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}
