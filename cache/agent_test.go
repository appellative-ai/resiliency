package cache

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"io"
	"net/http"
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

var cache = make(map[string]*http.Response)

func cachingExchange(r *http.Request) (*http.Response, error) {
	key := r.URL.String()
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

func routingExchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if r.Header.Get(iox.AcceptEncoding) == "" {
			r.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
		}
		resp, err = httpx.Do(r)
		return
	}
}

func ExampleCache() {
	url := "https://www.google.com/search?q=golang"
	resp, err := putCache(url, 0)
	fmt.Printf("test: cachingExchange.Put() [status:%v] [err:%v]\n", resp.StatusCode, err)

	// Get cached response
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)
	req.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
	resp, err = cachingExchange(req)
	fmt.Printf("test: cachingExchange.Get() [status:%v] [header:%v] [err:%v]\n", resp.StatusCode, resp.Header != nil, err)

	// verify that the response body can be read
	if err == nil {
		buf, err1 := io.ReadAll(resp.Body)
		fmt.Printf("test: io.ReadAll() [err:%v] [buf:%v]\n", err1, len(buf))
	}

	//Output:
	//fail

}

/*
	if err == nil {
		buf, err1 := io.ReadAll(resp.Body)
		fmt.Printf("test: io.ReadAll() [err:%v] [buf:%v]\n", err1, len(buf))
	}
	fmt.Printf("test: httpx.Do() [err:%v] [encoding:%v]\n", err, resp.Header.Get(iox.ContentEncoding))
*/

/*
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


*/
/*
	// create request and exchange
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)
	req.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
	resp, err := httpx.Do(req)

	// Read response body and send to caching exchange
	httpx.TransformBody(resp)
	cr, _ := http.NewRequest(http.MethodPut, url, resp.Body)
	cr.Header = httpx.CloneHeader(resp.Header)
	resp, err = cachingExchange(cr)
	fmt.Printf("test: cachingExchange() [err:%v]\n", err)



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

*/
