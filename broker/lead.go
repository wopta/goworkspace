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
	"github.com/wopta/goworkspace/network"
)

func LeadFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.Println("[LeadFx] Handler start -----------------------------------------")

	origin = r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Println("[LeadFx] loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[LeadFx] error getting authToken")
		return "", nil, err
	}

	log.Printf("[LeadFx] request: %s", string(body))
	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.Printf("[LeadFx] error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	err = lead(authToken, &policy)
	if err != nil {
		log.Printf("[LeadFx] error creating lead: %s", err.Error())
		return "", nil, err
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.Printf("[LeadFx] error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Printf("[LeadFx] response: %s", string(resp))

	return string(resp), &policy, err
}

func lead(authToken models.AuthToken, policy *models.Policy) error {
	var err error

	log.Println("[lead] start --------------------------------------------")

	policyFire := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	guaranteFire := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

	if policy.Uid != "" {
		log.Printf("[lead] creating lead for existing policy %s", policy.Uid)
		policyDoc, err := lib.GetFirestoreErr(models.PolicyCollection, policy.Uid)

		if err != nil {
			log.Printf("[lead] error getting policy %s from firebase: %s", policy.Uid, err.Error())
		}

		if  err == nil && policyDoc.Exists() {
			log.Printf("[lead] found policy %s on Firebase", policy.Uid)
			policyDoc.DataTo(policy)
		}
	} else {
		policy.Uid = lib.NewDoc(policyFire)
	}

	if policy.Channel == "" {
		policy.Channel = authToken.GetChannelByRole()
	}

	if policy.ProducerUid == "" {
		node := network.GetNetworkNodeByUid(authToken.UserID)
		if node != nil {
			setPolicyProducerNode(policy, node)
		}
	}

	log.Println("[lead] starting bpmn flow...")
	state := runBrokerBpmn(policy, leadFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[lead] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	*policy = *state.Data

	log.Println("[lead] saving lead to firestore...")
	err = lib.SetFirestoreErr(policyFire, policy.Uid, policy)
	lib.CheckError(err)

	log.Println("[lead] saving lead to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[lead] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "lead", guaranteFire)

	log.Println("[lead] end ----------------------------------------------")
	return err
}

func setPolicyProducerNode(policy *models.Policy, node *models.NetworkNode) {
	policy.ProducerUid = node.Uid
	policy.ProducerCode = node.Code
	policy.ProducerType = node.Type
	policy.NetworkUid = node.NetworkUid
}

func setLeadData(policy *models.Policy) {
	log.Println("[setLeadData] start -------------------")

	now := time.Now().UTC()

	if policy.CreationDate.IsZero() {
		policy.CreationDate = now
	}
	policy.Status = models.PolicyStatusInitLead
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	log.Printf("[setLeadData] policy status %s", policy.Status)

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

	log.Println("[setLeadData] add information stet")
	policy.Attachments = &[]models.Attachment{{
		Name: "Precontrattuale", FileName: "Precontrattuale.pdf",
		Link: "gs://documents-public-dev/information-sets/" + policy.Name + "/" + policy.ProductVersion + "/Precontrattuale.pdf",
	}}
	log.Println("[setLeadData] end -------------------")
}
