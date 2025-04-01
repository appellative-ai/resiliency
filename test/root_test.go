package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleRootExchange() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/searchq=golang", nil)
	rec := httptest.NewRecorder()
	handler := NewRootEndpoint()
	handler.Exchange(rec, req)

	fmt.Printf("test: RootEndpoint() -> [status:%v]\n", rec.Result().StatusCode)

	//Output:
	//fail

}
