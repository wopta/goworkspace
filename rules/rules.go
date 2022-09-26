package rules

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	base := "/rules"
	if strings.Contains(r.RequestURI, "/rules") {
		base = "/rules"
	} else {
		base = ""
	}
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	switch os := r.RequestURI; os {
	case base + "/allrisk":
		Allrisk(w, r)
	case base + "/pmi-allrisk":
		PmiAllrisk(w, r)
	default:
		fmt.Fprintf(w, "")
	}

	//lib.Files("")

}
