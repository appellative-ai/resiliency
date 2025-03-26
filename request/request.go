package request

import (
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"io"
	"net/http"
	"time"
)

var (
	//okResponse          = httpx.NewResponse(http.StatusOK, nil, nil)
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type Requester interface {
	Timeout() time.Duration
	Do() httpx.Exchange
}

func Do(agent Requester, method string, url string, h http.Header, r io.ReadCloser) (resp *http.Response, status *messaging.Status) {
	ctx, cancel := httpx.NewContext(agent.Timeout())
	defer cancel()
	start := time.Now().UTC()
	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return serverErrorResponse, messaging.NewStatusError(messaging.StatusInvalidArgument, err, "")
	}
	req.Header = h
	resp, err = agent.Do()(req)
	if err != nil {
		status = messaging.NewStatusError(http.StatusBadRequest, err, "")
		return
	}
	status = messaging.StatusOK()
	// transform the body as a cancel will close the connection and not allow reads
	err = httpx.TransformBody(resp)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		status = messaging.NewStatusError(messaging.StatusIOError, err, "")
	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, resp, access.Controller{Timeout: agent.Timeout()})
	return
}
