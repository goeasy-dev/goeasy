package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Counter struct {
	counter metric.Int64Counter
}

func NewCounter(name string) *Counter {
	counter, err := meter.Int64Counter(name)
	if err != nil {
		panic(err)
	}

	return &Counter{
		counter: counter,
	}
}

func (c *Counter) Inc(ctx context.Context, attributes ...attribute.KeyValue) {
	c.counter.Add(ctx, 1, metric.WithAttributes(attributes...))
}
