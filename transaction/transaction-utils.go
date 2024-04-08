package transaction

import (
	"github.com/wopta/goworkspace/models"
	"time"
)

func ReinitializePaymentInfo(tr *models.Transaction) {
	tr.IsPay = false
	tr.IsDelete = false
	tr.PaymentNote = ""
	tr.PaymentMethod = ""
	tr.PayDate = time.Time{}
	tr.TransactionDate = time.Time{}
	tr.Status = models.TransactionStatusToPay
	tr.StatusHistory = append(tr.StatusHistory, models.TransactionStatusToPay)
}
