module goeasy.dev

go 1.23.4

replace goeasy.dev/cmd/goeasy => ./cmd/goeasy

require (
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/otel v1.33.0
	go.opentelemetry.io/otel/metric v1.33.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.9.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.33.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.33.0
	go.opentelemetry.io/otel/log v0.9.0
	go.opentelemetry.io/otel/sdk v1.33.0
	go.opentelemetry.io/otel/sdk/log v0.9.0
	go.opentelemetry.io/otel/sdk/metric v1.33.0
	go.opentelemetry.io/otel/trace v1.33.0
	golang.org/x/sys v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
