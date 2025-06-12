package endpoint

import (
	"net/http"
)

type health struct {
	Pattern string
}

func newHealthEndpoint(pattern string) *health {
	h := new(health)
	h.Pattern = pattern
	return h
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}
