package fabrick

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/models"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func payTransaction(policy models.Policy, providerId, trSchedule, paymentMethod, collection string, networkNode *models.NetworkNode) (models.Transaction, error) {
	var (
		transaction models.Transaction
		err         error
	)

	if transaction, err = tr.GetTransactionToBePaid(policy.Uid, providerId, trSchedule, collection); err != nil {
		return models.Transaction{}, err
	}
	transaction.IsDelete = false
	transaction.IsPay = true
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
	transaction.PayDate = time.Now().UTC()
	transaction.TransactionDate = transaction.PayDate
	transaction.UpdateDate = transaction.PayDate
	transaction.PaymentMethod = paymentMethod
	transaction.PaymentNote = ""

	return transaction, nil
}
