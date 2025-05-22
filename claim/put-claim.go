package claim

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func PutClaimFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("PutClaimFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	claim, err := models.UnmarshalClaim(body)
	lib.CheckError(err)

	return PutClaim(r.Header.Get("Authorization"), r.Header.Get("Origin"), &claim)
}

func PutClaim(idToken string, origin string, claim *models.Claim) (string, interface{}, error) {
	var (
		user models.User
		obj  mail.MailRequest
		att  []mail.Attachment
		err  error
	)
	log.AddPrefix("PutClaim")
	defer log.PopPrefix()

	userAuthID, err := lib.GetUserIdFromIdToken(idToken)
	if err != nil {
		log.ErrorF("invalid idToken, error %s", err.Error())
		return "", nil, err
	}

	fireUsers := lib.UserCollection
	docsnap, err := lib.GetFirestoreErr(fireUsers, userAuthID)
	if err != nil {
		log.ErrorF("get user from DB error %s", err.Error())
		return "", nil, err
	}
	err = docsnap.DataTo(&user)
	if err != nil {
		log.ErrorF("data to DB error %s", err.Error())
		return "", nil, err
	}
	log.Printf("User: %v", user)

	claim.CreationDate = time.Now().UTC()
	claim.Updated = claim.CreationDate
	claim.ClaimUid = uuid.New().String()
	claim.Status = "open"
	claim.StatusHistory = []string{"open"}

	for index, document := range claim.Documents {
		splitName := strings.Split(document.FileName, ".")
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		claim.Documents[index].Name = document.FileName
		claim.Documents[index].FileName = splitName[0] + "_" + timestamp + "." + splitName[1]
	}

	obj.From = "noreply@wopta.it"
	obj.To = []string{"sinistri@wopta.it"}
	obj.Title = claim.PolicyDescription + " n° " + claim.PolicyNumber
	obj.SubTitle = "Notifica sinistro"
	obj.Message = `<p>Ciao il cliente ` + user.Name + ` ` + user.
		Surname + `</p> <p>desidera notificare un sinistro per la polizza ` + claim.Company + ` n° ` +
		claim.PolicyNumber + ` per i seguenti motivi: ` + claim.Description + `</p> `
	obj.Subject = claim.PolicyDescription + " n° " + claim.PolicyNumber + " Notifica sinistro"
	obj.IsHtml = true
	if len(claim.Documents) > 0 {
		obj.IsAttachment = true
	}

	log.Println("uploading attachments to google storage")
	for i, doc := range claim.Documents {
		byteFile, err := base64.StdEncoding.DecodeString(doc.Byte)
		if err != nil {
			log.ErrorF("error decoding base64 document encoding")
			return "", nil, err
		}
		gsLink := lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+user.Uid+"/claims/"+
			claim.ClaimUid+"/"+doc.FileName, byteFile)
		att = append(att, mail.Attachment{Byte: doc.Byte, Name: strings.ReplaceAll(doc.Name, "_", " "), FileName: doc.FileName, ContentType: doc.ContentType})
		claim.Documents[i].Byte = ""
		claim.Documents[i].Link = gsLink
	}
	obj.Attachments = &att
	log.Println("attachments uploaded to google storage")

	if user.Claims == nil {
		user.Claims = new([]models.Claim)
	}
	*user.Claims = append(*user.Claims, *claim)

	log.Printf("update user %s on firestore", user.Uid)
	err = lib.UpdateFirestoreErr(fireUsers, user.Uid, map[string]interface{}{
		"claims":  user.Claims,
		"updated": time.Now().UTC(),
	})
	if err != nil {
		log.ErrorF("error during user update")
		return "", nil, err
	}

	mail.SendMail(obj)

	err = claim.BigquerySave(origin)
	if err != nil {
		log.ErrorF("error bigquery save claim %s", claim.ClaimUid)
	}

	return "{}", nil, nil
}
