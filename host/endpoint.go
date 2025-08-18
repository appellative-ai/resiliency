package host

import (
	"github.com/appellative-ai/core/httpx"
	"github.com/appellative-ai/core/rest"
	"net/http"
)

func ExchangeHandler(w http.ResponseWriter, req *http.Request, resp *http.Response) {
	httpx.WriteResponse(w, resp.Header, resp.StatusCode, resp.Body, req.Header)
}

func Init(r *http.Request) {
	httpx.AddRequestId(r)
}

func NewEndpoint(pattern string, operatives []any) rest.Endpoint {
	net := rest.BuildNetwork(operatives)
	return rest.NewEndpoint(pattern, ExchangeHandler, Init, net)
}
