package runners

import (
	"context"
	"log"
	"net/http"

	"goeasy.dev"
	"goeasy.dev/errors"
	"goeasy.dev/status"
	"goeasy.dev/status/statustype"
)

// HttpOptions contains the information to handle HTTP connections
type HttpOptions struct {
	Address string
	Handler http.Handler
}

// Http begins listening for HTTP connections
func Http(config HttpOptions) goeasy.RunnerFunc {
	isRunning := status.SimpleCheck(statustype.Liveness | statustype.Readiness)

	if config.Address == "" {
		config.Address = ":8080"
	}

	return func(ctx context.Context) (goeasy.StopFunc, error) {
		srv := &http.Server{Addr: config.Address, Handler: config.Handler}
		go func() {
			*isRunning = true
			defer func() { *isRunning = false }()

			log.Printf("http server listening at %s\n", config.Address)
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(errors.Wrap(err, "error running http server"))
			}
		}()

		return func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				return errors.Wrap(err, "error shutting down HTTP server")
			}

			return nil
		}, nil
	}
}
