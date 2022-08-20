package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func jsonResponse(t *testing.T, responseCode int, json string) *http.Response {
	t.Helper()
	return &http.Response{
		StatusCode: responseCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer([]byte(json))),
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
				if expectURL != requestedURL {
					msg := fmt.Sprintf("unexpected request to %s, expected: %s", requestedURL, expectURL)
					panic(msg)
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
				if expectURL != requestedURL {
					msg := fmt.Sprintf("unexpected request to %s, expected: %s", requestedURL, expectURL)
					panic(msg)
				}
				return nil, stubTransportErr
			})).
			Fetch(ctx)
	}
}
