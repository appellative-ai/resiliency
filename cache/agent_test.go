package cache

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"net/http"
	"time"
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

type requesterT struct {
	hostName string
	timeout  time.Duration
	exchange httpx.Exchange
}

func newRequesterTest() *requesterT {
	a := new(requesterT)
	a.hostName = "www.google.com"
	a.timeout = 0
	a.exchange = httpx.Do
	return a
}
func (a *requesterT) Timeout() time.Duration   { return a.timeout }
func (a *requesterT) Exchange() httpx.Exchange { return a.exchange }
func (a *requesterT) Uri() string              { return NamespaceName }

func routingExchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		//agent := eventtest.New(nil)
		//rt := newRequesterTest()
		//if r.Header.Get(iox.AcceptEncoding) == "" {
		//	r.Header.Add(iox.AcceptEncoding, iox.GzipEncoding)
		//}
		//url
		//var status *messaging.Status
		h := make(http.Header)
		h.Add(iox.AcceptEncoding, iox.GzipEncoding)
		req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
		req.Header = h
		resp, err = httpx.Do(req)
		if err != nil {
			fmt.Printf("test: httx.Do() -> [err:%v]\n", err)
		}

		//url := uri.BuildURL("www.google.com", r.URL.Path, r.URL.Query())
		//resp, status = request.Do(rt, r.Method, url, httpx.CloneHeaderWithEncoding(r), r.Body)
		//if status.Err != nil {
		//status := messaging.NewStatusError(resp.StatusCode, err, NamespaceName)
		//	agent.Message(eventing.NewNotifyMessage(status))
		//	}
		return
	}
}
