package control

import (
	"context"
	"log"
	"net/http"

	"goeasy.dev"
	"goeasy.dev/status"
)

func NewRunner() goeasy.RunnerFunc {
	mux := http.NewServeMux()

	mux.Handle("/status/", http.StripPrefix("/status", status.Handler()))

	return func(ctx context.Context) (goeasy.StopFunc, error) {
		srv := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		go func() {
			log.Println("starting control server")
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Printf("error in controll server: %v\n", err)
			}
		}()

		return func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		}, nil
	}
}
