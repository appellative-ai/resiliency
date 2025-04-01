package test

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"strings"
)

const (
	googlePath = "/google/search"
	yahooPath  = "/yahoo/search"
)

func RoutingLink(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		uri := ""
		if strings.HasPrefix(r.URL.Path, googlePath) {
			values := r.URL.Query()
			s := values.Encode()
			uri = "https://www.google.com/search?" + s
		} else {
			if strings.HasPrefix(r.URL.Path, yahooPath) {
				values := r.URL.Query()
				s := values.Encode()
				uri = "https://search.yahoo.com/search?" + s
			} else {
				return httpx.NewResponse(http.StatusBadRequest, nil, nil), err
			}
		}
		h := make(http.Header)
		h.Add(iox.AcceptEncoding, iox.GzipEncoding)
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		req.Header = h
		resp, err = httpx.Do(req)
		if err != nil {
			fmt.Printf("test: httx.Do() -> [err:%v]\n", err)
		}
		return
	}
}
