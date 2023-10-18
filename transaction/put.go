package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func PutByPolicy(
	policy models.Policy,
	scheduleDate, origin, expireDate, customerId string,
	amount, amountNet float64,
	providerId, paymentMethod string,
	isPay bool,
) *models.Transaction {
	var (
		sd              string
		status          string
		statusHistory   = make([]string, 0)
		payDate         time.Time
		transactionDate time.Time
	)

	log.Println("[PutByPolicy] start -----------------------------------------")
	log.Printf("[PutByPolicy] Policy %s", policy.Uid)

	if scheduleDate == "" {
		sd = time.Now().UTC().Format(models.TimeDateOnly)
	} else {
		sd = scheduleDate
	}

	now := time.Now().UTC()

	if isPay {
		status = models.TransactionStatusPay
		statusHistory = append(statusHistory, models.TransactionStatusToPay, models.TransactionStatusPay)
		payDate = now
		transactionDate = now
	} else {
		status = models.TransactionStatusToPay
		statusHistory = append(statusHistory, models.TransactionStatusToPay)
	}

	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	transactionUid := lib.NewDoc(fireTransactions)

	// All commission get calculated apart on payment
	tr := models.Transaction{
		Amount:          amount,
		AmountNet:       amountNet,
		Id:              "",
		Uid:             transactionUid,
		PolicyName:      policy.Name,
		PolicyUid:       policy.Uid,
		CreationDate:    now,
		Status:          status,
		StatusHistory:   statusHistory,
		ScheduleDate:    sd,
		ExpirationDate:  expireDate,
		NumberCompany:   policy.CodeCompany,
		IsPay:           isPay,
		PayDate:         payDate,
		TransactionDate: transactionDate,
		Name:            policy.Contractor.Name + " " + policy.Contractor.Surname,
		Company:         policy.Company,
		IsDelete:        false,
		ProviderId:      providerId,
		UserToken:       customerId,
		ProviderName:    policy.Payment,
		AgentUid:        policy.AgentUid,
		AgencyUid:       policy.AgencyUid,
		PaymentMethod:   paymentMethod,
	}

	log.Println("[PutByPolicy] saving transaction to firestore...")
	err := lib.SetFirestoreErr(fireTransactions, transactionUid, tr)
	if err != nil {
		log.Printf("[PutByPolicy] error saving transaction to firestore: %s", err.Error())
		return nil
	}

	log.Println("[PutByPolicy] saving transaction to bigquery...")
	tr.BigQuerySave(origin)

	log.Println("[PutByPolicy] end -------------------------------------------")

	return &tr
}
