package routing

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"time"
)

func ExampleDo_Get() {
	url := "https://www.google.com/search?q=golang"
	a := newAgent(eventtest.New(nil))
	a.timeout = time.Second * 4

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	h := make(http.Header)
	h.Add(httpx.XRequestId, "1234-request-id")

	resp, err := a.do(req, url)
	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	if resp.StatusCode == http.StatusOK {
		buf, err1 := iox.ReadAll(resp.Body, resp.Header)
		fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err1)
	}

	//Output:
	//test: do() -> [resp:200] [err:<nil>]
	//test: iox.ReadAll() -> [buf:82676] [err:<nil>]

}

func ExampleDo_Get_Timeout() {
	url := "https://www.google.com/search?q=golang"
	a := newAgent(eventtest.New(nil))
	a.hostName = "localhost:8080"
	a.timeout = time.Millisecond * 10

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	h := make(http.Header)
	h.Add(iox.AcceptEncoding, "gzip")
	h.Add(httpx.XRequestId, "1234-request-id")

	resp, err := a.do(req, url)
	fmt.Printf("test: do() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: do() -> [resp:504] [err:<nil>]

}
