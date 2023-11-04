package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
)

func panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	panic(msg)
}

func jsonResponse(t *testing.T, responseCode int, json string) *http.Response {
	t.Helper()
	return &http.Response{
		StatusCode: responseCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer([]byte(json))),
	}
}

func htmlResponse(t *testing.T, responseCode int, html string) *http.Response {
	t.Helper()
	return &http.Response{
		StatusCode: responseCode,
		Header:     http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
		Body:       io.NopCloser(strings.NewReader(html)),
	}
}

func cmpURL(t *testing.T, url1, url2 string) bool {
	var u1, u2 *url.URL
	var err error
	if u1, err = url.Parse(url1); err != nil {
		panicf("invalid url '%s', err: %v", url1, err)
	}
	if u2, err = url.Parse(url2); err != nil {
		t.Fatalf("invalid url '%s', err: %v", url2, err)
	}
	return u1.EscapedPath() == u2.EscapedPath()
}

func StubHTMLFetcher(
	t *testing.T,
	expectURL string,
	stubStatusCode int,
	stubHTML string,
) func(ctx context.Context, builder *requests.Builder) error {
	t.Helper()
	resp := htmlResponse(t, stubStatusCode, stubHTML)

	return func(ctx context.Context, builder *requests.Builder) error {
		return builder.
			Transport(requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				requestedURL := req.URL.String()
				if !cmpURL(t, expectURL, requestedURL) {
					panicf("unexpected request to %s, expected: %s", requestedURL, expectURL)
				}
				return resp, nil
			})).
			Fetch(ctx)
	}
}

func StubJSONFetcher(
	t *testing.T,
	expectURL string,
	stubStatusCode int,
	stubJSON string,
) func(ctx context.Context, builder *requests.Builder) error {

	t.Helper()
	resp := jsonResponse(t, stubStatusCode, stubJSON)

	return func(ctx context.Context, builder *requests.Builder) error {
		return builder.
			Transport(requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				requestedURL := req.URL.String()
				if !cmpURL(t, expectURL, requestedURL) {
					panicf("unexpected request to %s, expected: %s", requestedURL, expectURL)
				}
				return resp, nil
			})).
			Fetch(ctx)
	}
}

func StubTransportError(
	t *testing.T,
	expectURL string,
	stubTransportErr error,
) func(ctx context.Context, builder *requests.Builder) error {
	return func(ctx context.Context, builder *requests.Builder) error {

		return builder.
			Transport(requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				requestedURL := req.URL.String()
				if !cmpURL(t, expectURL, requestedURL) {
					panicf("unexpected request to %s, expected: %s", requestedURL, expectURL)
				}
				return nil, stubTransportErr
			})).
			Fetch(ctx)
	}
}
