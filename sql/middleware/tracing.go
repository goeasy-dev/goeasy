package middleware

import (
	"database/sql"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"goeasy.dev/observability/tracing"
	"goeasy.dev/util"
)

func TraceMiddleware(next SQLRequestFunc) SQLRequestFunc {
	return func(request SQLRequest) error {
		ctx, span := tracing.StartWithName(
			request.context,
			"query",
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attribute.String("query", request.query)),
			trace.WithAttributes(attribute.String("query_hash", util.MD5(request.query))),
		)
		defer span.End()
		request.context = ctx

		err := next(request)
		if err != nil && err != sql.ErrNoRows {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return err
	}
}
