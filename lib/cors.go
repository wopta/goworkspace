package lib

import "net/http"

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST")
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type , Authorization")
		(*w).Header().Set("Access-Control-Max-Age", "3600")
		(*w).WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}
