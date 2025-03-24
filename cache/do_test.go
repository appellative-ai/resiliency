package cache

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"time"
)

func ExampleDo_Get() {
	a := newAgent(eventtest.New(nil))
	a.hostName = ""
	a.timeout = time.Second * 4

	h := make(http.Header)
	h.Add(iox.AcceptEncoding, "gzip")
	h.Add(httpx.XRequestId, "1234-request-id")
	resp, err := a.do("https://www.google.com/search?q=golang", h, http.MethodGet, nil)

	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	resp, err = a.do("https://www.google1234.com/search?q=golang", h, http.MethodGet, nil)
	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: do() -> [resp:200] [err:<nil>]
	//notify-> 2025-03-24T19:21:15.812Z [resiliency:agent/behavioral-ai/resiliency/cache] [core:messaging.status] [] [Bad Request] []
	//test: do() -> [resp:500] [err:<nil>]

}

func ExampleDo_Get_Timeout() {
	a := newAgent(eventtest.New(nil))
	a.timeout = time.Millisecond * 10

	h := make(http.Header)
	h.Add(iox.AcceptEncoding, "gzip")
	h.Add(httpx.XRequestId, "1234-request-id")
	resp, err := a.do("https://www.google.com/search?q=golang", h, http.MethodGet, nil)

	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: do() -> [resp:504] [err:<nil>]

}

func testExchange(r *http.Request) (*http.Response, error) {
	//fmt.Printf()
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func ExampleDo_Put() {
	a := newAgent(eventtest.New(nil))
	a.timeout = time.Second * 4
	a.exchange = testExchange
	h := make(http.Header)
	h.Add(httpx.XRequestId, "1234-request-id")
	resp, err := a.do("https://www.google.com/search?q=golang", h, http.MethodPut, nil)

	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: do() -> [resp:200] [err:<nil>]

}
