package broker

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func ProposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[ProposalFx] Handler start -----------------------------------------")

	var (
		err    error
		policy models.Policy
	)

	origin = r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[ProposalFx] Request: %s", string(body))
	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.Printf("[ProposalFx] error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	err = Proposal(&policy)
	if err != nil {
		log.Printf("[ProposalFx] error creating proposal: %s", err.Error())
		return "", nil, err
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.Printf("[ProposalFx] error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Printf("[ProposalFx] response: %s", string(resp))

	return string(resp), &policy, err
}

func Proposal(policy *models.Policy) error {
	log.Println("[Proposal] start --------------------------------------------")

	var (
		err error
	)

	policyFire := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	guaranteFire := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

	log.Println("[Proposal] starting bpmn flow...")
	state := runBrokerBpmn(policy, proposalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[Proposal] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	*policy = *state.Data

	log.Println("[Proposal] saving proposal to firestore...")
	policyUid := lib.NewDoc(policyFire)
	policy.Uid = policyUid
	err = lib.SetFirestoreErr(policyFire, policyUid, policy)
	lib.CheckError(err)

	log.Println("[Proposal] saving proposal to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[Proposal] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "proposal", guaranteFire)

	log.Println("[Proposal] end ----------------------------------------------")
	return err
}

func setProposalData(policy *models.Policy) {
	log.Println("[setProposalData]")

	policyFire := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	now := time.Now().UTC()

	policy.CreationDate = now
	policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	policy.Status = models.PolicyStatusInitLead
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)

	numb := GetSequenceProposal("", policyFire)
	policy.ProposalNumber = numb
	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = now

	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}

	// TODO delete me when PMI is fixed
	if policy.Name == models.PmiProduct {
		policy.NameDesc = "Wopta per te Artigiani & Imprese"
	}

	policy.Attachments = &[]models.Attachment{{
		Name: "Precontrattuale", FileName: "Precontrattuale.pdf",
		Link: "gs://documents-public-dev/information-sets/" + policy.Name + "/" + policy.ProductVersion + "/Precontrattuale.pdf",
	}}

	log.Printf("[setProposalData] proposal number: %d", policy.ProposalNumber)
}
