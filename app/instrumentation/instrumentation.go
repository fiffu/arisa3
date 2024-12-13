package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/app/log"
	honeycomb "github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	envvarServiceName = "OTEL_SERVICE_NAME"
	envvarAPIKey      = "HONEYCOMB_API_KEY"
)

type ScopedName interface {
	scope() supportedScope
	name() string
}

type Client interface {
	Shutdown()
}

type instrumentationClient struct {
	shutdownFunc func()
}

func NewInstrumentationClient(ctx context.Context) (Client, error) {
	serviceName, _ := os.LookupEnv(envvarServiceName)
	_, haveAPIKey := os.LookupEnv(envvarAPIKey)
	log.Infof(ctx, "Instrumentation client enabled: %v serviceName: %s", haveAPIKey, serviceName)

	// use honeycomb distro to setup OpenTelemetry SDK
	// enable multi-span attributes
	bsp := honeycomb.NewBaggageSpanProcessor()
	shutdownFunc, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithSpanProcessor(bsp),
	)
	if err != nil {
		log.Errorf(ctx, err, "Error setting up OTel SDK")
	}
	return &instrumentationClient{shutdownFunc}, nil
}

func (o *instrumentationClient) Shutdown() {
	o.shutdownFunc()
}

func SpanInContext(ctx context.Context, sn ScopedName) (context.Context, trace.Span) {
	tracer := fromCtx(ctx, sn.scope())
	ctx, span := tracer.Start(ctx, sn.name())
	return ctx, span
}

func EmitErrorf(ctx context.Context, msg string, args ...any) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(fmt.Errorf(msg, args...), WithStackTrace())
}

// WithStackTrace is a wrapper over `trace.EventOption`.
// This is meant to be used with span.RecordError().
func WithStackTrace() trace.EventOption {
	return trace.WithStackTrace(true)
}

func fromCtx(ctx context.Context, scope supportedScope) trace.Tracer {
	scopeName := string(scope)
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span.TracerProvider().Tracer(scopeName)
	}
	return otel.GetTracerProvider().Tracer(scopeName)
}

func NewHTTPTransport(tpt http.RoundTripper) http.RoundTripper {
	return requests.RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		spanName := httpSpanNameFormatter(req)

		ctx := req.Context()
		ctx, span := SpanInContext(ctx, ExternalHTTP(spanName))
		req = req.WithContext(ctx)
		defer span.End()

		reqSize := req.ContentLength
		span.SetAttributes(
			KV.HTTPHost(req.Host),
			KV.HTTPMethod(req.Method),
			KV.HTTPPath(req.URL.EscapedPath()),
		)

		res, err = tpt.RoundTrip(req)
		if res != nil {
			resSize := res.ContentLength

			span.SetAttributes(
				KV.HTTPTotalContentLength(reqSize+resSize),
				KV.HTTPRespStatusCode(res.StatusCode),
			)
		}
		if err != nil {
			span.SetAttributes(
				KV.Error(err),
			)
		}
		return
	})
}

func httpSpanNameFormatter(r *http.Request) string {
	path := r.URL.Host + r.URL.EscapedPath()
	if discordAPIPath := MatchDiscordAPIPath(r.Context(), path); discordAPIPath != "" {
		path = discordAPIPath
	}
	return fmt.Sprintf("%s %s", r.Method, path)
}
