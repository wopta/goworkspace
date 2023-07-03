package claim

import (
	b64 "encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Claim")
	functions.HTTP("Claim", Claim)
}

func Claim(w http.ResponseWriter, r *http.Request) {
	log.Println("Claim")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "",
				Handler: PutClaimFx,
				Method:  http.MethodPut,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}

/*func get(w http.ResponseWriter, r *http.Request) {

	var user model.User
	docsnap := lib.GetFirestore("users", "")
	docsnap.DataTo(&user)

}
*/

func PutClaimFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("PutClaim")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println(string(body))

	claim, err := models.UnmarshalClaim(body)
	lib.CheckError(err)

	return PutClaim(r.Header.Get("Origin"), &claim)
}

func PutClaim(origin string, claim *models.Claim) (string, interface{}, error) {
	var (
		user models.User
		obj  mail.MailRequest
		att  []mail.Attachment
	)

	log.Printf("Get user %s from firestore", claim.UserUid)
	fireUsers := lib.GetDatasetByEnv(origin, "users")
	docsnap, err := lib.GetFirestoreErr(fireUsers, claim.UserUid)
	if err != nil {
		log.Printf("[PutClaim] error retrieving user %s from firestore, error message: %s", claim.UserUid, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		log.Println("[PutClaim] error convert docsnap to user")
		return `{"success":false}`, `{"success":false}`, nil
	}

	claim.CreationDate = time.Now().UTC()
	claim.Updated = time.Now().UTC()
	claim.ClaimUid = uuid.New().String()
	claim.Status = "open"

	log.Println("User: ", user)

	obj.From = "noreply@wopta.it"
	obj.To = []string{"sinistri@wopta.it"}
	obj.Message = `<p>ciao il cliente ` + claim.Name + ` ` + claim.Surname + `</p> <p>desidera notificare un sinistro per la polizza: ` + claim.PolicyId + ` per i seguenti motivi: ` + claim.Description + `</p> `
	obj.Subject = "Notifica sinisto " + claim.PolicyId
	obj.IsHtml = true
	if len(claim.Documents) > 0 {
		obj.IsAttachment = true
	}

	log.Println("[PutClaim] uploading attachments to google storage")
	for i, doc := range claim.Documents {
		byteFile, err := b64.StdEncoding.DecodeString(doc.Byte)
		if err != nil {
			log.Println("[PutClaim] error decoding base64 document encoding")
			return `{"success":false}`, `{"success":false}`, nil
		}
		gsLink := lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+claim.UserUid+"/claims/"+
			claim.ClaimUid+"/"+doc.FileName, byteFile)
		att = append(att, mail.Attachment{Byte: doc.Byte, Name: doc.FileName, ContentType: doc.ContentType})
		claim.Documents[i].Byte = ""
		claim.Documents[i].Link = gsLink
	}
	obj.Attachments = &att
	log.Println("[PutClaim] attachments uploaded to google storage")

	if user.Claims == nil {
		*user.Claims = make([]models.Claim, 0)
	}
	*user.Claims = append(*user.Claims, *claim)

	log.Printf("[PutClaim] update user %s on firestore", claim.UserUid)
	err = lib.UpdateFirestoreErr(fireUsers, claim.UserUid, map[string]interface{}{
		"claims":  user.Claims,
		"updated": time.Now().UTC(),
	})
	if err != nil {
		log.Println("[PutClaim] error during user update")
		return `{"success":false}`, `{"success":false}`, nil
	}

	mail.SendMail(obj)

	return `{"success":true}`, `{"success":true}`, nil
}
