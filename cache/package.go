package cache

import (
	http2 "github.com/behavioral-ai/core/http"
	"net/http"
)

// Agent needs to be configured with the URL of the caching service
var (
	Agent = New()
)

func Exchange(r *http.Request) (*http.Response, error) {
	resp, err := Agent.Exchange(r)
	if err != nil {
		return resp, err
	}
	return nil, nil
}

func NewIntermediary() http2.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		resp, err = Agent.Exchange(r)
		return
	}
}
