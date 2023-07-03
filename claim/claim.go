package claim

import (
	b64 "encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
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
	if r.Method == http.MethodPost {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
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
	docsnap := lib.GetFirestore("users", claim.UserUid)
	e = docsnap.DataTo(&user)
	lib.CheckError(e)
	claim.CreationDate = time.Now()
	claim.Updated = time.Now()
	uidClaim := uuid.New().String()
	claim.ClaimUid = uidClaim
	log.Println(user)
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{"sinistri@wopta.it"}
	obj.Message = `<p>ciao il cliente ` + claim.Name + ` ` + claim.Surname + `</p> <p>desidera notificare un sinistro per la polizza: ` + claim.PolicyId + ` per i seguenti motivi: ` + claim.Description + `</p> `
	obj.Subject = "Notifica sinisto " + claim.PolicyId
	obj.IsHtml = true
	if len(claim.Documents) > 0 {
		obj.IsAttachment = true
	}
	var att []mail.Attachment
	for i, doc := range claim.Documents {
		byteFile, e := b64.StdEncoding.DecodeString(doc.Byte)
		lib.CheckError(e)
		link := lib.PutToStorage(os.Getenv("USER_BUCKET"), "users/"+claim.UserUid+"/claims/"+uidClaim+"/"+doc.FileName, byteFile)
		att = append(att, mail.Attachment{Byte: doc.Byte, Name: doc.FileName, ContentType: doc.ContentType})
		claim.Documents[i].Byte = ""
		claim.Documents[i].Link = link
	}
	obj.Attachments = &att

	userClaims := make([]models.Claim, 0)
	if user.Claims != nil {
		userClaims = append(userClaims, *user.Claims...)
	}
	userClaims = append(userClaims, claim)

	user.Claims = &userClaims
	log.Println("SetFirestore")
	lib.UpdateFirestoreErr("users", claim.UserUid, map[string]interface{}{
		"claims": userClaims,
	})
	mail.SendMail(obj)
	// lib.PutFirestore("users")
}
