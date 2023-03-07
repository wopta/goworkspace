package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	models "github.com/wopta/goworkspace/models"
)

func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Proposal")
	var policy models.Policy
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	e := json.Unmarshal([]byte(req), &policy)

	defer r.Body.Close()
	policy.Updated = time.Now()
	policy.CreationDate = time.Now()
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusInitLead)
	policy.Status = models.PolicyStatusInitLead
	numb := GetSequenceProposal("")
	policy.ProposalNumber = numb
	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = time.Now()
	log.Println("Proposal save")
	ref, _ := lib.PutFirestore("policy", policy)
	ref, _ = lib.PutFirestore("users", policy.Contractor)
	policy.BigStartDate = civil.DateTimeOf(policy.StartDate)
	policy.BigEndDate = civil.DateTimeOf(policy.EndDate)
	e = lib.InsertRowsBigQuery("wopta", "policy", policy)
	log.Println(ref.ID + " Proposal sand mail")
	mail.SendMail(getProposalMailObj(policy))
	log.Println(ref.ID)

	return `{"uid":"` + ref.ID + `"}`, policy, e
}

func getProposalMailObj(policy models.Policy) mail.MailRequest {
	//att1 := lib.GetFromStorage("function-data", "information-sets/"+policy.Name+"/v1/CGA.pdf", "")
	//att2 := lib.GetFromStorage("function-data", "information-sets/"+policy.Name+"/v1/DIP.pdf", "")
	link := "https://storage.googleapis.com/documents-public-dev/information-set/information-sets/" + policy.Name + "/v1/Precontrattuale.pdf"
	var obj mail.MailRequest
	//att1bs64 := b64.StdEncoding.EncodeToString([]byte(att1))
	//att2bs64 := b64.StdEncoding.EncodeToString([]byte(att2))
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>Ciao ` + policy.Contractor.Name + ` ` + policy.Contractor.Surname + ` </p>

	<p>Leggi il set informativo precontrattuale a tua disposizione nel link qui sotto </p> 
	<p><a class="button" href='` + link + ` '>Leggi set informativo</a></p> 	<p>Grazie </p> <p></p>`
	obj.Subject = " Wopta set informantivo"
	obj.IsHtml = true
	obj.IsAttachment = false
	/*obj.Attachments = []mail.Attachment{
	{
		Byte:        att1bs64,
		Name:        "CGA.pdf",
		ContentType: "application/pdf",
	},
	{
		Byte:        att2bs64,
		Name:        "DIP.pdf",
		ContentType: "application/pdf",
	}}*/

	return obj
}
