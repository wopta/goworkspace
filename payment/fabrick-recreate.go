package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

type FabrickRefreshPayByLinkRequest struct {
	PolicyUid string `json:"policyUid"`
}

func FabrickRecreateFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[FabrickRecreateFx] Handler start ---------------------------")

	var (
		request     FabrickRefreshPayByLinkRequest
		err         error
		policy      *models.Policy
		flowName    string
		networkNode *models.NetworkNode
		warrant     *models.Warrant
	)

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[FabrickRecreateFx] request body: %s", string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[FabrickRecreateFx] error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	policy, err = FabrickRecreate(request.PolicyUid, origin, networkNode, warrant)
	if err != nil {
		log.Printf("[FabrickRecreateFx] error recreating payment: %s", err.Error())
		return "", nil, err
	}

	flowName = models.ECommerceFlow
	if policy.Channel == models.MgaChannel {
		flowName = models.MgaFlow
	} else {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode != nil {
			warrant = networkNode.GetWarrant()
			if warrant != nil {
				flowName = warrant.GetFlowName(policy.Name)
			}
		}
	}
	log.Printf("[FabrickRecreateFx] flowName '%s'", flowName)

	log.Println("[FabrickRecreateFx] send pay mail to contractor...")
	mail.SendMailPay(
		*policy,
		mail.AddressAnna,
		mail.GetContractorEmail(policy),
		mail.Address{},
		flowName,
	)

	models.CreateAuditLog(r, string(body))

	return `{"success":true}`, `{"success":true}`, nil
}

func FabrickRecreate(policyUid, origin string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Policy, error) {
	log.Println("[FabrickRecreate]")
	var (
		err    error
		policy models.Policy
	)

	policy = plc.GetPolicyByUid(policyUid, origin)
	if policy.IsPay {
		log.Printf("[FabrickRecreate] policy %s already paid cannot recreate payment(s)", policy.Uid)
		return nil, fmt.Errorf("policy %s already paid cannot recreate payment(s)", policy.Uid)
	}

	oldTransactions := tr.GetPolicyTransactions(origin, policy.Uid)

	log.Println("[FabrickRecreate] recreating payment...")
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	payUrl, err := PaymentController(origin, &policy, product)
	if err != nil {
		log.Printf("[FabrickRecreate] error creating payment: %s", err.Error())
		return nil, err
	}

	now := time.Now().UTC()
	policy.PayUrl = payUrl
	policy.Updated = now

	// TODO: automatically delete the transacations on fabrick DB (expireBill)
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	log.Println("[FabrickRecreate] deleting transaction(s)...")
	for _, transaction := range oldTransactions {
		log.Printf("[FabrickRecreate] deleting transaction %s", transaction.Uid)
		transaction.IsDelete = true
		transaction.ExpirationDate = now.AddDate(0, 0, 1).Format(models.TimeDateOnly)
		transaction.Status = models.PolicyStatusDeleted
		transaction.StatusHistory = append(transaction.StatusHistory, transaction.Status)
		transaction.PaymentNote = "Cancellata per ricreazione titoli"

		log.Println("[FabrickRecreate] saving transaction to firestore...")
		err = lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
		if err != nil {
			log.Printf("[FabrickRecreate] error saving transaction to firestore: %s", err.Error())
			return nil, err
		}
		log.Println("[FabrickRecreate] saving transaction to bigquery...")
		transaction.BigQuerySave(origin)
	}

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Println("[FabrickRecreate] saving policy to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
	if err != nil {
		log.Printf("[FabrickRecreate] error saving policy to firestore: %s", err.Error())
		return nil, err
	}

	log.Println("[FabrickRecreate] saving policy to bigquery...")
	policy.BigquerySave(origin)

	return &policy, nil
}
