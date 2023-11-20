package instrumentation

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func CaptureInstrumentation(t *testing.T, callback func()) *testSpan {
	t.Helper()

	oldProvider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(oldProvider)

	customProvider := testProvider{t, newTestSpan(t)}
	otel.SetTracerProvider(customProvider)
	callback()

	return customProvider.baggage
}

type testProvider struct {
	*testing.T
	baggage *testSpan
}

func (tp testProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return testTracer{tp.T, tp.baggage} //lint:ignore S1016 - it's clearer this way
}

type testTracer struct {
	*testing.T
	baggage *testSpan
}

func (tt testTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	var span trace.Span
	if tt.baggage != nil {
		span = tt.baggage
	} else {
		span = newTestSpan(tt.T)
	}
	span.SetName(spanName)
	return ctx, span
}

func newTestSpan(t *testing.T) *testSpan {
	return &testSpan{
		T:          t,
		Events:     make([]string, 0),
		Errors:     make([]error, 0),
		Attributes: make(map[string]attribute.KeyValue),
	}
}

type testSpan struct {
	*testing.T
	Name        string
	Ended       bool
	Code        codes.Code
	Description string
	Events      []string
	Errors      []error
	Attributes  AttrDict
}

type AttrDict map[string]attribute.KeyValue

func (ad AttrDict) GetAsString(attrKey string) string {
	kv, ok := ad[attrKey]
	if !ok {
		return ""
	}
	return kv.Value.Emit()
}

// End completes the Span. The Span is considered complete and ready to be
// delivered through the rest of the telemetry pipeline after this method
// is called. Therefore, updates to the Span are not allowed after this
// method has been called.
func (ts *testSpan) End(options ...trace.SpanEndOption) {
	ts.Ended = true
}

// AddEvent adds an event with the provided name and options.
func (ts *testSpan) AddEvent(name string, options ...trace.EventOption) {
	if !ts.IsRecording() {
		ts.T.Fatalf("AddEvent() not allowed after Span.End() has been called")
	}
	ts.Events = append(ts.Events, name)
}

// IsRecording returns the recording state of the Span. It will return
// true if the Span is active and events can be recorded.
func (ts *testSpan) IsRecording() bool { return !ts.Ended }

// RecordError will record err as an exception span event for this span. An
// additional call to SetStatus is required if the Status of the Span should
// be set to Error, as this method does not change the Span status. If this
// span is not being recorded or err is nil then this method does nothing.
func (ts *testSpan) RecordError(err error, options ...trace.EventOption) {
	if !ts.IsRecording() {
		ts.T.Fatalf("RecordError() not allowed after Span.End has been called")
	}
}

// SpanContext returns the SpanContext of the Span. The returned SpanContext
// is usable even after the End method has been called for the Span.
func (ts *testSpan) SpanContext() trace.SpanContext { return trace.SpanContext{} } // wtf is this

// SetStatus sets the status of the Span in the form of a code and a
// description, provided the status hasn't already been set to a higher
// value before (OK > Error > Unset). The description is only included in a
// status when the code is for an error.
func (ts *testSpan) SetStatus(code codes.Code, description string) {
	ts.Code = code
	ts.Description = description
}

// SetName sets the Span name.
func (ts *testSpan) SetName(name string) { ts.Name = name }

// SetAttributes sets kv as attributes of the Span. If a key from kv
// already exists for an attribute of the Span it will be overwritten with
// the value contained in kv.
func (ts *testSpan) SetAttributes(kv ...attribute.KeyValue) {
	for _, attr := range kv {
		ts.T.Logf("SetAttributes: %s => %s", attr.Key, attr.Value.Emit())
		ts.Attributes[string(attr.Key)] = attr
	}
}

// TracerProvider returns a TracerProvider that can be used to generate
// additional Spans on the same telemetry pipeline as the current Span.
func (ts *testSpan) TracerProvider() trace.TracerProvider {
	return testProvider{ts.T, ts}
}
