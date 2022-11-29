package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	models "github.com/wopta/goworkspace/models"
)

func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("Proposal")
	var policy models.Policy
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))

	e := json.Unmarshal([]byte(req), &policy)
	lib.CheckError(e)
	defer r.Body.Close()
	//policy, e := models.UnmarshalPolicy(req)
	policy.Updated = time.Now()
	policy.CreationDate = time.Now()
	policy.Status = models.Proposal
	numb := GetSequenceProposal("global")
	policy.ProposalNumber = numb
	log.Println("save")
	ref, _ := lib.PutFirestore("policy", policy)
	log.Println("saved")
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>ciao </p> `
	obj.Subject = "Wopta Proposta e set informantivo"
	obj.IsHtml = true
	mail.SendMail(obj)
	log.Println(ref.ID)

	return `{"uid":"` + ref.ID + `"}`, policy
}
