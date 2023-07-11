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
		policy     models.Policy
		policyFire string
	)
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	e := json.Unmarshal([]byte(req), &policy)
	j, e := policy.Marshal()
	log.Println("Proposal request proposal: ", string(j))
	defer r.Body.Close()
	policyFire = lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	guaranteFire := lib.GetDatasetByEnv(r.Header.Get("origin"), "guarante")
	policy.CreationDate = time.Now().UTC()
	policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusInitLead)
	policy.Status = models.PolicyStatusInitLead
	numb := GetSequenceProposal("", policyFire)
	policy.ProposalNumber = numb
	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = time.Now()

	//------------------------------------------
	//Precontrattuale.pdf
	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}
	policy.Attachments = &[]models.Attachment{
		{
			Name: "Precontrattuale", FileName: "Precontrattuale.pdf",
			Link: "gs://documents-public-dev/information-sets/" + policy.Name + "/" + policy.ProductVersion + "/Precontrattuale.pdf",
		},
	}
	log.Println("Proposal save")
	ref, _ := lib.PutFirestore(policyFire, policy)
	policy.BigquerySave(r.Header.Get("origin"))
	models.SetGuaranteBigquery(policy, "proposal", guaranteFire)
	log.Println(ref.ID + " Proposal sand mail")
	mail.SendMailProposal(policy)

	log.Println("Proposal ", ref.ID)

	policy.Uid = ref.ID

	resp, e := policy.Marshal()

	return string(resp), policy, e
}
