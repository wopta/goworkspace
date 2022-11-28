package callback

import (
	"log"
	"net/http"
)

func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())

	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	log.Println(action)
	log.Println(envelope)

	if action == "workstepFinished" {

	}

	return "", nil
}
