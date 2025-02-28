package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LeadFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.SetPrefix("[LeadFx]")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	authToken := r.Context().Value(lib.CtxAuthToken).(lib.AuthToken)
	nn := r.Context().Value(models.CtxRequesterNetworkNode).(models.NetworkNode)
	networkNode = &nn
	// token := r.Header.Get("Authorization")
	// authToken, err := lib.GetAuthTokenFromIdToken(token)
	// if err != nil {
	// 	log.Printf("error getting authToken")
	// 	return "", nil, err
	// }
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.Println("error decoding body")
		return "", nil, err
	}
	// origin = r.Header.Get("Origin")
	// body := lib.ErrorByte(io.ReadAll(r.Body))
	// defer r.Body.Close()

	// err = json.Unmarshal([]byte(body), &policy)
	// if err != nil {
	// 	log.Printf("error unmarshaling policy: %s", err.Error())
	// 	return "", nil, err
	// }

	policy.Normalize()

	err = lead(authToken, &policy)
	if err != nil {
		log.Printf("error creating lead: %s", err.Error())
		return "", nil, err
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(resp), &policy, err
}

func lead(authToken models.AuthToken, policy *models.Policy) error {
	var err error

	log.Println("[lead] start ------------------------------------------------")

	if policy.Uid != "" {
		if err = checkIfPolicyIsLead(policy); err != nil {
			return err
		}
	} else {
		policy.Uid = lib.NewDoc(lib.PolicyCollection)
	}

	if policy.Channel == "" {
		policy.Channel = authToken.GetChannelByRoleV2()
		log.Printf("[lead] setting policy channel to '%s'", policy.Channel)
	}

	// networkNode = network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("[lead] starting bpmn flow...")
	state := runBrokerBpmn(policy, leadFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[lead] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	*policy = *state.Data

	log.Println("[lead] saving lead to firestore...")
	err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
	lib.CheckError(err)

	log.Println("[lead] saving lead to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[lead] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "lead", lib.GuaranteeCollection)

	log.Println("[lead] end --------------------------------------------------")
	return err
}

func checkIfPolicyIsLead(policy *models.Policy) error {
	var recoveredPolicy models.Policy
	policyDoc, err := lib.GetFirestoreErr(lib.PolicyCollection, policy.Uid)
	if err != nil {
		log.Printf("[checkIfPolicyIsLead] error getting policy %s from firebase: %s", policy.Uid, err.Error())
		return nil
	} else if !policyDoc.Exists() {
		log.Printf("[checkIfPolicyIsLead] policy %s not found on Firebase", policy.Uid)
		return nil
	}

	if err := policyDoc.DataTo(&recoveredPolicy); err != nil {
		log.Printf("[checkIfPolicyIsLead] error converting policy %s data: %s", policy.Uid, err.Error())
		return nil
	}

	allowedStatus := []string{models.PolicyStatusInit, models.PolicyStatusPartnershipLead, models.PolicyStatusInitLead}
	if !slices.Contains(allowedStatus, recoveredPolicy.Status) {
		log.Printf("[checkIfPolicyIsLead] error policy %s is not a lead", policy.Uid)
		return errors.New("policy is not a lead")
	}

	log.Printf("[checkIfPolicyIsLead] found lead for existing policy %s", policy.Uid)

	policy.CreationDate = recoveredPolicy.CreationDate
	policy.Status = recoveredPolicy.Status
	policy.StatusHistory = recoveredPolicy.StatusHistory
	policy.ProducerUid = recoveredPolicy.ProducerUid
	policy.ProducerCode = recoveredPolicy.ProducerCode
	policy.ProducerType = recoveredPolicy.ProducerType
	policy.NetworkUid = recoveredPolicy.NetworkUid

	return nil
}

func setLeadData(policy *models.Policy, product models.Product) {
	log.Println("[setLeadData] start -----------------------------------------")

	now := time.Now().UTC()

	if policy.CreationDate.IsZero() {
		policy.CreationDate = now
	}
	if policy.Status != models.PolicyStatusInitLead {
		policy.Status = models.PolicyStatusInitLead
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}
	log.Printf("[setLeadData] policy status %s", policy.Status)

	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = now

	if networkNode != nil {
		setPolicyProducerNode(policy, networkNode)
	}

	// TODO delete me when PMI is fixed
	if policy.Name == models.PmiProduct {
		policy.NameDesc = "Wopta per te Artigiani & Imprese"
	}
	if policy.ProductVersion == "" {
		policy.ProductVersion = "v1"
	}

	setRenewInfo(policy, product)

	log.Println("[setLeadData] add information set")
	informationSet := models.Attachment{
		Name:     "Precontrattuale",
		FileName: "Precontrattuale.pdf",
		Link: fmt.Sprintf(
			"gs://documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf",
			policy.Name,
			policy.ProductVersion,
		),
	}
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	attIdx := slices.IndexFunc(*policy.Attachments, func(a models.Attachment) bool {
		return a.Name == informationSet.Name
	})
	if attIdx == -1 {
		*policy.Attachments = append(*policy.Attachments, informationSet)
	}

	log.Println("[setLeadData] end -------------------------------------------")
}

func setPolicyProducerNode(policy *models.Policy, node *models.NetworkNode) {
	policy.ProducerUid = node.Uid
	policy.ProducerCode = node.Code
	policy.ProducerType = node.Type
	policy.NetworkUid = node.NetworkUid
}

func setRenewInfo(policy *models.Policy, product models.Product) {
	policy.Annuity = 0
	policy.IsRenewable = product.IsRenewable
	policy.IsAutoRenew = product.IsAutoRenew
	policy.PolicyType = product.PolicyType
	policy.QuoteType = product.QuoteType
}
