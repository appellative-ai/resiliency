package common

import (
	"context"
	"net/http"
	"time"
)

const (
	localhost   = "localhost"
	httpScheme  = "http"
	httpsScheme = "https"
)

var (
	OkResponse = &http.Response{StatusCode: http.StatusOK}
	cancel     = func() {}
)

func NewContext(timeout time.Duration) (context.Context, func()) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.Background(), cancel
}
