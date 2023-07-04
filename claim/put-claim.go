package claim

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

func PutClaimFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PutClaimFx]")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Println("[PutClaimFx] " + string(body))

	claim, err := models.UnmarshalClaim(body)
	lib.CheckError(err)

	return PutClaim(r.Header.Get("Authorization"), r.Header.Get("Origin"), &claim)
}

func PutClaim(idToken string, origin string, claim *models.Claim) (string, interface{}, error) {
	var (
		user   models.User
		obj    mail.MailRequest
		att    []mail.Attachment
		err    error
		policy models.Policy
	)

	userAuthID, err := lib.GetUserIdFromIdToken(idToken)
	if err != nil {
		log.Printf("[PutClaim] invalid idToken, error %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	fireUsers := lib.GetDatasetByEnv(origin, models.UserCollection)
	docsnap, err := lib.GetFirestoreErr(fireUsers, userAuthID)
	if err != nil {
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		return `{"success":false}`, `{"success":false}`, nil
	}
	log.Println("[PutClaim] User: ", user)

	log.Printf("[PutClaim] get policy %s from firestore", claim.PolicyId)
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	docsnap, err = lib.GetFirestoreErr(firePolicy, claim.PolicyId)
	if err != nil {
		log.Printf("[PutClaim] error retrieving policy %s from firestore, error message: %s", claim.PolicyId, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&policy)
	if err != nil {
		log.Println("[PutClaim] error convert docsnap to policy")
		return `{"success":false}`, `{"success":false}`, nil
	}

	if userAuthID != policy.Contractor.Uid {
		log.Println("[PutClaim] claim requester and policy contractor are not the same")
		return `{"success":false}`, `{"success":false}`, nil
	}

	claim.CreationDate = time.Now().UTC()
	claim.Updated = claim.CreationDate
	claim.ClaimUid = uuid.New().String()
	claim.Status = "open"

	for index, document := range claim.Documents {
		splitName := strings.Split(document.FileName, ".")
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		claim.Documents[index].Name = document.FileName
		claim.Documents[index].FileName = splitName[0] + "_" + timestamp + "." + splitName[1]
	}

	obj.From = "noreply@wopta.it"
	obj.To = []string{"sinistri@wopta.it"}
	obj.Message = `<p>Ciao il cliente ` + user.Name + ` ` + user.Surname + `</p> <p>desidera notificare un sinistro per la polizza: ` + policy.CodeCompany + ` per i seguenti motivi: ` + claim.Description + `</p> `
	obj.Subject = "Notifica sinisto polizza " + policy.CodeCompany
	obj.IsHtml = true
	if len(claim.Documents) > 0 {
		obj.IsAttachment = true
	}

	log.Println("[PutClaim] uploading attachments to google storage")
	for i, doc := range claim.Documents {
		byteFile, err := base64.StdEncoding.DecodeString(doc.Byte)
		if err != nil {
			log.Println("[PutClaim] error decoding base64 document encoding")
			return `{"success":false}`, `{"success":false}`, nil
		}
		gsLink := lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+userAuthID+"/claims/"+
			claim.ClaimUid+"/"+doc.FileName, byteFile)
		att = append(att, mail.Attachment{Byte: doc.Byte, Name: doc.Name, FileName: doc.FileName, ContentType: doc.ContentType})
		claim.Documents[i].Byte = ""
		claim.Documents[i].Link = gsLink
	}
	obj.Attachments = &att
	log.Println("[PutClaim] attachments uploaded to google storage")

	if user.Claims == nil {
		user.Claims = new([]models.Claim)
	}
	*user.Claims = append(*user.Claims, *claim)

	log.Printf("[PutClaim] update user %s on firestore", userAuthID)
	err = lib.UpdateFirestoreErr(fireUsers, userAuthID, map[string]interface{}{
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
