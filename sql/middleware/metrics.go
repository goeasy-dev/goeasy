package middleware

import (
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"goeasy.dev/observability/metrics"
)

var queryErrorMetric = metrics.NewCounter("sql_query_errors")
var queryResponseMetric = metrics.NewDuration("sql_query_time_seconds")

func MetricsMiddleware(next SQLRequestFunc) SQLRequestFunc {
	return func(request SQLRequest) error {
		attributes := []attribute.KeyValue{
			attribute.String("query", request.query),
		}

		durration := queryResponseMetric.Start(request.context, attributes...)
		defer durration()

		err := next(request)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			queryErrorMetric.Inc(request.context, attributes...)
		}

		return err
	}
}
