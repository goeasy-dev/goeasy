package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	defaultBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
)

type Duration struct {
	duration metric.Float64Histogram
}

func (d *Duration) Start(ctx context.Context, attributes ...attribute.KeyValue) func() {
	start := time.Now()

	return func() {
		d.duration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(attributes...))
	}
}

func NewDuration(name string) *Duration {
	duration, err := meter.Float64Histogram(name,
		metric.WithExplicitBucketBoundaries(defaultBuckets...),
	)
	if err != nil {
		panic(err)
	}

	return &Duration{
		duration: duration,
	}
}
