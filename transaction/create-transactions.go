package transaction

import (
	"fmt"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

var numTransactionsMap = map[string]int{
	string(models.PaySplitMonthly):           12,
	string(models.PaySplitYear):              1,
	string(models.PaySplitYearly):            1,
	string(models.PaySplitSingleInstallment): 1,
	string(models.PaySplitSemestral):         2,
}

func CreateTransactions(policy models.Policy, mgaProduct models.Product, uidGenerator func() string) []models.Transaction {
	var (
		numTransactions int
		transactions    = make([]models.Transaction, 0)
	)

	numTransactions = numTransactionsMap[policy.PaymentSplit]
	if numTransactions == 0 {
		return transactions
	}

	for i := 0; i < numTransactions; i++ {
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
	}

	return transactions
}

func createTransaction(policy models.Policy, uidGenerator func() string) models.Transaction {
	contractorName := lib.TrimSpace(fmt.Sprintf("%s %s", policy.Contractor.Name, policy.Contractor.Surname))
	return models.Transaction{
		Uid:           uidGenerator(),
		PolicyName:    policy.Name,
		Name:          contractorName,
		Annuity:       policy.Annuity,
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

	startDate := lib.AddMonths(policy.StartDate, 12*transaction.Annuity)
	transaction.EffectiveDate = lib.AddMonths(startDate, index)
	transaction.ScheduleDate = transaction.EffectiveDate.Format(time.DateOnly)
	// TODO: code smell - this info is needed for Fabrick only but we don't know the params for other scenarios
	transaction.ExpirationDate = lib.AddMonths(startDate, 18).Format(time.DateOnly)
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
