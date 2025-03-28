package cache

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"net/http/httptest"
)

func ExampleChain() {
	agent := newAgent(eventtest.New(nil))
	agent.exchange = cachingExchange
	agent.hostName = "localhost:8082"

	url := "https://localhost:8081/search/google?q=golang"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)
	httpx.AddRequestId(req)

	chain := httpx.BuildChain(agent, routingExchange)
	r := httptest.NewRecorder()
	host.Exchange(r, req, chain)
	r.Flush()
	buf, err := iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	r = httptest.NewRecorder()
	host.Exchange(r, req, chain)
	r.Flush()
	buf, err = iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

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
