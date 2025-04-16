package endpoint

import (
	"fmt"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"net/http/httptest"
	"time"
)

func ExampleNewRootEndpoint() {
	h := make(http.Header)
	//h.Add(host.Authorization, "authorization")
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/google/search?q=pascal", nil)
	req.Header = h

	rec := httptest.NewRecorder()
	handler := NewRootEndpoint()
	handler.ServeHTTP(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err := iox.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [content-type:%v] [err:%v]\n", len(buf), http.DetectContentType(buf), err)

	time.Sleep(time.Second * 2)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err = iox.ReadAll(rec.Result().Body, rec.Result().Header)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [content-type:%v] [err:%v]\n", len(buf), http.DetectContentType(buf), err)

	//Output:
	//fail

}

func _ExampleSearch_Yahoo() {
	h := make(http.Header)
	//h.Add(host.Authorization, "authorization")
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/yahoo/search?q=golang", nil)
	req.Header = h

	rec := httptest.NewRecorder()
	handler := NewRootEndpoint()
	handler.ServeHTTP(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err := iox.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err = iox.ReadAll(rec.Result().Body, rec.Result().Header)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	//Output:
	//fail

}
