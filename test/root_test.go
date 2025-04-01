package test

import (
	"fmt"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"net/http/httptest"
)

func ExampleSearch_Google() {
	h := make(http.Header)
	h.Add(host.Authorization, "authorization")
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/google/search?q=golang", nil)
	req.Header = h

	rec := httptest.NewRecorder()
	handler := NewRootEndpoint()
	handler.Exchange(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err := iox.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	rec = httptest.NewRecorder()
	handler.Exchange(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err = iox.ReadAll(rec.Result().Body, rec.Result().Header)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	//Output:
	//fail

}

func _ExampleSearch_Yahoo() {
	h := make(http.Header)
	h.Add(host.Authorization, "authorization")
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/yahoo/search?q=golang", nil)
	req.Header = h

	rec := httptest.NewRecorder()
	handler := NewRootEndpoint()
	handler.Exchange(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err := iox.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	rec = httptest.NewRecorder()
	handler.Exchange(rec, req)
	fmt.Printf("test: RootEndpoint() -> [status:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header.Get(iox.ContentEncoding))

	buf, err = iox.ReadAll(rec.Result().Body, rec.Result().Header)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err)

	//Output:
	//fail

}
