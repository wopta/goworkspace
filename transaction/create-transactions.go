package transaction

import (
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"time"
)

func CreateTransactions(policy models.Policy, mgaProduct models.Product, uidGenerator func() string) []models.Transaction {
	var (
		transactions = make([]models.Transaction, 0)
	)

	numTransactions := 1
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		numTransactions = 12
	}

	i := 0
	for i < numTransactions {
		// create transaction
		tr := createTransaction(policy, uidGenerator)

		// enrich transaction with date info
		tr = setDateInfo(i, tr, policy)

		// enrich transaction with price info
		// TODO: define a way to determine correct price based on product configuration (e.g., byGuarantee, flat or by offer)
		tr = setPriceInfo(tr, policy)

		// enrich transaction with commissions
		// TODO: define a way to determine right amount of commissions based on renewal
		commissionMga := lib.RoundFloat(product.GetCommissionByProduct(&policy, &mgaProduct, false), 2)
		tr.Commissions = commissionMga

		transactions = append(transactions, tr)
		i++
	}

	return transactions
}

func createTransaction(policy models.Policy, uidGenerator func() string) models.Transaction {
	contractorName := lib.TrimSpace(fmt.Sprintf("%s %s", policy.Contractor.Name, policy.Contractor.Surname))
	return models.Transaction{
		Uid:           uidGenerator(),
		PolicyName:    policy.Name,
		Name:          contractorName,
		PolicyUid:     policy.Uid,
		Company:       policy.Company,
		NumberCompany: policy.CodeCompany,
		ProviderName:  policy.Payment,
		Status:        models.TransactionStatusToPay,
		StatusHistory: []string{models.TransactionStatusToPay},
	}
}

func setDateInfo(index int, transaction models.Transaction, policy models.Policy) models.Transaction {
	now := time.Now().UTC()

	startDate := policy.StartDate
	transaction.EffectiveDate = startDate.AddDate(0, index, 0)
	transaction.ScheduleDate = transaction.EffectiveDate.Format(time.DateOnly)
	transaction.ExpirationDate = startDate.AddDate(10, index, 0).Format(time.DateOnly)
	transaction.CreationDate = now
	transaction.UpdateDate = now

	return transaction
}

func setPriceInfo(transaction models.Transaction, policy models.Policy) models.Transaction {
	priceGross := policy.PriceGross
	priceNet := policy.PriceNett

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		priceGross = policy.PriceGrossMonthly
		priceNet = policy.PriceNettMonthly
	}

	transaction.Amount = priceGross
	transaction.AmountNet = priceNet

	return transaction
}
