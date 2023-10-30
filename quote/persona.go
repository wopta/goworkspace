package quote

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/sellable"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func PersonaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  models.Policy
		warrant *models.Warrant
	)

	log.Println("[PersonaFx] handler start ----------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[PersonaFx] body: %s", string(body))

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[PersonaFx] error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("[PersonaFx] error getting authToken from idToken: %s", err.Error())
		return "", nil, err
	}

	log.Println("[PersonaFx] loading network node")
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("[PersonaFx] start quoting")

	err = Persona(&policy, authToken.GetChannelByRoleV2(), networkNode, warrant)

	policyJson, err := policy.Marshal()

	log.Printf("[PersonaFx] response: %s", string(policyJson))

	log.Println("[PersonaFx] handler end ------------------------------------")

	return string(policyJson), policy, err
}

func Persona(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) error {
	var personaRates map[string]json.RawMessage

	log.Println("[Persona] function start -----------------------------------")

	personProduct := sellable.Persona(*policy, channel, networkNode, warrant)

	b := lib.GetFilesByEnv(fmt.Sprintf("%s%s/%s/taxes.json", models.ProductsFolder, policy.Name,
		policy.ProductVersion))
	err := json.Unmarshal(b, &personaRates)
	if err != nil {
		log.Printf("[Persona] error unmarshaling persona rates: %s", err.Error())
		return err
	}

	policy.StartDate = time.Now().UTC()
	policy.EndDate = policy.StartDate.AddDate(1, 0, 0)

	log.Println("[Persona] populating policy guarantees list")

	guaranteesList := make([]models.Guarante, 0)
	for _, guarantee := range personProduct.Companies[0].GuaranteesMap {
		guaranteesList = append(guaranteesList, *guarantee)
	}
	policy.Assets[0].Guarantees = guaranteesList

	log.Println("[Persona] init offer prices struct")

	initOfferPrices(policy, personProduct)

	log.Println("[Persona] calculate guarantees prices")

	for _, guarantee := range policy.Assets[0].Guarantees {
		switch guarantee.Slug {
		case "IPI":
			calculateIPIPrices(policy.Contractor, &guarantee, personaRates)
		case "D":
			calculateDPrices(policy.Contractor, &guarantee, personaRates)
		case "DRG":
			calculateDRGPrices(policy.Contractor, &guarantee, personaRates)
		case "ITI":
			calculateITIPrices(policy.Contractor, &guarantee, personaRates)
		case "DC":
			calculateDCPrices(policy.Contractor, &guarantee, personaRates)
		case "RSC":
			calculateRSCPrices(policy.Contractor, &guarantee, personaRates)
		case "IPM":
			contractorAge, err := policy.CalculateContractorAge()
			if err != nil {
				log.Printf("[Persona] error calculate contractor age: %s", err.Error())
				return err
			}
			if contractorAge < 66 {
				calculateIPMPrices(contractorAge, &guarantee, personaRates)
			}
		}
	}

	log.Println("[Persona] applying discounts")

	applyDiscounts(policy)

	log.Println("[Persona] calculate offers prices")

	calculatePersonaOfferPrices(policy)

	log.Println("[Persona] round offers prices")

	roundMonthlyOfferPrices(policy, "IPI", "DRG")

	roundYearlyOfferPrices(policy, "IPI", "DRG")

	roundOfferPrices(policy.OffersPrices)

	roundToTwoDecimalPlaces(policy)

	log.Println("[Persona] filter by minimum price")

	filterOffersByMinimumPrice(policy, 120.0, 50.0)

	log.Println("[Persona] function end -----------------------------------")

	return nil
}

