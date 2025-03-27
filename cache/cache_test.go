package cache

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"io"
	"net/http"
	"time"
)

func putCache(url string, timeout time.Duration) (*http.Response, error) {
	// create request and process exchange
	ctx, cancel := httpx.NewContext(nil, timeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header = make(http.Header)
	req.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
	resp, err1 := httpx.Do(req)
	if err1 != nil {
		return resp, err1
	}

	// read response body and send to caching exchange
	httpx.TransformBody(resp)
	cr, _ := http.NewRequest(http.MethodPut, url, resp.Body)
	cr.Header = httpx.CloneHeader(resp.Header)
	return cachingExchange(cr)
}

func ExampleCache_No_Timeout() {
	url := "https://www.google.com/search?q=golang"
	timeout := time.Millisecond * 0
	fmt.Printf("test: ExampleCache() [url:%v] [timeout:%v]\n", url, timeout)

	resp, err := putCache(url, timeout)
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
	//test: ExampleCache() [url:https://www.google.com/search?q=golang] [timeout:0]
	//test: cachingExchange.Put() [status:200] [err:<nil>]
	//test: cachingExchange.Get() [status:200] [header:true] [err:<nil>]
	//test: io.ReadAll() [err:<nil>] [buf:40845]

}

func ExampleCache_Timeout_504() {
	url := "https://www.google.com/search?q=erlang"
	timeout := time.Millisecond * 10
	fmt.Printf("test: ExampleCache() [url:%v] [timeout:%v]\n", url, timeout)

	resp, err := putCache(url, timeout)
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
	//test: ExampleCache() [url:https://www.google.com/search?q=erlang] [timeout:10ms]
	//test: cachingExchange.Put() [status:504] [err:Get "https://www.google.com/search?q=erlang": context deadline exceeded]
	//test: cachingExchange.Get() [status:404] [header:true] [err:<nil>]
	//test: io.ReadAll() [err:<nil>] [buf:0]

}

func ExampleCache_Timeout_200() {
	url := "https://www.google.com/search?q=pascal"
	timeout := time.Second * 5
	fmt.Printf("test: ExampleCache() [url:%v] [timeout:%v]\n", url, timeout)

	resp, err := putCache(url, timeout)
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
	//test: ExampleCache() [url:https://www.google.com/search?q=pascal] [timeout:5s]
	//test: cachingExchange.Put() [status:200] [err:<nil>]
	//test: cachingExchange.Get() [status:200] [header:true] [err:<nil>]
	//test: io.ReadAll() [err:<nil>] [buf:40853]

}
