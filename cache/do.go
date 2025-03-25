package cache

import (
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"io"
	"net/http"
	"time"
)

func (a *agentT) do(url string, h http.Header, method string, r io.ReadCloser) (resp *http.Response, err error) {
	ctx, cancel := httpx.NewContext(a.timeout)
	defer cancel()
	start := time.Now().UTC()
	req, err1 := http.NewRequestWithContext(ctx, method, url, r)
	if err1 != nil {
		status := messaging.NewStatusError(messaging.StatusInvalidArgument, err1, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
		return serverErrorResponse, err
	}
	req.Header = h
	resp, err = a.exchange(req)
	// Handle error, but continue as response status code reflects error
	if err != nil {
		status := messaging.NewStatusError(http.StatusBadRequest, err, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
	} else {
		// If there is a timeout context, then transform the body
		if a.timeout > 0 {
			err = httpx.TransformBody(resp)
			if err != nil {
				status := messaging.NewStatusError(messaging.StatusIOError, err, a.Uri())
				a.handler.Message(eventing.NewNotifyMessage(status))
			}
		}
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: a.timeout})
	return
}
