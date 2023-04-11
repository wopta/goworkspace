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

func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Proposal")
	log.Println("--------------------------Proposal-------------------------------------------")
	var (
		policy  models.Policy
		useruid string
	)
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
	docsnap := lib.WhereFirestore("users", "fiscalCode", "==", policy.Contractor.FiscalCode)
	user, _ := models.FirestoreDocumentToUser(docsnap)

	if len(user.Uid) == 0 {
		ref2, _ := lib.PutFirestore("users", policy.Contractor)
		log.Println("Proposal User uid", ref2)
		useruid = ref2.ID
	} else {
		useruid = user.Uid
	}
	log.Println("Proposal User uid ", useruid)
	policy.Contractor.Uid = useruid
	//Precontrattuale.pdf
	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}
	policy.Attachments = &[]models.Attachment{{Name: "Precontrattuale", Link: "gs://documents-public-dev/information-sets/" + policy.Name + "/" + policy.ProductVersion + "v1/Precontrattuale.pdf"}}
	log.Println("Proposal save")
	ref, _ := lib.PutFirestore("policy", policy)
	policy.BigquerySave()
	log.Println(ref.ID + " Proposal sand mail")
	mail.SendMailProposal(policy)

	log.Println("Proposal ", ref.ID)

	return `{"uid":"` + ref.ID + `"}`, policy, e
}
