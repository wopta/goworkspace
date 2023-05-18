package quote

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
	"io"
	"net/http"
	"strconv"
)

func PersonFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		personaTassi map[string]json.RawMessage
	)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	policy := sellable.Person(body)

	b := lib.GetByteByEnv("quote/persona-tassi.json", false)
	err := json.Unmarshal(b, &personaTassi)
	lib.CheckError(err)

	for _, guarantee := range policy.Assets[0].Guarantees {
		switch guarantee.Slug {
		case "IPI":
			calculateIPIPrices(policy.Contractor, &guarantee, personaTassi)
		case "D":
			calculateDPrices(policy.Contractor, &guarantee, personaTassi)
		case "DRG":
			calculateDRGPrices(policy.Contractor, &guarantee, personaTassi)
		case "ITI":
			calculateITIPrices(policy.Contractor, &guarantee, personaTassi)
		case "DC":
			calculateDCPrices(policy.Contractor, &guarantee, personaTassi)
		case "RSC":
			calculateRSCPrices(policy.Contractor, &guarantee, personaTassi)
		case "IPM":
			contractorAge, err := policy.CalculateContractorAge()
			lib.CheckError(err)
			if contractorAge < 66 {
				calculateIPMPrices(contractorAge, &guarantee, personaTassi)
			}
		}
	}

	policyJson, err := policy.Marshal()

	return string(policyJson), policy, err
}

func calculateIPIPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = lib.RoundFloat((offer.SumInsuredLimitOfIndemnity/1000.0)*tassi[guarantee.Type][contractor.RiskClass][offer.DeductibleType][offer.Deductible], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["D"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = lib.RoundFloat((offer.SumInsuredLimitOfIndemnity/1000.0)*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDRGPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DRG"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateITIPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["ITI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[contractor.RiskClass][guarantee.Offer[offerKey].Deductible], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDCPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DC"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateRSCPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["RSC"], &tassi)
	lib.CheckError(err)

	sumInsuredLimitOfIndemnity := strconv.FormatFloat(guarantee.Offer["premium"].SumInsuredLimitOfIndemnity, 'f', -1, 64)

	guarantee.Offer["premium"].PremiumNetYearly = lib.RoundFloat(guarantee.Offer["premium"].SumInsuredLimitOfIndemnity*tassi[guarantee.Type][contractor.RiskClass][sumInsuredLimitOfIndemnity], 2)
	guarantee.Offer["premium"].PremiumTaxAmountYearly = lib.RoundFloat((guarantee.Tax*guarantee.Offer["premium"].PremiumNetYearly)/100, 2)
	guarantee.Offer["premium"].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer["premium"].PremiumTaxAmountYearly+guarantee.Offer["premium"].PremiumNetYearly, 2)

	guarantee.Offer["premium"].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer["premium"].PremiumNetYearly/12, 2)
	guarantee.Offer["premium"].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer["premium"].PremiumTaxAmountYearly/12, 2)
	guarantee.Offer["premium"].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer["premium"].PremiumGrossYearly/12, 2)

}

func calculateIPMPrices(contractorAge int, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPM"], &tassi)
	lib.CheckError(err)

	age := strconv.Itoa(contractorAge)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly = offer.SumInsuredLimitOfIndemnity * tassi[age]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly = (guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly = guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}
