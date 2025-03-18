package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleExchange() {
	req, _ := http.NewRequest("", "http://localhost:8083/github/advanced-go/agency:resiliency?action=start", nil)
	rec := httptest.NewRecorder()
	Exchange(rec, req)
	rec.Flush()
	buf, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		fmt.Printf("test: io.ReadAlle() -> [err:%v]\n", err)
	} else {

		fmt.Printf("test: Exchange() -> [code:%v] [%v]\n", rec.Result().StatusCode, string(buf))
	}

	req, _ = http.NewRequest("", "http://localhost:8083/resiliency", nil)
	rec = httptest.NewRecorder()
	Exchange(rec, req)
	rec.Flush()
	buf, err = io.ReadAll(rec.Result().Body)
	if err != nil {
		fmt.Printf("test: io.ReadAlle() -> [err:%v]\n", err)
	} else {

		fmt.Printf("test: Exchange() -> [code:%v] [%v]\n", rec.Result().StatusCode, string(buf))
	}

	req, _ = http.NewRequest("", "http://localhost:8083/resiliency?dummy=1", nil)
	rec = httptest.NewRecorder()
	Exchange(rec, req)
	rec.Flush()
	buf, err = io.ReadAll(rec.Result().Body)
	if err != nil {
		fmt.Printf("test: io.ReadAlle() -> [err:%v]\n", err)
	} else {

		fmt.Printf("test: Exchange() -> [code:%v] [%v]\n", rec.Result().StatusCode, string(buf))
	}

	req, _ = http.NewRequest("", "http://localhost:8083/resiliency?event=invalid", nil)
	rec = httptest.NewRecorder()
	Exchange(rec, req)
	rec.Flush()
	buf, err = io.ReadAll(rec.Result().Body)
	if err != nil {
		fmt.Printf("test: io.ReadAlle() -> [err:%v]\n", err)
	} else {
		fmt.Printf("test: Exchange() -> [code:%v] [%v]\n", rec.Result().StatusCode, string(buf))
	}

	//Output:
	//test: Exchange() -> [code:400] [error: invalid path]
	//test: Exchange() -> [code:400] [error: no query args]
	//test: Exchange() -> [code:400] [error: event query key not found]
	//test: Exchange() -> [code:400] [AgentMessage() -> [error: invalid event] [event:invalid]]

}
