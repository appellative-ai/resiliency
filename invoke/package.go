package invoke

import "net/http"

func Exchange(r*http.Request) (*http.Response,error) {
	return &http.Response{StatusCode: http.StatusOK},nil
}
