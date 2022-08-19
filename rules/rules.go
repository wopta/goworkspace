package rules

import (
	"log"
	"net/http"
	"strings"
)

func Rules(w http.ResponseWriter, r *http.Request) {
	log.Println("Rules")
	requestURI := strings.Split(r.RequestURI, "/")
	log.Println(requestURI)
	log.Println(len(requestURI))
	log.Println("newwww")
	log.Println(requestURI[0])
	log.Println(requestURI[1])

	if len(requestURI) >= 2 {
		log.Println(requestURI[2])
		if strings.EqualFold(requestURI[2], "allrisk") {

			Allrisk(w, r)
		}
	}
	//lib.Files("")

}
