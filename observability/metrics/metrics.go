package metrics

import (
	"go.opentelemetry.io/otel"
)

var meter = otel.Meter("goeasy.dev")
