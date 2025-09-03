package transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/product"
)

func CreateTransactions(policy models.Policy, mgaProduct models.Product, uidGenerator func() string) []models.Transaction {
	var (
		split           = policy.PaymentComponents.Split
		numTransactions int
		transactions    = make([]models.Transaction, 0)
	)

	// Retrocompatibility with policies without PaymentComponents
	if split == "" {
		split = models.PaySplit(policy.PaymentSplit)
	}

	numTransactions = models.PaySplitRateMap[split]
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
		tr = setPriceInfo(i, tr, policy)

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

	provider := policy.PaymentComponents.Provider
	// Retrocompatibility with policies without PaymentComponents
	if provider == "" {
		provider = policy.Payment
	}

	return models.Transaction{
		Uid:           uidGenerator(),
		PolicyName:    policy.Name,
		Name:          contractorName,
		Annuity:       policy.Annuity,
		PolicyUid:     policy.Uid,
		Company:       policy.Company,
		NumberCompany: policy.CodeCompany,
		ProviderName:  provider,
		Status:        models.TransactionStatusToPay,
		StatusHistory: []string{models.TransactionStatusToPay},
	}
}

func setDateInfo(index int, transaction models.Transaction, policy models.Policy) models.Transaction {
	now := time.Now().UTC()

	startDate := lib.AddMonths(policy.StartDate, 12*transaction.Annuity)
	var split models.PaySplit = models.PaySplit(policy.PaymentSplit)
	if split == "" {
		split = models.PaySplitMonthly
	}
	transaction.EffectiveDate = lib.AddMonths(startDate, index*models.PaySplitMonthsMap[split])
	transaction.ScheduleDate = transaction.EffectiveDate.Format(time.DateOnly)
	// TODO: code smell - this info is needed for Fabrick only but we don't know the params for other scenarios
	transaction.ExpirationDate = lib.AddMonths(now, 18).Format(time.DateOnly)
	transaction.CreationDate = now
	transaction.UpdateDate = now

	return transaction
}

const (
	ItemRate        = "rate"
	ItemConsultancy = "consultancy"
)

func setPriceInfo(index int, transaction models.Transaction, policy models.Policy) models.Transaction {
	priceComponent := policy.PaymentComponents.PriceSplit
	if index == 0 {
		priceComponent = policy.PaymentComponents.PriceFirstSplit
	}

	transaction.Amount = priceComponent.Total
	transaction.AmountNet = lib.RoundFloat(priceComponent.Total-priceComponent.Tax, 2)

	// Retrocompatibility with policies without PaymentComponents
	if transaction.Amount == 0 {
		priceGross := policy.PriceGross
		priceNet := policy.PriceNett

		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			priceGross = policy.PriceGrossMonthly
			priceNet = policy.PriceNettMonthly
		}

		transaction.Amount = priceGross
		transaction.AmountNet = priceNet
	}

	if priceComponent.Consultancy > 0 {
		transaction.Items = make([]models.Item, 0, 2)
		transaction.Items = append(transaction.Items, models.Item{
			Type:          fmt.Sprintf("%s-%s", ItemRate, policy.Name),
			Uid:           uuid.NewString(),
			EffectiveDate: transaction.EffectiveDate,
			AmountGross:   priceComponent.Gross,
			AmountNett:    priceComponent.Nett,
			AmountTax:     priceComponent.Tax,
		})
		transaction.Items = append(transaction.Items, models.Item{
			Type:          ItemConsultancy,
			Uid:           uuid.NewString(),
			EffectiveDate: transaction.EffectiveDate,
			AmountGross:   priceComponent.Consultancy,
			AmountNett:    priceComponent.Consultancy,
			AmountTax:     0,
		})
	}

	return transaction
}
