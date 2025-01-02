package status

import (
	"encoding/json"
	"log"
	"net/http"

	"goeasy.dev/status"
	"goeasy.dev/status/statustype"
)

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/liveness", getHandlerFunc(statustype.Liveness))
	mux.HandleFunc("/readiness", getHandlerFunc(statustype.Readiness))
	mux.HandleFunc("/startup", getHandlerFunc(statustype.Startup))

	return mux
}

type statusResult struct {
	Status map[string]bool `json:"status"`
}

func getHandlerFunc(kind statustype.Type) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodGet || r.Method == http.MethodHead) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("content-type", "application/json")

		status, ok := status.CheckStatus(r.Context(), kind)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

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
