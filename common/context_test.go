package common

import (
	"fmt"
	"github.com/behavioral-ai/core/uri"
	"net/url"
)

func ExampleNewUrl() {
	s := "http://localhost:8080/search?q=golang"
	u, _ := url.Parse(s)
	u2 := uri.BuildURL("www.google.com", "", u.Path, u.Query())
	fmt.Printf("test: NewUrl(\"%v\") -> [%v]\n", u, u2)

	//Output:
	//fail

}
