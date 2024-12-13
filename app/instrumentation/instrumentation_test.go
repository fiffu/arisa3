package instrumentation

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/app/log"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func Test_NewInstrumentationClient(t *testing.T) {

	var c Client
	var err error
	logs := log.CaptureLogging(t, func() {
		c, err = NewInstrumentationClient(context.Background())
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Contains(t, logs, "Instrumentation client enabled: false serviceName:")
}

func Test_httpSpanNameFormatter(t *testing.T) {
	testCases := []struct {
		path   string
		expect string
	}{
		{
			path:   "http://user:password@example.com/test?query=ignored",
			expect: "POST example.com/test",
		},
		{
			path:   "http://discord.com/api/v9/applications/964085462748774401/commands",
			expect: "POST discord.com/api/.+/applications/.+/commands",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			actual := httpSpanNameFormatter(req)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func Test_NewHTTPTransport(t *testing.T) {
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Write(body)
	}))

	req, err := http.NewRequest(http.MethodPatch, serv.URL, bytes.NewBufferString("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	span := CaptureInstrumentation(t, func() {
		res, err := NewHTTPTransport(http.DefaultTransport).RoundTrip(req)
		assert.NoError(t, err)

		reply, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello world"), reply)
	})
	expects := map[string]attribute.KeyValue{
		"http_host":                 KV.HTTPHost(strings.TrimPrefix(serv.URL, "http://")),
		"http_method":               KV.HTTPMethod("PATCH"),
		"http_resp_status":          KV.HTTPRespStatusCode(200),
		"http_total_content_length": KV.HTTPTotalContentLength(int64(len("hello world") * 2)),
	}
	for k, x := range expects {
		assert.Equal(t, x.Value.Emit(), span.Attributes.GetAsString(k))
	}
	assert.True(t, span.Ended)
}

func Test_NewHTTPTransport_withError(t *testing.T) {
	someError := errors.New("http transport error")
	erroringRoundTripper := requests.RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		return nil, someError
	})

	req, err := http.NewRequest(http.MethodPatch, "example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	span := CaptureInstrumentation(t, func() {
		res, err := NewHTTPTransport(erroringRoundTripper).RoundTrip(req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	assert.Equal(t, someError.Error(), span.Attributes.GetAsString("error"))
	assert.True(t, span.Ended)
}
