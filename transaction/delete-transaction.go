package transaction

import (
	"time"

	"github.com/wopta/goworkspace/models"
)

func DeleteTransaction(transaction *models.Transaction, note string) {
	now := time.Now().UTC()

	transaction.IsDelete = true
	transaction.Status = models.TransactionStatusDeleted
	transaction.StatusHistory = append(transaction.StatusHistory, transaction.Status)
	transaction.UpdateDate = now
	transaction.ExpirationDate = now.AddDate(0, 0, -1).Format(models.TimeDateOnly)
	transaction.PaymentNote = note
}
