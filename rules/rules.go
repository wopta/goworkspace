package rules

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
)

func Rules(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	log.Println("Rules")

	if strings.Contains(r.RequestURI, "/allrisk") {

		Allrisk(w, r)

	} else {
		fmt.Fprintf(w, "")
	}
	//lib.Files("")

}
