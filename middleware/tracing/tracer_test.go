package tracing

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/internal/testdata/binding"
	"go.opentelemetry.io/otel/trace"
)

func Test_NewTracer(t *testing.T) {
	tracer := NewTracer(trace.SpanKindClient, func(o *options) {
		o.tracerProvider = trace.NewNoopTracerProvider()
	})

	if tracer.kind != trace.SpanKindClient {
		t.Errorf("The tracer kind must be equal to trace.SpanKindClient, %v given.", tracer.kind)
	}

	defer func() {
		if recover() == nil {
			t.Error("The NewTracer with an invalid SpanKindMustCrash must panic")
		}
	}()
	_ = NewTracer(666, func(o *options) {
		o.tracerProvider = trace.NewNoopTracerProvider()
	})
}

func Test_Tracer_End(t *testing.T) {
	tracer := NewTracer(trace.SpanKindClient, func(o *options) {
		o.tracerProvider = trace.NewNoopTracerProvider()
	})
	ctx, span := trace.NewNoopTracerProvider().Tracer("noop").Start(context.Background(), "noopSpan")

	// Handle with error case
	tracer.End(ctx, span, nil, errors.New("dummy error"))

	// Handle without error case
	tracer.End(ctx, span, nil, nil)

	m := &binding.HelloRequest{}

	// Handle the trace KindServer
	tracer = NewTracer(trace.SpanKindServer, func(o *options) {
		o.tracerProvider = trace.NewNoopTracerProvider()
	})
	tracer.End(ctx, span, m, nil)
	tracer = NewTracer(trace.SpanKindClient, func(o *options) {
		o.tracerProvider = trace.NewNoopTracerProvider()
	})
	tracer.End(ctx, span, m, nil)
}