func initOfferPrices(policy *models.Policy, personProduct *models.Product) {
	policy.OffersPrices = make(map[string]map[string]*models.Price)

	for offerKey, _ := range personProduct.Offers {
		policy.OffersPrices[offerKey] = map[string]*models.Price{
			"monthly": {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
			"yearly": {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
		}
	}
}

func calculateIPIPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			(offer.SumInsuredLimitOfIndemnity / 1000.0) *
				tassi[guarantee.Type][contractor.RiskClass][offer.DeductibleType][offer.Deductible]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly =
			guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func calculateDPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["D"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			(offer.SumInsuredLimitOfIndemnity / 1000.0) * tassi[guarantee.Type][contractor.RiskClass]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly =
			guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func calculateDRGPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DRG"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			offer.SumInsuredLimitOfIndemnity * tassi[guarantee.Type][contractor.RiskClass]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly =
			guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func calculateITIPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["ITI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			offer.SumInsuredLimitOfIndemnity * tassi[contractor.RiskClass][guarantee.Offer[offerKey].Deductible]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly =
			guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func calculateDCPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DC"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			offer.SumInsuredLimitOfIndemnity * tassi[guarantee.Type][contractor.RiskClass]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly =
			guarantee.Offer[offerKey].PremiumTaxAmountYearly + guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func calculateRSCPrices(contractor models.User, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["RSC"], &tassi)
	lib.CheckError(err)

	sumInsuredLimitOfIndemnity :=
		strconv.FormatFloat(guarantee.Offer["premium"].SumInsuredLimitOfIndemnity, 'f', -1, 64)

	guarantee.Offer["premium"].PremiumNetYearly =
		tassi[guarantee.Type][contractor.RiskClass][sumInsuredLimitOfIndemnity]
	guarantee.Offer["premium"].PremiumTaxAmountYearly =
		(guarantee.Tax * guarantee.Offer["premium"].PremiumNetYearly) / 100
	guarantee.Offer["premium"].PremiumGrossYearly =
		guarantee.Offer["premium"].PremiumTaxAmountYearly + guarantee.Offer["premium"].PremiumNetYearly

	guarantee.Offer["premium"].PremiumNetMonthly = guarantee.Offer["premium"].PremiumNetYearly / 12
	guarantee.Offer["premium"].PremiumTaxAmountMonthly = guarantee.Offer["premium"].PremiumTaxAmountYearly / 12
	guarantee.Offer["premium"].PremiumGrossMonthly = guarantee.Offer["premium"].PremiumGrossYearly / 12

}

func calculateIPMPrices(contractorAge int, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPM"], &tassi)
	lib.CheckError(err)

	age := strconv.Itoa(contractorAge)

	for offerKey, _ := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			(guarantee.Offer[offerKey].SumInsuredLimitOfIndemnity / 1000) * tassi[age]
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			(guarantee.Tax * guarantee.Offer[offerKey].PremiumNetYearly) / 100
		guarantee.Offer[offerKey].PremiumGrossYearly = guarantee.Offer[offerKey].PremiumTaxAmountYearly +
			guarantee.Offer[offerKey].PremiumNetYearly

		guarantee.Offer[offerKey].PremiumNetMonthly = guarantee.Offer[offerKey].PremiumNetYearly / 12
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = guarantee.Offer[offerKey].PremiumTaxAmountYearly / 12
		guarantee.Offer[offerKey].PremiumGrossMonthly = guarantee.Offer[offerKey].PremiumGrossYearly / 12
	}

}

func applyDiscounts(policy *models.Policy) {
	numberOfGuarantees := map[string]int{
		"base": 0, "your": 0, "premium": 0,
	}
	numberOfInsured := len(policy.Assets)

	guaranteesDiscount := map[int]float64{
		0: 1.0, 1: 1.0, 2: 1.0, 3: 0.97, 4: 0.95, 5: 0.91, 6: 0.90,
	}

	insuredDiscount := map[int]float64{
		1: 1.0, 2: 0.97, 3: 0.92, 4: 0.88, 5: 0.85, 6: 0.80, 7: 0.80, 8: 0.80, 9: 0.80, 10: 0.80,
	}

	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			for offerKey, offer := range guarantee.Offer {
				if offer.PremiumGrossYearly > 0.0 {
					numberOfGuarantees[offerKey]++
				}
			}
		}
	}

	for assetIndex, _ := range policy.Assets {
		for guaranteeIndex, guarantee := range policy.Assets[assetIndex].Guarantees {
			for offerKey, offer := range guarantee.Offer {
				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossYearly =
					offer.PremiumGrossYearly * insuredDiscount[numberOfInsured] * guaranteesDiscount[numberOfGuarantees[offerKey]]
				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountYearly =
					policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossYearly * (guarantee.Tax / 100)
				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetYearly =
					policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossYearly -
						policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountYearly

				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossMonthly =
					policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossYearly / 12
				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountMonthly =
					policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountYearly / 12
				policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetMonthly =
					policy.Assets[assetIndex].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetYearly / 12
			}

		}
	}

}

func calculatePersonaOfferPrices(policy *models.Policy) {
	for offerKey, _ := range policy.OffersPrices {
		for _, guarantee := range policy.Assets[0].Guarantees {
			if guarantee.Offer[offerKey] != nil {
				policy.OffersPrices[offerKey]["monthly"].Net += guarantee.Offer[offerKey].PremiumNetMonthly
				policy.OffersPrices[offerKey]["monthly"].Tax += guarantee.Offer[offerKey].PremiumTaxAmountMonthly
				policy.OffersPrices[offerKey]["monthly"].Gross += guarantee.Offer[offerKey].PremiumGrossMonthly
				policy.OffersPrices[offerKey]["yearly"].Net += guarantee.Offer[offerKey].PremiumNetYearly
				policy.OffersPrices[offerKey]["yearly"].Tax += guarantee.Offer[offerKey].PremiumTaxAmountYearly
				policy.OffersPrices[offerKey]["yearly"].Gross += guarantee.Offer[offerKey].PremiumGrossYearly
			}
		}
	}
}

