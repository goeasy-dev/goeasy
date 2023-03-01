package control

import (
	"net/http"

	"goeasy.dev/status"
)

func NewRunner() http.HandlerFunc {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", status.HandlerFunc())

	return mux.ServeHTTP
}
