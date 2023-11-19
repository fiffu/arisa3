package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/fiffu/arisa3/app/log"
	honeycomb "github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

func fromCtx(ctx context.Context, scope supportedScope) trace.Tracer {
	scopeName := string(scope)
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span.TracerProvider().Tracer(scopeName)
	}
	return otel.GetTracerProvider().Tracer(scopeName)
}

func NewHTTPTransport(base http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(base, otelhttp.WithSpanNameFormatter(httpSpanNameFormatter))
}

func httpSpanNameFormatter(operation string, r *http.Request) string {
	return fmt.Sprintf("HTTP %s %s%s", r.Method, r.URL.Host, r.URL.EscapedPath())
}
