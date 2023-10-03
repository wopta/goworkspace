package transaction

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CreateNetworkTransaction(
	policy *models.Policy,
	transaction *models.Transaction,
	commission float64, // Percentual
	accountType string,
	paymentType string,
) *models.NetworkTransaction {
	log.Printf(
		"[CreateNetworkTransaction] accountType '%s' paymentType '%s' commission '%f' amount '%f'",
		accountType,
		paymentType,
		commission,
		transaction.Amount,
	)

	var amount float64

	switch paymentType {
	case models.PaymentTypeRemittanceCompany, models.PaymentTypeCommission:
		amount = transaction.Amount * commission
	case models.PaymentTypeRemittanceMga:
		amount = transaction.Amount - (transaction.Amount * commission)
	}

	if accountType == models.AccountTypePassive {
		amount = -amount
	}

	netTransaction := models.NetworkTransaction{
		Uid:              uuid.New().String(),
		PolicyUid:        policy.Uid,
		TransactionUid:   transaction.Uid,
		NetworkUid:       "", // extract from policy.NetworkUid
		NetworkNodeUid:   "", // extract from policy.ProducerUid
		NetworkNodeType:  "", // extract from policy.ProducerType
		AccountType:      accountType,
		PaymentType:      paymentType,
		Amount:           amount,
		AmountNet:        amount, // TBD
		Name:             "",
		Status:           models.NetworkTransactionStatusCreated,
		StatusHistory:    models.NetworkTransactionStatusCreated,
		IsPay:            false,
		IsConfirmed:      false,
		CreationDate:     lib.GetBigQueryNullDateTime(transaction.CreationDate),
		PayDate:          lib.GetBigQueryNullDateTime(time.Time{}),
		TransactionDate:  lib.GetBigQueryNullDateTime(time.Time{}),
		ConfirmationDate: lib.GetBigQueryNullDateTime(time.Time{}),
	}

	jsonLog, _ := json.Marshal(&netTransaction)

	err := netTransaction.SaveBigQuery()
	if err != nil {
		log.Printf("[CreateNetworkTransaction] error saving network transaction to bigquery: %s", err.Error())
		return nil
	}

	log.Printf("[CreateNetworkTransaction] network transaction created! %s", string(jsonLog))

	return &netTransaction
}
