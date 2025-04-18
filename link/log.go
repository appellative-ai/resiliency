package operations

import (
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/rest"
	"net/http"
	"time"
)

const (
	Route = "host"
)

// Logger - chainable exchange
func Logger(next rest.Exchange) rest.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		start := time.Now().UTC()
		resp, err = next(r)
		access.Log(access.IngressTraffic, start, time.Since(start), Route, r, resp, newThreshold(resp))
		return
	}
}

func newThreshold(resp *http.Response) access.Threshold {
	limit := resp.Header.Get(access.XRateLimit)
	resp.Header.Del(access.XRateLimit)
	timeout := resp.Header.Get(access.XTimeout)
	resp.Header.Del(access.XTimeout)
	redirect := resp.Header.Get(access.XRedirect)
	resp.Header.Del(access.XRedirect)
	return access.Threshold{Timeout: timeout, RateLimit: limit, Redirect: redirect}
}
