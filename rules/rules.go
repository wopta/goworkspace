package rules

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Rules")

	functions.HTTP("Rules", Rules)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	log.Println("Rules")
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)

	if strings.Contains(r.RequestURI, "/allrisk") {

		Allrisk(w, r)

	} else {
		fmt.Fprintf(w, "")
	}
	//lib.Files("")

}
