package claim

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	model "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Claim")
	functions.HTTP("Claim", Claim)
}

func Claim(w http.ResponseWriter, r *http.Request) {
	log.Println("Claim")
	lib.EnableCors(&w, r)

	if r.Method == http.MethodGet {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		get(w, r)
	}
	if r.Method == http.MethodPut {
		w.Header().Set("Access-Control-Allow-Methods", "PUT")
		put(w, r)
	}
}
func get(w http.ResponseWriter, r *http.Request) {

	var user model.User
	docsnap := lib.GetFirestore("users", "")
	docsnap.DataTo(&user)

}
func post(w http.ResponseWriter, r *http.Request) {

}
func put(w http.ResponseWriter, r *http.Request) {
	var user model.User

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	defer r.Body.Close()
	claim, e := model.UnmarshalClaim(req)
	lib.CheckError(e)
	docsnap := lib.GetFirestore("users", claim.Uid)
	docsnap.DataTo(&user)
	claims := append(user.Claims, claim)
	user.Claims = claims
	lib.SetFirestore("users", claim.Uid, user)
	log.Println(user)

	// lib.PutFirestore("users")
}
