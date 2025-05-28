package policy

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func UpdatePolicy(policy *models.Policy) (map[string]any, error) {
	input := make(map[string]any, 0)

	input["assets"] = policy.Assets
	input["contractor"] = policy.Contractor
	input["fundsOrigin"] = policy.FundsOrigin
	if policy.Surveys != nil {
		input["surveys"] = policy.Surveys
	}
	if policy.Statements != nil {
		input["statements"] = policy.Statements
	}
	input["step"] = policy.Step
	if policy.OfferlName != "" {
		input["offerName"] = policy.OfferlName
	}
	input["consultancyValue"] = map[string]any{
		"percentage": policy.ConsultancyValue.Percentage,
		"price":      lib.RoundFloat(policy.PriceGross*policy.ConsultancyValue.Percentage, 2),
	}

	switch policy.Name {
	case models.PersonaProduct:
		input["taxAmount"] = policy.TaxAmount
		input["priceNett"] = policy.PriceNett
		input["priceGross"] = policy.PriceGross
		input["taxAmountMonthly"] = policy.TaxAmountMonthly
		input["priceNettMonthly"] = policy.PriceNettMonthly
		input["priceGrossMonthly"] = policy.PriceGrossMonthly
	case models.CommercialCombinedProduct:
		input["startDate"] = policy.StartDate
		input["endDate"] = policy.EndDate
		input["declaredClaims"] = policy.DeclaredClaims
		input["hasBond"] = policy.HasBond
		input["bond"] = policy.Bond
		input["clause"] = policy.Clause
		input["contractors"] = policy.Contractors
		input["priceGroup"] = policy.PriceGroup
	case models.CatNatProduct:
		input["startDate"] = policy.StartDate
		input["endDate"] = policy.EndDate
		input["quoteQuestions"] = policy.QuoteQuestions
		input["offersPrices"] = policy.OffersPrices
		input["contractors"] = policy.Contractors
		input["taxAmount"] = policy.TaxAmount
		input["priceNett"] = policy.PriceNett
		input["priceGross"] = policy.PriceGross
		input["paymentSplit"] = policy.PaymentSplit
	}

	input["updated"] = time.Now().UTC()

	return input, nil
}
