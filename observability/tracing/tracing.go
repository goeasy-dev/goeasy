package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"goeasy.dev/util"
)

var tracer = otel.Tracer("goeasy.dev")

func Start(ctx context.Context, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	caller := util.GetCaller()
	return tracer.Start(ctx, caller.Name, options...)
}

func StartWithName(ctx context.Context, name string, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, options...)
}
