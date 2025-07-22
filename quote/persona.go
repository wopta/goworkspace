package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

func PersonaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  models.Policy
		warrant *models.Warrant
	)

	log.AddPrefix("PersonaFx")
	defer log.PopPrefix()
	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.ErrorF("error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("error getting authToken from idToken: %s", err.Error())
		return "", nil, err
	}

	flow := authToken.GetChannelByRoleV2()

	log.Println("loading network node")
	nodeUid := policy.PartnershipName
	if strings.EqualFold(policy.Channel, models.NetworkChannel) {
		nodeUid = authToken.UserID
	}
	networkNode := network.GetNetworkNodeByUid(nodeUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		if warrant != nil {
			flow = warrant.GetFlowName(policy.Name)
		}
	}

	log.Println("start quoting")

	if err = Persona(&policy, authToken.GetChannelByRoleV2(), networkNode, warrant, flow); err != nil {
		log.ErrorF("error on quote: %s", err.Error())
		return "", nil, err
	}

	policyJson, err := policy.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(policyJson), policy, err
}

func Persona(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, flow string) error {
	var personaRates map[string]json.RawMessage
	log.AddPrefix("Persona")
	defer log.PopPrefix()
	log.Println("function start -----------------------------------")

	personProduct := sellable.Persona(*policy, channel, networkNode, warrant)

	availableRate := internal.GetAvailableRates(personProduct, flow)

	b := lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/taxes.json", models.ProductsFolder,
		policy.Name, policy.ProductVersion))
	err := json.Unmarshal(b, &personaRates)
	if err != nil {
		log.ErrorF("error unmarshaling persona rates: %s", err.Error())
		return err
	}

	policy.StartDate = lib.SetDateToStartOfDay(time.Now().UTC())
	policy.EndDate = lib.AddMonths(policy.StartDate, 12)

	log.Println("populating policy guarantees list")

	guaranteesList := make([]models.Guarante, 0)
	for _, guarantee := range personProduct.Companies[0].GuaranteesMap {
		guaranteesList = append(guaranteesList, *guarantee)
	}
	policy.Assets[0].Guarantees = guaranteesList

	log.Println("init offer prices struct")

	initOfferPrices(policy, personProduct)

	log.Println("calculate guarantees prices")

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
				log.ErrorF("error calculate contractor age: %s", err.Error())
				return err
			}
			if contractorAge < 66 {
				calculateIPMPrices(contractorAge, &guarantee, personaRates)
			}
		}
	}

	log.Println("applying discounts")

	applyDiscounts(policy)

	log.Println("calculate offers prices")

	calculatePersonaOfferPrices(policy)

	log.Println("round offers prices")

	roundMonthlyOfferPrices(policy, "IPI", "DRG")

	roundYearlyOfferPrices(policy, "IPI", "DRG")

	roundToTwoDecimalPlaces(policy)

	log.Println("apply consultacy price")

	internal.AddConsultacyPrice(policy, personProduct)

	log.Println("filter by minimum price")

	companyIdx := slices.IndexFunc(personProduct.Companies, func(c models.Company) bool {
		return c.Name == policy.Company
	})

	filterOffersByMinimumPrice(policy, personProduct.Companies[companyIdx].MinimumYearlyPrice, personProduct.Companies[companyIdx].MinimumMonthlyPrice)

	log.Println("filtering available rates")

	internal.RemoveOfferRate(policy, availableRate)

	log.Println("function end -----------------------------------")

	return nil
}

func initOfferPrices(policy *models.Policy, personProduct *models.Product) {
	policy.OffersPrices = make(map[string]map[string]*models.Price)

	for offerKey := range personProduct.Offers {
		policy.OffersPrices[offerKey] = map[string]*models.Price{
			string(models.PaySplitMonthly): {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
			string(models.PaySplitYearly): {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
		}
	}
}

func calculateIPIPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat((offer.SumInsuredLimitOfIndemnity/1000.0)*
				tassi[guarantee.Type][contractor.RiskClass][offer.DeductibleType][offer.Deductible], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["D"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat((offer.SumInsuredLimitOfIndemnity/1000.0)*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDRGPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DRG"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateITIPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["ITI"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[contractor.RiskClass][guarantee.Offer[offerKey].Deductible], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateDCPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["DC"], &tassi)
	lib.CheckError(err)

	for offerKey, offer := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat(offer.SumInsuredLimitOfIndemnity*tassi[guarantee.Type][contractor.RiskClass], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func calculateRSCPrices(contractor models.Contractor, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]map[string]map[string]float64
	)

	err := json.Unmarshal(personaTassi["RSC"], &tassi)
	lib.CheckError(err)

	for offerKey := range guarantee.Offer {
		sumInsuredLimitOfIndemnity :=
			strconv.FormatFloat(guarantee.Offer[offerKey].SumInsuredLimitOfIndemnity, 'f', -1, 64)

		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat(tassi[guarantee.Type][contractor.RiskClass][sumInsuredLimitOfIndemnity], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly =
			lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}
}

func calculateIPMPrices(contractorAge int, guarantee *models.Guarante, personaTassi map[string]json.RawMessage) {
	var (
		tassi map[string]float64
	)

	err := json.Unmarshal(personaTassi["IPM"], &tassi)
	lib.CheckError(err)

	age := strconv.Itoa(contractorAge)

	for offerKey := range guarantee.Offer {
		guarantee.Offer[offerKey].PremiumNetYearly =
			lib.RoundFloat((guarantee.Offer[offerKey].SumInsuredLimitOfIndemnity/1000)*tassi[age], 2)
		guarantee.Offer[offerKey].PremiumTaxAmountYearly =
			lib.RoundFloat((guarantee.Tax*guarantee.Offer[offerKey].PremiumNetYearly)/100, 2)
		guarantee.Offer[offerKey].PremiumGrossYearly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly+
			guarantee.Offer[offerKey].PremiumNetYearly, 2)

		guarantee.Offer[offerKey].PremiumNetMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumNetYearly/12, 2)
		guarantee.Offer[offerKey].PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumTaxAmountYearly/12, 2)
		guarantee.Offer[offerKey].PremiumGrossMonthly = lib.RoundFloat(guarantee.Offer[offerKey].PremiumGrossYearly/12, 2)
	}

}

func applyDiscounts(policy *models.Policy) {
	numberOfGuarantees := make(map[string]int)
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

	for assetIndex := range policy.Assets {
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
	for offerKey := range policy.OffersPrices {
		for _, guarantee := range policy.Assets[0].Guarantees {
			if guarantee.Offer[offerKey] != nil {
				policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Net = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Net+guarantee.Offer[offerKey].PremiumNetMonthly, 2)
				policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Tax = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Tax+guarantee.Offer[offerKey].PremiumTaxAmountMonthly, 2)
				policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Gross = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitMonthly)].Gross+guarantee.Offer[offerKey].PremiumGrossMonthly, 2)
				policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Net = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Net+guarantee.Offer[offerKey].PremiumNetYearly, 2)
				policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Tax = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Tax+guarantee.Offer[offerKey].PremiumTaxAmountYearly, 2)
				policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Gross = lib.RoundFloat(policy.OffersPrices[offerKey][string(models.PaySplitYearly)].Gross+guarantee.Offer[offerKey].PremiumGrossYearly, 2)
			}
		}
	}

}

