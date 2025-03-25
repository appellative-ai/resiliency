package cache

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"net/http"
)

func ExampleNew() {
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

func testExchange2(r *http.Request) (resp *http.Response, err error) {
	resp = httpx.NewResponse(http.StatusOK, nil, nil)

	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
	}
	return
}
