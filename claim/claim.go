package claim

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
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
	log.Println("Put")
	var user model.User

	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println(string(req))
	claim, e := model.UnmarshalClaim(req)
	lib.CheckError(e)
	log.Println("GetFirestore")
	docsnap := lib.GetFirestore("users", claim.Uid)
	docsnap.DataTo(&user)
	claim.CreationDate = time.Now().String()
	claim.Updated = time.Now().String()
	claims := append(user.Claims, claim)
	user.Claims = claims
	log.Println("SetFirestore")
	lib.SetFirestore("users", claim.Uid, user)

	log.Println(user)
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{"sinistri@wopta.it"}
	obj.Message = `<p>ciao ` + claim.Name + ` ` + claim.Surname + `</p> <p>desidera notificare un sinistro per la polizza: ` + claim.PolicyId + ` per i seguenti motivi: ` + claim.Description + `</p> `
	obj.Subject = "Notifica sinisto " + claim.PolicyId
	obj.IsHtml = true
	obj.IsAttachment = true
	var att []mail.Attachment
	for _, doc := range claim.Documents {
		att = append(att, mail.Attachment{Byte: doc.Byte})
	}
	obj.Attachments = att
	mail.SendMail(obj)
	// lib.PutFirestore("users")
}
