package policy

import (
	"time"

	"github.com/wopta/goworkspace/models"
)

func UpdatePolicy(policy *models.Policy) map[string]interface{} {
	input := make(map[string]interface{}, 0)

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
	}

	input["updated"] = time.Now().UTC()

	return input
}
