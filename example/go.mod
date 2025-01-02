module goeasy.dev/example

replace goeasy.dev => ./..

replace goeasy.dev/cmd/goeasy => ./../cmd/goeasy

go 1.23.4

require (
	go.opentelemetry.io/otel v1.33.0
	goeasy.dev v0.0.0-00010101000000-000000000000
)

require (
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofiber/fiber/v2 v2.52.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	goeasy.dev/cmd/goeasy v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
