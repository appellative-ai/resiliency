package invoke

import (
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func Exchange(r *http.Request) (*http.Response, error) {
	status := messaging.StatusOK
	if status != nil {

	}
	return &http.Response{StatusCode: http.StatusOK}, nil
}
