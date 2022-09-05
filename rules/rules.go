package rules

import (
	"fmt"
	"log"
	"net/http"

	"strings"
)

func Rules(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Rules")
	requestURI := strings.Split(r.RequestURI, "/")
	log.Println(requestURI)
	log.Println(len(requestURI))
	log.Println(requestURI[0])
	log.Println(requestURI[1])

	if len(requestURI) >= 2 {
		log.Println(requestURI[2])
		if strings.EqualFold(requestURI[2], "allrisk") {

			Allrisk(w, r)
		}
	} else {
		fmt.Fprintf(w, "")
	}
	//lib.Files("")

}
