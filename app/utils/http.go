package utils

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib/functional"
)

func ReadAndReplaceBody(r *http.Response) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func NewInstrumentedTransport() http.RoundTripper {
	return requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		reqID := newRequestID()
		startTime := time.Now()

		ctx := req.Context()
		ctx = log.Put(ctx, log.TraceSubID, reqID)
		defer log.Pop(ctx, log.TraceSubID)

		logRequest(req, startTime)
		res, err := instrumentation.NewHTTPTransport(http.DefaultTransport).RoundTrip(req)
		logResponse(req, res, startTime, err)

		return res, err
	})
}

func newRequestID() string {
	seed := time.Now().UnixMicro()
	s := strconv.FormatInt(seed, 36)
	return strings.ToUpper(s)
}

func logRequest(req *http.Request, startTime time.Time) {
	ctx := req.Context()
	log.Infof(ctx, "%s %s", req.Method, req.URL.String())
}

func logResponse(req *http.Request, res *http.Response, startTime time.Time, reqErr error) {
	ctx := req.Context()
	elapsed := time.Since(startTime)

	body, readErr := ReadAndReplaceBody(res)

	if reqErr != nil {
		log.Errorf(ctx, reqErr, "%s %s in %dms - request error: %s", req.Method, res.Status, elapsed.Milliseconds(), reqErr)
	}
	if readErr != nil {
		log.Errorf(ctx, readErr, "%s %s in %dms - io error: %s", req.Method, res.Status, elapsed.Milliseconds(), readErr)
	}

	// Only peek the repsonse if no errors
	if reqErr == nil && readErr == nil {
		peek100Bytes := functional.SliceOf(body).Take(100) // First 100 bytes
		log.Infof(ctx, "%s %s in %dms - body: %s ...", req.Method, res.Status, elapsed.Milliseconds(), peek100Bytes)
	}
}
