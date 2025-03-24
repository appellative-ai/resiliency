package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func NewUrl(hostName string, url *url.URL) string {
	s := httpsScheme
	if strings.HasPrefix(hostName, localhost) {
		s = httpsScheme
	}
	return fmt.Sprintf("%v://%v%v", s, hostName, url.String())
}