func roundMonthlyOfferPrices(policy *models.Policy, roundingGuarantees ...string) {
	guarantees := policy.GuaranteesToMap()

	for offerKey, offer := range policy.OffersPrices {
		nonRoundedGrossPrice := offer["yearly"].Gross
		roundedMonthlyGrossPrice := math.Round(offer["monthly"].Gross)
		yearlyGrossPrice := roundedMonthlyGrossPrice * 12
		offer["monthly"].Delta = (yearlyGrossPrice - nonRoundedGrossPrice) / 12
		offer["monthly"].Gross = roundedMonthlyGrossPrice

		for _, roundingGuarantee := range roundingGuarantees {
			hasGuarantee := guarantees[roundingGuarantee].Offer[offerKey] != nil &&
				guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly > 0
			if hasGuarantee {
				guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly += offer["monthly"].Delta
				newNetPrice := guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly /
					(1 + (guarantees[roundingGuarantee].Tax / 100))
				newTax := guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly - newNetPrice
				offer["monthly"].Net += newNetPrice - guarantees[roundingGuarantee].Offer[offerKey].PremiumNetMonthly
				offer["monthly"].Tax += newTax - guarantees[roundingGuarantee].Offer[offerKey].PremiumTaxAmountMonthly
				guarantees[roundingGuarantee].Offer[offerKey].PremiumNetMonthly = newNetPrice
				guarantees[roundingGuarantee].Offer[offerKey].PremiumTaxAmountMonthly = newTax
				break
			}
		}

	}

	guaranteesList := make([]models.Guarante, 0)

	for _, guarantee := range guarantees {
		guaranteesList = append(guaranteesList, guarantee)
	}

	policy.Assets[0].Guarantees = guaranteesList
}

func roundYearlyOfferPrices(policy *models.Policy, roundingGuarantees ...string) {
	guarantees := policy.GuaranteesToMap()

	for offerKey, offer := range policy.OffersPrices {
		ceilGrossPrice := math.Ceil(offer["yearly"].Gross)
		offer["yearly"].Delta = ceilGrossPrice - offer["yearly"].Gross
		offer["yearly"].Gross = ceilGrossPrice
		for _, roundingCoverage := range roundingGuarantees {
			hasGuarantee := guarantees[roundingCoverage].Offer[offerKey] != nil &&
				guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly > 0
			if hasGuarantee {
				guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly += offer["yearly"].Delta
				newNetPrice := guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly /
					(1 + (guarantees[roundingCoverage].Tax / 100))
				newTax := guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly - newNetPrice
				offer["yearly"].Net += newNetPrice - guarantees[roundingCoverage].Offer[offerKey].PremiumNetYearly
				offer["yearly"].Tax += newTax - guarantees[roundingCoverage].Offer[offerKey].PremiumTaxAmountYearly
				guarantees[roundingCoverage].Offer[offerKey].PremiumNetYearly = newNetPrice
				guarantees[roundingCoverage].Offer[offerKey].PremiumTaxAmountYearly = newTax
				break
			}
		}
	}

	guaranteesList := make([]models.Guarante, 0)

	for _, guarantee := range guarantees {
		guaranteesList = append(guaranteesList, guarantee)
	}

	policy.Assets[0].Guarantees = guaranteesList
}

func roundToTwoDecimalPlaces(policy *models.Policy) {
	for guaranteeIndex, guarantee := range policy.Assets[0].Guarantees {
		for offerKey, _ := range guarantee.Offer {
			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetMonthly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetMonthly, 2)
			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountMonthly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountMonthly, 2)
			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossMonthly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossMonthly, 2)

			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetYearly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly, 2)
			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountYearly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly, 2)
			policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossYearly =
				lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly, 2)
		}
	}

	for offerKey, offerValue := range policy.OffersPrices {
		for paymentKey, _ := range offerValue {
			policy.OffersPrices[offerKey][paymentKey].Net =
				lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Net, 2)
			policy.OffersPrices[offerKey][paymentKey].Tax =
				lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Tax, 2)
			policy.OffersPrices[offerKey][paymentKey].Gross =
				lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Gross, 2)
			policy.OffersPrices[offerKey][paymentKey].Delta =
				lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Delta, 2)
		}
	}
}

func filterOffersByMinimumPrice(policy *models.Policy, yearlyPriceMinimum float64, monthlyPriceMinimum float64) {
	for offerKey, offer := range policy.OffersPrices {
		hasNotOfferMinimumYearlyPrice := offer["yearly"].Gross < yearlyPriceMinimum
		hasNotOfferMinimumMonthlyPrice := offer["monthly"].Gross < monthlyPriceMinimum
		if hasNotOfferMinimumMonthlyPrice && hasNotOfferMinimumYearlyPrice {
			delete(policy.OffersPrices, offerKey)
			for guaranteeIndex, _ := range policy.Assets[0].Guarantees {
				delete(policy.Assets[0].Guarantees[guaranteeIndex].Offer, offerKey)
			}
			continue
		}
		if hasNotOfferMinimumMonthlyPrice {
			delete(policy.OffersPrices[offerKey], "monthly")
			for guaranteeIndex, _ := range policy.Assets[0].Guarantees {
				if policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey] != nil {
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossMonthly = 0.0
				}
			}
		}
		if hasNotOfferMinimumYearlyPrice {
			delete(policy.OffersPrices[offerKey], "yearly")
			for guaranteeIndex, _ := range policy.Assets[0].Guarantees {
				if policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey] != nil {
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossMonthly = 0.0
				}
			}
		}

	}
}
