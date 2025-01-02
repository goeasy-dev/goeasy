package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"goeasy.dev/util"
)

var tracer = otel.Tracer("goeasy.dev")

func Start(ctx context.Context) (context.Context, trace.Span) {
	caller := util.GetCaller()
	return tracer.Start(ctx, caller.Name)
}
