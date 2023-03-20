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
	log.Println("--------------------------Proposal-------------------------------------------")
	var policy models.Policy
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	e := json.Unmarshal([]byte(req), &policy)
	j, e := policy.Marshal()
	log.Println("Proposal request proposal: ", string(j))
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
	ref2, _ := lib.PutFirestore("users", policy.Contractor)
	log.Println("Proposal User uid", ref2)
	policy.BigStartDate = civil.DateTimeOf(policy.StartDate)
	policy.BigEndDate = civil.DateTimeOf(policy.EndDate)

	policy.BigEmitDate = civil.DateTimeOf(policy.Updated)
	policyJson, e := policy.Marshal()
	policy.Data = string(policyJson)
	e = lib.InsertRowsBigQuery("wopta", "policy", policy)
	log.Println(ref.ID + " Proposal sand mail")
	mail.SendMail(getProposalMailObj(policy))
	log.Println("Proposal ", ref.ID)

	return `{"uid":"` + ref.ID + `"}`, policy, e
}

func getProposalMailObj(policy models.Policy) mail.MailRequest {
	var name string
	var linkForm string
	if policy.Name == "pmi" {
		name = "Artigiani & Imprese"
		linkForm = "https://www.wopta.it/it/multi-rischio/"
	}

	link := "https://storage.googleapis.com/documents-public-dev/information-set/information-sets/" + policy.Name + "/v1/Precontrattuale.pdf"
	var obj mail.MailRequest

	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p></p><p>Gentile ` + policy.Contractor.Name + ` ` + policy.Contractor.Surname + ` </p>
	<p></p>
	<p>richiedendo un preventivo per la soluzione assicurativa Wopta per Te ` + name + ` , dimostri interesse nel proteggere la tua Attività. </p> 
	<p>Per poter valutare completamente la soluzione che sceglierai, ti alleghiamo tutti i documenti che ti consentiranno di prendere una decisione pienamente consapevole ed informata.</p>
	<p>Prima della sottoscrizione, leggi quanto trovi in questo <a class="button" href='` + link + ` '>Link</a></p>
	<p>Un saluto.</p>
	<p>ll Team Wopta. Proteggiamo chi sei</p> 	
	<p></p>
	<p></p>
	<p>Se hai bisogno di ulteriore supporto, non scrivere a questo indirizzo email, puoi compilare il <a class="button" href='` + linkForm + ` '>Form</a> oppure scrivere alla mail e verrai contattato da un nostro esperto.</p>
	<p></p>
	`
	obj.Subject = "Wopta per Te " + name + " Documenti informativi precontrattuali"
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
