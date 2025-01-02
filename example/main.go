package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"goeasy.dev"
	"goeasy.dev/application/runners"
	"goeasy.dev/bootstrap"
	"goeasy.dev/cache"
	"goeasy.dev/errors"
	"goeasy.dev/observability/metrics"
	"goeasy.dev/status"
	"goeasy.dev/status/statustype"

	"go.opentelemetry.io/otel/attribute"
)

//go:generate go run ./../cmd/goeasy gen
func main() {
	ctx, done := bootstrap.Bootstrap()
	defer done()

	application := goeasy.NewApplication()
	statusCheck := status.SimpleCheck(statustype.Startup | statustype.Readiness)

	toggleRunner := runners.Timed(runners.TimedOptions{
		HandleFunc: func(ctx context.Context) error {
			log.Println("toggling status check")
			*statusCheck = !*statusCheck
			return nil
		},
		Interval: time.Second * 5,
	})

	duration := metrics.NewDuration("http_request_duration")
	httpRunner := runners.Http(runners.HttpOptions{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			record := duration.Start(r.Context(), attribute.String("path", r.URL.Path))
			defer record()

			method := r.Method
			if method == http.MethodGet {
				var dest string
				err := cache.Get(r.Context(), r.URL.Path, &dest)
				if errors.Is(err, errors.ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					return
				} else if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(dest))
			} else if method == http.MethodPost {
				defer r.Body.Close()
				key := r.URL.Path
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}

				err = cache.Put(r.Context(), key, string(body))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}

				w.WriteHeader(http.StatusCreated)
			}
		}),
	})

	err := application.Start(ctx, toggleRunner, httpRunner)
	if err != nil {
		log.Fatal(err)
	}
}
