package routing

import (
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
	"time"
)

func (a *agentT) do(r *http.Request, url string) (*http.Response, error) {
	ctx, cancel := httpx.NewContext(a.timeout)
	defer cancel()
	start := time.Now().UTC()
	req, err1 := http.NewRequestWithContext(ctx, r.Method, url, r.Body)
	if err1 != nil {
		status := messaging.NewStatusError(messaging.StatusInvalidArgument, err1, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
		return serverErrorResponse, err1
	}
	req.Header = httpx.CloneHeader(r.Header)
	if req.Method == http.MethodGet && req.Header.Get(iox.AcceptEncoding) == "" {
		req.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
	}
	resp, err := a.exchange(req)
	// Handle error, but continue as response status code reflects error
	if err != nil {
		status := messaging.NewStatusError(http.StatusBadRequest, err, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
	} else {
		//var cnt int
		// Read the body into a new reader, as the connection is still active and can timeout later
		resp.ContentLength, err = httpx.TransformBody(resp)
		//resp.ContentLength = int64(cnt)
		if err != nil {
			status := messaging.NewStatusError(messaging.StatusIOError, err, a.Uri())
			a.handler.Message(eventing.NewNotifyMessage(status))
			resp = serverErrorResponse
		}
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: a.timeout})
	return resp, err
}
