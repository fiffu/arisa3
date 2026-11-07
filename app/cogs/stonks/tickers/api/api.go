// Package api has the HTTP clients for the ticker data sources.
// Each client takes the http.RoundTripper that sends its requests, so callers can add
// hooks on the request lifecycle, such as instrumentation.
package api

import (
	"io"
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
)

const (
	userAgent = "arisa3 (+https://github.com/fiffu/arisa3)"
	timeout   = 30 * time.Second
	maxBody   = 4 << 20
)

func readBody(dst *[]byte) requests.ResponseHandler {
	return func(res *http.Response) error {
		b, err := io.ReadAll(io.LimitReader(res.Body, maxBody))
		*dst = b
		return err
	}
}
