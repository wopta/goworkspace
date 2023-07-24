package callback

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

const fabrickBillPaid string = "PAID"

func PaymentV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PaymentV2Fx] Handler start -----------------------------------")

	var (
		response        string
		err             error
		fabrickCallback FabrickCallback
	)

	policyUid := r.URL.Query().Get("uid")
	schedule := r.URL.Query().Get("schedule")
	origin := r.URL.Query().Get("origin")
	log.Printf("[PaymentV2Fx] uid %s, schedule %s", policyUid, schedule)

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[PaymentV2Fx] request payload: %s", string(request))
	err = json.Unmarshal([]byte(request), &fabrickCallback)
	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR unmarshaling request: %s", err.Error())
		return response, nil, err
	}

	if policyUid == "" || origin == "" {
		ext := strings.Split(fabrickCallback.ExternalID, "_")
		policyUid = ext[0]
		schedule = ext[1]
		origin = ext[2]
	}

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		err = fabrickPayment(origin, policyUid, schedule)
	default:
	}

	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR: %s", err.Error())
		return response, nil, err
	}

	response = `{
		"result": true,
		"requestPayload": ` + string(request) + `,
		"locale": "it"
	}`
	log.Printf("[PaymentV2Fx] response: %s", response)

	return response, nil, nil
}

func fabrickPayment(origin, policyUid, schedule string) error {
	log.Printf("[fabrickPayment] Policy %s", policyUid)

	policy := plc.GetPolicyByUid(policyUid, origin)

	if !policy.IsPay && policy.Status == models.PolicyStatusToPay {
		// Create/Update document on user collection based on contractor fiscalCode
		user.SetUserIntoPolicyContractor(&policy, origin)

		// Get Policy contract
		gsLink := <-document.GetFileV6(policy, policyUid)
		log.Println("[fabrickPayment] contractGsLink: ", gsLink)

		// Update Policy as paid
		plc.SetPolicyPaid(&policy, gsLink, origin)

		// Update the first transaction in policy as paid
		transaction.SetPolicyFirstTransactionPaid(policyUid, schedule, origin)

		// Update agency if present
		updateAgencyPortfolio(&policy, origin)

		// Update agent if present
		updateAgentPortfolio(&policy, origin)

		// Send mail with the contract to the user
		log.Printf("[fabrickPayment] Policy %s send mail", policyUid)
		var contractbyte []byte
		name := policy.Uid + ".pdf"
		contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
			policy.Contractor.Uid+"/contract_"+name)

		lib.CheckError(err)

		mail.SendMailContract(policy, &[]mail.Attachment{{
			Byte:        base64.StdEncoding.EncodeToString(contractbyte),
			ContentType: "application/pdf",
			Name: policy.Contractor.Name + "_" + policy.Contractor.Surname + "_" +
				strings.ReplaceAll(policy.NameDesc, " ", "_") + "_contratto.pdf",
		}})

		return nil
	}

	log.Printf("[fabrickPayment] ERROR Policy %s with status %s and isPay %t cannot be paid", policyUid, policy.Status, policy.IsPay)
	return errors.New("cannot pay policy")
}

func updateAgencyPortfolio(policy *models.Policy, origin string) {
	if policy.AgencyUid == "" {
		return
	}

	var agency models.Agency
	fireAgency := lib.GetDatasetByEnv(origin, models.AgencyCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgency, policy.AgentUid)
	lib.CheckError(err)
	docsnap.DataTo(&agency)
	agency.Policies = append(agency.Policies, policy.Uid)

	if !lib.SliceContains(agency.Portfolio, policy.Contractor.Uid) {
		agency.Portfolio = append(agency.Portfolio, policy.Contractor.Uid)
	}

	err = lib.SetFirestoreErr(fireAgency, agency.Uid, agency)
	lib.CheckError(err)
}

func updateAgentPortfolio(policy *models.Policy, origin string) {
	if policy.AgentUid == "" {
		return
	}

	var agent models.Agent
	fireAgent := lib.GetDatasetByEnv(origin, models.AgentCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgent, policy.AgentUid)
	lib.CheckError(err)
	docsnap.DataTo(&agent)
	agent.Policies = append(agent.Policies, policy.Uid)

	if !lib.SliceContains(agent.Portfolio, policy.Contractor.Uid) {
		agent.Portfolio = append(agent.Portfolio, policy.Contractor.Uid)
	}

	err = lib.SetFirestoreErr(fireAgent, agent.Uid, agent)
	lib.CheckError(err)
}
