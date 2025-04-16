package authz

import (
	"github.com/behavioral-ai/core/rest"
	"net/http"
)

func Link(next rest.Exchange) rest.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		/*auth := r.Header.Get(Authorization)
		if auth == "" {
			return &http.Response{StatusCode: http.StatusUnauthorized}, nil
		}

		*/
		if next != nil {
			resp, err = next(r)
		} else {
			return &http.Response{StatusCode: http.StatusOK}, nil
		}
		return
	}
}