func roundMonthlyOfferPrices(policy *models.Policy, roundingGuarantees ...string) {
	guarantees := policy.GuaranteesToMap()

	for offerKey, offer := range policy.OffersPrices {
		nonRoundedGrossPrice := offer[string(models.PaySplitYearly)].Gross
		roundedMonthlyGrossPrice := math.Round(offer[string(models.PaySplitMonthly)].Gross)
		yearlyGrossPrice := roundedMonthlyGrossPrice * 12
		offer[string(models.PaySplitMonthly)].Delta = (yearlyGrossPrice - nonRoundedGrossPrice) / 12
		offer[string(models.PaySplitMonthly)].Gross = roundedMonthlyGrossPrice

		for _, roundingGuarantee := range roundingGuarantees {
			hasGuarantee := guarantees[roundingGuarantee].Offer[offerKey] != nil &&
				guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly > 0
			if hasGuarantee {
				guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly += offer[string(models.PaySplitMonthly)].Delta
				newNetPrice := guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly /
					(1 + (guarantees[roundingGuarantee].Tax / 100))
				newTax := guarantees[roundingGuarantee].Offer[offerKey].PremiumGrossMonthly - newNetPrice
				offer[string(models.PaySplitMonthly)].Net += newNetPrice - guarantees[roundingGuarantee].Offer[offerKey].PremiumNetMonthly
				offer[string(models.PaySplitMonthly)].Tax += newTax - guarantees[roundingGuarantee].Offer[offerKey].PremiumTaxAmountMonthly
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
		// TODO: sthe original production approved used .Ceil but the test cases use .Round
		ceilGrossPrice := math.Round(offer[string(models.PaySplitYearly)].Gross)
		offer[string(models.PaySplitYearly)].Delta = ceilGrossPrice - offer[string(models.PaySplitYearly)].Gross
		offer[string(models.PaySplitYearly)].Gross = ceilGrossPrice
		for _, roundingCoverage := range roundingGuarantees {
			hasGuarantee := guarantees[roundingCoverage].Offer[offerKey] != nil &&
				guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly > 0
			if hasGuarantee {
				guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly += offer[string(models.PaySplitYearly)].Delta
				newNetPrice := guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly /
					(1 + (guarantees[roundingCoverage].Tax / 100))
				newTax := guarantees[roundingCoverage].Offer[offerKey].PremiumGrossYearly - newNetPrice
				offer[string(models.PaySplitYearly)].Net += newNetPrice - guarantees[roundingCoverage].Offer[offerKey].PremiumNetYearly
				offer[string(models.PaySplitYearly)].Tax += newTax - guarantees[roundingCoverage].Offer[offerKey].PremiumTaxAmountYearly
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
		for offerKey := range guarantee.Offer {
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
		hasNotOfferMinimumYearlyPrice := offer[string(models.PaySplitYearly)].Gross < yearlyPriceMinimum
		hasNotOfferMinimumMonthlyPrice := offer[string(models.PaySplitMonthly)].Gross < monthlyPriceMinimum
		if hasNotOfferMinimumMonthlyPrice && hasNotOfferMinimumYearlyPrice {
			delete(policy.OffersPrices, offerKey)
			for guaranteeIndex, _ := range policy.Assets[0].Guarantees {
				delete(policy.Assets[0].Guarantees[guaranteeIndex].Offer, offerKey)
			}
			continue
		}
		if hasNotOfferMinimumMonthlyPrice {
			delete(policy.OffersPrices[offerKey], string(models.PaySplitMonthly))
			for guaranteeIndex, _ := range policy.Assets[0].Guarantees {
				if policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey] != nil {
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumNetMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumTaxAmountMonthly = 0.0
					policy.Assets[0].Guarantees[guaranteeIndex].Offer[offerKey].PremiumGrossMonthly = 0.0
				}
			}
		}
		if hasNotOfferMinimumYearlyPrice {
			delete(policy.OffersPrices[offerKey], string(models.PaySplitYearly))
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
