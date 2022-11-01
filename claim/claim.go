package claim

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("Claim", Claim)
}

func Claim(w http.ResponseWriter, r *http.Request) {
	log.Println("Document")
	lib.EnableCors(&w, r)

	if r.Method == http.MethodGet {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		get()
	}
	if r.Method == http.MethodPut {
		w.Header().Set("Access-Control-Allow-Methods", "PUT")
		post()
	}
}
func get() {

	var user model.User
	docsnap := lib.GetFirestore("", "")
	docsnap.DataTo(&user)

}
func post() {

}
