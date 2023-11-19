package instrumentation

import (
	"context"
	"os"

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

type TraceScope string

const (
	CommandScope TraceScope = "command"
	EventScope   TraceScope = "event"
)

type Client interface {
	Shutdown()
	SpanInContext(ctx context.Context, traceScope TraceScope, spanName string) (context.Context, trace.Span)
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

func (o *instrumentationClient) SpanInContext(ctx context.Context, traceScope TraceScope, spanName string) (context.Context, trace.Span) {
	tracer := o.getTracer(ctx, string(traceScope))
	ctx, span := tracer.Start(ctx, spanName)
	return ctx, span
}

func (o *instrumentationClient) getTracer(ctx context.Context, scopeName string) trace.Tracer {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span.TracerProvider().Tracer(scopeName)
	}
	return otel.GetTracerProvider().Tracer(scopeName)
}
