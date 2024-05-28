package fabrick

import (
	"time"

	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

func payTransaction(policy models.Policy, providerId, trSchedule, paymentMethod string, networkNode *models.NetworkNode) (models.Transaction, error) {
	var (
		transaction models.Transaction
		mgaProduct  *models.Product
		err         error
	)

	if transaction, err = tr.GetTransactionToBePaid(policy.Uid, providerId, trSchedule, ""); err != nil {
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

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	// TODO: this method still saves all entities inside it.
	// We should extract to batch save them only at the end
	if err = tr.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct); err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}
