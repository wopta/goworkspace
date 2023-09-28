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
	commission float64,
	accountType string,
	paymentType string,
) *models.NetworkTransaction {
	amount := transaction.Amount - commission
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

	log.Printf("[CreateNetworkTransaction] network transaction created! %s", string(jsonLog))

	return &netTransaction
}
