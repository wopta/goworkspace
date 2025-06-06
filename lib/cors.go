package lib

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		log.Println("---------------http.MethodOptions OPTION----------------------------------------------------------------")
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "*")
		(*w).Header().Set("Access-Control-Allow-Headers", "*")
		(*w).Header().Set("Access-Control-Allow-Credentials", "true")
		(*w).Header().Set("Access-Control-Max-Age", "3600")
		(*w).WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "false")

	// Only for local development
	if env.IsLocal() {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Content-Type", "application/json")
	}
}
