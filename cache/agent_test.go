package cache

import (
	"context"
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/resiliency/common"
	"io"
	"net/http"
	"strings"
)

func _ExampleNew() {
	//url := "https://www.google.com/search"
	a := newAgent(nil)

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())
	m := make(map[string]string)
	m[common.CacheHostKey] = "google.com"
	a.Message(messaging.NewConfigMessage(m))
	fmt.Printf("test: Message() -> %v\n", a.hostName)

	//Output:
	//test: newAgent() -> resiliency:agent/behavioral-ai/resiliency/cache
	//test: Message() -> google.com

}

func pathKey(r *http.Request) string {
	newUrl := strings.Builder{}
	//newUrl.WriteString(http.r.Header.Get("method")
	path := r.URL.Path
	if len(path) > 0 {
		if path[:1] != "/" {
			path += "/"
		}
	}
	newUrl.WriteString(path)
	query := r.URL.Query()
	q := uri.BuildQuery(query)
	if q != "" {
		newUrl.WriteString("?")
		newUrl.WriteString(q)
	}
	return newUrl.String()
}

var cache = make(map[string]*http.Response)

func cachingExchange(r *http.Request) (*http.Response, error) {
	key := r.URL.String() //pathKey(r)
	switch r.Method {
	case http.MethodGet:
		if resp, ok := cache[key]; ok {
			return resp, nil
		}
		return httpx.NewResponse(http.StatusNotFound, nil, nil), nil
	case http.MethodPut:
		cache[key] = &http.Response{StatusCode: http.StatusOK, Header: httpx.CloneHeader(r.Header), Body: r.Body} //io.NopCloser(bytes.NewReader(buf))}
	}
	return httpx.NewResponse(http.StatusOK, nil, nil), nil
}

func ExampleCache() {
	url := "https://www.google.com/search?q=golang"
	ctx, cancel := httpx.NewContext(nil, 0)
	defer cancel()

	// create request and process
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header = make(http.Header)
	req.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
	resp, err := httpx.Do(req)

	/*
		if err == nil {
			buf, err1 := io.ReadAll(resp.Body)
			fmt.Printf("test: io.ReadAll() [err:%v] [buf:%v]\n", err1, len(buf))
		}
		fmt.Printf("test: httpx.Do() [err:%v] [encoding:%v]\n", err, resp.Header.Get(iox.ContentEncoding))
	*/

	httpx.TransformBody(resp)
	cr, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, url, resp.Body)
	cr.Header = httpx.CloneHeader(resp.Header)
	resp, err = cachingExchange(cr)
	fmt.Printf("test: cachingExchange() [err:%v]\n", err)

	cr, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	//cr.Header = httpx.CloneHeader(resp.Header)
	resp, err = cachingExchange(cr)
	fmt.Printf("test: cachingExchange() [err:%v] [encoding:%v]\n", err, resp.Header.Get(iox.ContentEncoding))

	if err == nil {
		buf, err1 := io.ReadAll(resp.Body)
		fmt.Printf("test: io.ReadAll() [err:%v] [buf:%v]\n", err1, len(buf))
	}
	//Output:
	//fail

}

func routingExchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if r.Header.Get(iox.AcceptEncoding) == "" {
			r.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
		}
		resp, err = httpx.Do(r)
		return
	}
}
