package cache

import (
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"io"
	"net/http"
	"time"
)

func (a *agentT) do(url string, h http.Header, method string, r io.ReadCloser) (*http.Response, error) {
	ctx, cancel := common.NewContext(a.timeout)
	defer cancel()
	start := time.Now().UTC()
	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		status := messaging.NewStatusError(messaging.StatusInvalidArgument, err, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
		return &http.Response{StatusCode: http.StatusInternalServerError}, err
	}
	req.Header = h
	resp, err1 := httpx.Do(req)
	// Handle error, but continue as response status code reflects error
	if err1 != nil {
		status := messaging.NewStatusError(http.StatusBadRequest, err, a.Uri())
		a.handler.Message(eventing.NewNotifyMessage(status))
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: a.timeout})
	return resp, nil
}
