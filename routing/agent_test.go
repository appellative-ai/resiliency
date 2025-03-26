package routing

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"net/http"
	"time"
)

func ExampleNew() {
	a := newAgent(nil)

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())

	m := make(map[string]string)
	m[common.AppHostKey] = "google.com"
	a.Message(messaging.NewConfigMessage(m))
	time.Sleep(time.Second * 2)
	fmt.Printf("test: Message() -> %v\n", a.hostName)

	//Output:
	//test: newAgent() -> resiliency:agent/behavioral-ai/resiliency/routing
	//test: Message() -> google.com

}

func ExampleExchange() {
	url := "http://localhost:8080/search?q=golang"
	a := newAgent(eventtest.New(nil))
	ex := a.Link(nil)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err := ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	a.hostName = "www.google.com"
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err = ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//notify-> 2025-03-25T14:44:49.521Z [resiliency:agent/behavioral-ai/resiliency/routing] [core:messaging.status] [] [Invalid Argument] [host configuration is empty]
	//test: Exchange() -> [resp:500] [err:host configuration is empty]
	//test: Exchange() -> [resp:200] [err:<nil>]

}
