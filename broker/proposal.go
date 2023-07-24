package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	models "github.com/wopta/goworkspace/models"
)

func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[Proposal] Handler start -----------------------------------------")

	var (
		policy     models.Policy
		policyFire string
		origin     string = r.Header.Get("origin")
	)

	req := lib.ErrorByte(io.ReadAll(r.Body))
	e := json.Unmarshal([]byte(req), &policy)
	j, e := policy.Marshal()
	log.Println("[Proposal] request body: ", string(j))
	defer r.Body.Close()

	policyFire = lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")

	policy.CreationDate = time.Now().UTC()
	policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusInitLead)
	policy.Status = models.PolicyStatusInitLead
	numb := GetSequenceProposal("", policyFire)
	policy.ProposalNumber = numb
	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = time.Now().UTC()

	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}
	policy.Attachments = &[]models.Attachment{{
		Name: "Precontrattuale", FileName: "Precontrattuale.pdf",
		Link: "gs://documents-public-dev/information-sets/" + policy.Name + "/" + policy.ProductVersion + "/Precontrattuale.pdf",
	}}

	log.Println("[Proposal] save")
	policyUid := lib.NewDoc(policyFire)
	policy.Uid = policyUid
	err := lib.SetFirestoreErr(policyFire, policyUid, policy)
	lib.CheckError(err)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(policy, "proposal", guaranteFire)

	log.Printf("[Proposal] Policy %s send mail", policy.Uid)
	mail.SendMailProposal(&policy)

	resp, e := policy.Marshal()
	log.Println("[Proposal] response: ", string(resp))

	return string(resp), policy, e
}
