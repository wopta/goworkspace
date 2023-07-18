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

// TODO review payload
type DeleteRefundPolicyPayload struct {
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	RefundType  string    `json:"refundType,omitempty"`
}

func DeleteRefundPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[DeleteRefundPolicyFx] Handler start ------------------------")

	var (
		err     error
		policy  models.Policy
		request DeleteRefundPolicyPayload
	)

	policyUid := r.Header.Get("uid")
	origin := r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(body, &request)

	errorResponse := fmt.Sprintf(`{"uid":"%s","success":false}`, policyUid)
	successResponse := fmt.Sprintf(`{"uid":"%s","success":true}`, policyUid)

	if err != nil {
		log.Printf("[DeleteRefundPolicyFx] ERROR Unmarshal body %s", policyUid)
		return errorResponse, errorResponse, nil
	}

	policy, err = GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("[DeleteRefundPolicyFx] ERROR getting policy %s", policyUid)
		return errorResponse, errorResponse, nil
	}

	operations := make(map[string]map[string]interface{})
	modifiedTransactions := make(map[string]interface{})
	modifiedPolicy := make(map[string]interface{})
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionCollection)
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	policyTransactions := transaction.GetPolicyTransactions(origin, policyUid)

	if idx := lib.SliceFindIndex(policyTransactions, func(t models.Transaction) bool { return t.IsPay }); idx < 0 {
		log.Printf("[DeleteRefundPolicyFx] no paid transaction for policy %s", policyUid)
		err = deletePolicy(origin, &policy, request)
		if err != nil {
			log.Printf("[DeleteRefundPolicyFx] ERROR %s", err.Error())
			return errorResponse, errorResponse, nil
		}
		err = <-deleteTransactions(&policyTransactions)
		if err != nil {
			log.Printf("[DeleteRefundPolicyFx] ERROR %s", err.Error())
			return errorResponse, errorResponse, nil
		}
		for _, t := range policyTransactions {
			modifiedTransactions[t.Uid] = t
		}
	} else {
		log.Printf("[DeleteRefundPolicyFx] found paid transaction %s for policy %s", policyTransactions[idx].Uid, policyUid)
		err = refundPolicy(origin, &policy, request)
		if err != nil {
			log.Printf("[DeleteRefundPolicyFx] ERROR %s", err.Error())
			return errorResponse, errorResponse, nil
		}
		refundedTransactions := refundTransactions(&policyTransactions, policy.StartDate, request)
		err = <-deleteTransactions(&policyTransactions)
		if err != nil {
			log.Printf("[DeleteRefundPolicyFx] ERROR %s", err.Error())
			return errorResponse, errorResponse, nil
		} else {
			for _, t := range policyTransactions {
				if t.IsPay {
					modifiedTransactions[t.Uid] = t
				}
			}
		}

		for _, rt := range refundedTransactions {
			modifiedTransactions[rt.(models.Transaction).Uid] = rt.(models.Transaction)
		}
	}

	policy.Updated = time.Now().UTC()
	modifiedPolicy[policy.Uid] = &policy
	operations[firePolicy] = modifiedPolicy
	operations[fireTransactions] = modifiedTransactions

	jsonOp, err := json.Marshal(operations)
	lib.CheckError(err)
	log.Printf("[DeleteRefundPolicyFx] firestore set operations: %s", string(jsonOp))

	err = lib.SetBatchFirestoreErr(operations)
	if err != nil {
		log.Printf("[DeleteRefundPolicyFx] ERROR %s", err.Error())
		return errorResponse, errorResponse, nil
	}

	return successResponse, successResponse, nil
}

func deletePolicy(origin string, policy *models.Policy, payload DeleteRefundPolicyPayload) error {
	log.Printf("[deletePolicy] %s", policy.Uid)

	if policy.IsDeleted || !policy.CompanyEmit {
		return fmt.Errorf("cannot delete policy %s", policy.Uid)
	}

	policy.IsDeleted = true
	policy.Status = models.PolicyStatusDeleted
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusDeleted)
	// TODO review how the payload data will be sent and saved
	policy.DeleteCode = payload.Code        // REVIEW
	policy.DeleteDesc = payload.Description // REVIEW
	policy.DeleteDate = payload.Date        // REVIEW
	policy.RefundType = payload.RefundType  // REVIEW

	return nil
}

func deleteTransactions(transactions *[]models.Transaction) <-chan error {
	result := make(chan error)

	go func() {
		defer close(result)
		for _, t := range *transactions {
			if t.IsPay {
				continue
			}
			if t.ProviderName == "fabrick" {
				err := payment.FabrickExpireBill(&t)
				if err != nil {
					result <- errors.New("error deleting transaction")
					break
				}
			}
		}
		result <- nil
	}()

	return result
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

func refundPolicy(origin string, policy *models.Policy, payload DeleteRefundPolicyPayload) error {
	log.Printf("[refundPolicy] %s", policy.Uid)

	if !policy.IsPay {
		return errors.New("cannot refund unpaid policy")
	}

	if !isPolicyInRefundPeriod(policy) {
		return errors.New("policy out of refund period")
	}

	policy.Status = models.PolicyStatusToRefund
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)

	return nil
}

func refundTransactions(transactions *[]models.Transaction, startDate time.Time, payload DeleteRefundPolicyPayload) map[string]interface{} {
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionCollection)
	refundedTransactions := make(map[string]interface{})

	for _, refundedTransaction := range *transactions {
		if !refundedTransaction.IsPay {
			continue
		}

		transactionUid := lib.NewDoc(fireTransactions)

		// TODO: check refundCode to evaluate amount to be refunded or if refund is disabled
		refundedTransactions[transactionUid] = models.Transaction{
			Uid:           transactionUid,
			Amount:        -refundedTransaction.Amount,
			AmountNet:     -refundedTransaction.AmountNet,
			Commissions:   -refundedTransaction.Commissions,
			CreationDate:  time.Now().UTC(),
			StartDate:     startDate,
			Status:        models.TransactionStatusToPay,
			StatusHistory: []string{models.TransactionStatusToPay},
			IsPay:         false,
			IsDelete:      false,
			ProviderName:  "manual",
		}
	}

	return refundedTransactions
}
