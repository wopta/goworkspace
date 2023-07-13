package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
	"github.com/wopta/goworkspace/transaction"
)

type RefundPolicyRequest struct {
	RefundCode string `json:"refundCode"`
}

func RefundPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[RefundPolicyFx] Handler start ----------------------------------------")

	var (
		err     error
		policy  models.Policy
		request RefundPolicyRequest
	)

	origin := r.Header.Get("origin")
	policyUid := r.Header.Get("policyUid")
	log.Printf("[RefundPolicyFx] Policy: %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[RefundPolicyFx] Request body: %s", string(body))
	err = json.Unmarshal([]byte(body), &request)
	lib.CheckError(err)

	policy, err = GetPolicy(policyUid, origin)
	lib.CheckError(err)

	err = RefundPolicy(&policy, origin, request.RefundCode)

	if err != nil {
		log.Printf("[RefundPolicyFx] ERROR Policy %s: %s", policyUid, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	return `{"success":true}`, `{"success":true}`, nil
}

// TODO: this should be set in product
// ex.: Product.RefundPeriod = 30
func isPolicyInRefundPeriod(policy *models.Policy) bool {
	now := time.Now().UTC()
	elapsedDaysFromStartDate := now.Sub(policy.StartDate).Hours() / 24

	switch policy.Name {
	case "life":
		return elapsedDaysFromStartDate <= 30
	default:
		return elapsedDaysFromStartDate <= 14
	}
}

func RefundPolicy(policy *models.Policy, origin, refundCode string) error {
	log.Println("[RefundPolicy]")

	if !policy.IsPay {
		return errors.New("cannot refund unpaid policy")
	}

	if !isPolicyInRefundPeriod(policy) {
		return errors.New("policy out of refund period")
	}

	// TODO: check refundCode to evaluate amount to be refunded or if refund is disabled

	operations := make(map[string]map[string]interface{})
	refundedTransactions := make(map[string]interface{})
	transactionsFire := lib.GetDatasetByEnv(origin, "transactions")
	policyTransactions := transaction.GetPolicyTransactions(origin, policy.Uid)

	if !policyTransactions[0].IsPay {
		return fmt.Errorf("no transaction to be refunded")
	}

	for _, refundedTransaction := range policyTransactions {
		if !refundedTransaction.IsPay && refundedTransaction.ProviderName == "fabrick" {
			log.Printf("[RefundPolicy] Expire fabrick transaction %s", refundedTransaction.ProviderId)
			err := payment.FabrickExpireBill(&refundedTransaction)
			lib.CheckError(err)

			continue
		}

		transactionUid := lib.NewDoc(transactionsFire)

		refundedTransactions[transactionUid] = models.Transaction{
			Uid:           transactionUid,
			Amount:        -refundedTransaction.Amount,
			AmountNet:     -refundedTransaction.AmountNet,
			Commissions:   -refundedTransaction.Commissions,
			CreationDate:  time.Now().UTC(),
			StartDate:     policy.StartDate,
			Status:        models.TransactionStatusToPay,
			StatusHistory: []string{models.TransactionStatusToPay},
			IsPay:         false,
			IsDelete:      false,
			ProviderName:  "manual",
		}
	}

	log.Printf("[RefundPolicy] %d transactions to be refunded", len(refundedTransactions))

	policyFire := lib.GetDatasetByEnv(origin, "policy")
	policy.Status = models.PolicyStatusToRefund
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)

	operations[transactionsFire] = refundedTransactions
	policyMap := make(map[string]interface{})
	policyMap[policy.Uid] = policy
	operations[policyFire] = policyMap

	jsonOp, err := json.Marshal(operations)
	lib.CheckError(err)
	log.Printf("[RefundPolicy] firestore set operations: %s", string(jsonOp))

	err = lib.SetBatchFirestoreErr(operations)

	return err
}
