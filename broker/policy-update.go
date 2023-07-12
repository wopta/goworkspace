package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

func UpdatePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		input     map[string]interface{}
	)
	log.Println("UpdatePolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &policy)
	if err != nil {
		log.Println("UpdatePolicy: unable to unmarshal request body")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	input = make(map[string]interface{}, 0)
	input["assets"] = policy.Assets
	input["contractor"] = policy.Contractor
	input["fundsOrigin"] = policy.FundsOrigin
	if policy.Surveys != nil {
		input["surveys"] = policy.Surveys
	}
	if policy.Statements != nil {
		input["statements"] = policy.Statements
	}
	input["updated"] = time.Now().UTC()

	lib.FireUpdate(firePolicy, policyUID, input)

	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

func PatchPolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err          error
		policyUID    string
		updateValues map[string]interface{}
	)
	log.Println("UpdatePolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &updateValues)
	if err != nil {
		log.Println("PatchPolicy: unable to unmarshal request body")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		log.Println("PatchPolicy: error during policy update in firestore")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

func DeletePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		request   PolicyDeleteReq
	)
	log.Println("DeletePolicy")
	policyUID = r.Header.Get("uid")
	guaranteFire := lib.GetDatasetByEnv(r.Header.Get("origin"), "guarante")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(req, &request)
	if err != nil {
		log.Printf("DeletePolicy: unable to delete policy %s", policyUID)
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	docsnap := lib.GetFirestore(firePolicy, policyUID)
	docsnap.DataTo(&policy)
	if policy.IsDeleted || !policy.CompanyEmit {
		log.Printf("DeletePolicy: can't delete policy %s", policyUID)
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil

	}
	policy.IsDeleted = true
	policy.DeleteCode = request.Code
	policy.DeleteDesc = request.Description
	policy.DeleteDate = request.Date
	policy.RefundType = request.RefundType
	policy.Status = models.PolicyStatusDeleted
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusDeleted)
	lib.SetFirestore(firePolicy, policyUID, policy)
	policy.BigquerySave(r.Header.Get("origin"))
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)
	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

type PolicyDeleteReq struct {
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	RefundType  string    `json:"refundType,omitempty"`
}

func RefundPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[RefundPolicyFx] Handler start ----------------------------------------")
	var (
		err    error
		policy models.Policy
	)

	policyUid := r.Header.Get("policyUid")
	origin := r.Header.Get("origin")

	log.Printf("[RefundPolicyFx] Uid: %s", policyUid)
	policy, err = GetPolicy(policyUid, origin)
	lib.CheckError(err)
	policyJsonLog, err := policy.Marshal()
	lib.CheckError(err)
	log.Printf("[RefundPolicyFx] Policy %s JSON: %s", policyUid, string(policyJsonLog))

	if !policy.IsPay {
		log.Printf("[RefundPolicyFx] ERROR Policy %s is not paid, cannot refund", policyUid)
		return `{"success":false}`, `{"success":false}`, nil
	}

	if !isPolicyInRefundPeriod(&policy) {
		log.Printf("[RefundPolicyFx] ERROR Policy %s is out of refund period with startDate: %s", policyUid, policy.StartDate.String())
		return `{"success":false}`, `{"success":false}`, nil
	}

	err = RefundPolicy(&policy, origin)

	if err != nil {
		log.Printf("[RefundPolicyFx] ERROR Policy %s: %s", policyUid, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	return `{"success":true}`, `{"success":true}`, nil
}

// TODO: this should be set in product
// ex.: product.RefundPeriod = 30
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

func RefundPolicy(policy *models.Policy, origin string) error {
	transactionsFire := lib.GetDatasetByEnv(origin, "transactions")

	policyTransactions := transaction.GetPolicyTransactions(origin, policy.Uid)

	refundedTransactions := []models.Transaction{}
	for _, refundedTransaction := range policyTransactions {
		if !refundedTransaction.IsPay {
			break
		}

		transactionUid := lib.NewDoc(transactionsFire)

		refundedTransaction.Uid = transactionUid
		refundedTransaction.Amount = -refundedTransaction.Amount
		refundedTransaction.AmountNet = refundedTransaction.AmountNet * -1
		refundedTransaction.IsPay = false
		refundedTransaction.IsDelete = false
		refundedTransaction.PayDate = time.Time{}
		refundedTransaction.CreationDate = time.Now().UTC()
		refundedTransaction.StartDate = policy.StartDate
		refundedTransaction.Status = models.TransactionStatusToPay
		refundedTransaction.StatusHistory = []string{models.TransactionStatusToPay}
		refundedTransaction.ScheduleDate = ""       //?
		refundedTransaction.ExpirationDate = ""     //?
		refundedTransaction.ProviderId = ""         //?
		refundedTransaction.UserToken = ""          //?
		refundedTransaction.ProviderName = "manual" //?

		refundedTransactions = append(refundedTransactions, refundedTransaction)

		jsonRt, _ := json.Marshal(refundedTransaction)

		log.Printf("Refunded transaction: %v", string(jsonRt))
	}

	if len(refundedTransactions) == 0 {
		return fmt.Errorf("no transaction to be refunded")
	}

	return nil
}
