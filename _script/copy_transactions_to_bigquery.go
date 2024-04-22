package _script

import (
	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
	"time"
)

func CopyTransactionsToBigQuery() {
	//transactions := transaction.GetPolicyTransactions("", "wT6LRDMwSHViSbTCj5GM") DEV
	transactions := transaction.GetPolicyTransactions("", "6tztJcR7KwqBT7JHwnwI")

	for index, tr := range transactions {
		if tr.IsPay {
			t := deepcopy.Copy(tr).(models.Transaction)
			t.Status = models.TransactionStatusToPay
			t.StatusHistory = []string{models.TransactionStatusToPay}
			t.IsPay = false
			t.PaymentMethod = ""
			t.UpdateDate = t.CreationDate
			t.PayDate = time.Time{}
			t.TransactionDate = time.Time{}
			t.BigQuerySave("")
		}
		transactions[index].BigQuerySave("")
	}
}
