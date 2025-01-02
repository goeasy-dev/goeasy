package control

import (
	"net/http"

	"goeasy.dev"
	"goeasy.dev/application/runners"
	"goeasy.dev/internal/control/status"
)

func NewRunner() goeasy.RunnerFunc {
	mux := http.NewServeMux()

	mux.Handle("/status/", http.StripPrefix("/status", status.Handler()))

	return runners.Http(runners.HttpOptions{
		Address: ":8081",
		Handler: mux,
	})
}
