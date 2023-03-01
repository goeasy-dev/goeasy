package status

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlerFunc() http.HandlerFunc {
	mux := http.NewServeMux()

	mux.HandleFunc("/liveness", getHandlerFunc(Liveness))
	mux.HandleFunc("/readiness", getHandlerFunc(Readiness))
	mux.HandleFunc("/startup", getHandlerFunc(Startup))

	return mux.ServeHTTP
}

type statusResult struct {
	Status map[string]bool `json:"status"`
}

func getHandlerFunc(kind CheckType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodGet || r.Method == http.MethodHead) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		status, ok := CheckStatus(r.Context(), kind)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Add("content-type", "application/json")

		if r.Method == http.MethodHead {
			return
		}

		encoded, err := json.Marshal(statusResult{Status: status})
		if err != nil {
			log.Println("unable to marshal status json: %w", err)
			return
		}

		w.Write(encoded)
	}
}
