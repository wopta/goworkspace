package callback

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

/**
workstepFinished : when the workstep was finished
workstepRejected : when the workstep was rejected
workstepDelegated : whe the workstep was delegated
workstepOpened : when the workstep was opened
sendSignNotification : when the sign notification was sent
envelopeExpired : when the envelope was expired
workstepDelegatedSenderActionRequired : when an action from the sender is required because of the delegation
*/
func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())
	var e error
	uid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")

	log.Println(action)
	log.Println(envelope)
	log.Println(uid)
	if action == "workstepFinished" {
		log.Println("workstepFinished")
		policyF := lib.GetFirestore("policy", uid)
		var policy models.Policy
		policyF.DataTo(&policy)
		policy.IsSign = true
		policy.Updated = time.Now()
		policy.Status = models.PolicyStatusToPay
		policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
		policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
		log.Println("workstepFinished")
		lib.SetFirestore("policy", uid, policy)
		e = lib.InsertRowsBigQuery("wopta", "policy", policy)
		mail.SendMail(getEmitMailObj(policy, policy.PayUrl))
		log.Println("workstepFinished")
		s := <-GetFileV6(policy.IdSign, uid)
		log.Println(s)
	}

	return "", nil, e
}
func getEmitMailObj(policy models.Policy, payUrl string) mail.MailRequest {
	var obj mail.MailRequest
	log.Println(policy.Contractor.Mail)
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>Ciao ` + policy.Contractor.Name + `` + policy.Contractor.Surname + ` </p>
	<p>Polizza n° ` + policy.NumberCompany + `</p> 
	<p>Grazie per aver scelto uno dei nostri prodotti Wopta per te</p> 
	<p>Puoi ora procedere alla firma della polizza in oggetto. Qui trovi il link per
	 accedere alla procedura semplice e guidata di firma elettronica avanzata tramite utilizzo di
	  un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato. 
	Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. 
	Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.</p> 
	<p><a class="button" href='` + payUrl + `'>Firma la tua polizza:</a></p>
	<p>Ultimata la procedura di firma potrai procedere al pagamento.</p>
	<p>Grazie per aver scelto Wopta </p> 
	<p>Proteggiamo chi sei</p> 
	`
	obj.Subject = " Wopta Contratto e pagamento"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
func PutToStorage(bucketname string, path string, file []byte) (string, error) {

	log.Println("start PutToStorage")
	ctx := context.Background()
	client, e := storage.NewClient(ctx)
	bucket := client.Bucket(bucketname)
	write := bucket.Object(path).NewWriter(ctx)
	defer write.Close()
	write.Write(file)

	return "gs://" + bucketname + "/" + path, e

}
func GetFileV6(id string, uid string) chan string {
	r := make(chan string)
	log.Println("Get file: ", id)
	go func() {

		defer close(r)
		files := <-GetFilesV6(id)

		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/file/" + files.Documents[0].FileID
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
		log.Println("url parse:", req.Header)
		res, err := client.Do(req)
		lib.CheckError(err)
		if res != nil {
			body, _ := ioutil.ReadAll(res.Body)

			//log.Println("Get body: ", string(body))
			_, e := PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "document/contracts/"+uid, body)
			lib.CheckError(e)
			defer res.Body.Close()
			r <- "upload done"

		}

	}()
	return r
}
func GetFilesV6(envelopeId string) chan NamirialFiles {
	r := make(chan NamirialFiles)

	go func() {
		defer close(r)
		var urlstring = os.Getenv("ESIGN_BASEURL") + "v6/envelope/" + envelopeId + "/files"
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest(http.MethodGet, urlstring, nil)
		req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))

		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {

			body, _ := ioutil.ReadAll(res.Body)
			resp, _ := UnmarshalNamirialFiles(body)
			res.Body.Close()

			log.Println("body:", string(body))
			r <- resp
		}

	}()
	return r
}
func UnmarshalNamirialFiles(data []byte) (NamirialFiles, error) {
	var r NamirialFiles
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NamirialFiles) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NamirialFiles struct {
	Documents      []Documents     `json:"Documents"`
	AuditTrail     AuditTrail      `json:"AuditTrail"`
	LegalDocuments []LegalDocument `json:"LegalDocuments"`
}

type AuditTrail struct {
	FileID    string `json:"FileId"`
	XMLFileID string `json:"XmlFileId"`
}

type Documents struct {
	FileID           string       `json:"FileId"`
	FileName         string       `json:"FileName"`
	AuditTrailFileID string       `json:"AuditTrailFileId"`
	Attachments      []Attachment `json:"Attachments"`
	PageCount        int64        `json:"PageCount"`
	DocumentNumber   int64        `json:"DocumentNumber"`
}

type Attachment struct {
	FileID   string `json:"FileId"`
	FileName string `json:"FileName"`
}

type LegalDocument struct {
	FileID     string `json:"FileId"`
	FileName   string `json:"FileName"`
	ActivityID string `json:"ActivityId"`
	Email      string `json:"Email"`
}
