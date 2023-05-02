package rules

import (
	"encoding/json"
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"math"
	"strconv"
	"strings"
)

func rulesFromJson(groule []byte, out interface{}, in []byte, data []byte) (string, interface{}) {

	log.Println("RulesFromJson")
	//rules := lib.CheckEbyte(ioutil.ReadFile("pmi-allrisk.json"))

	fx := &Fx{}
	fxSurvey := &FxSurvey{}
	var err error
	// create new instance of DataContext
	dataContext := ast.NewDataContext()
	// add your JSON Fact into data context using AddJSON() function.
	if in != nil {
		err = dataContext.AddJSON("in", in)
		log.Println("RulesFromJson in")
		lib.CheckError(err)
	}

	if out != nil {
		err = dataContext.Add("out", out)
		log.Println("RulesFromJson out")
		lib.CheckError(err)
	}

	if data != nil {
		err = dataContext.AddJSON("data", data)
		log.Println("RulesFromJson data loaded")
		lib.CheckError(err)
	}

	err = dataContext.Add("fx", fx)
	log.Println("RulesFromJson fx loaded")
	lib.CheckError(err)

	err = dataContext.Add("fxSurvey", fxSurvey)
	log.Println("RulesFromJson fxSurvey loaded")
	lib.CheckError(err)

	underlying := pkg.NewBytesResource(groule)
	lib.CheckError(err)

	resource := pkg.NewJSONResourceFromResource(underlying)
	lib.CheckError(err)
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//bs := pkg.NewBytesResource([]byte(fileRes))

	err = ruleBuilder.BuildRuleFromResource("rules", "0.0.1", resource)
	lib.CheckError(err)
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("rules", "0.0.1")
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataContext, knowledgeBase)
	lib.CheckError(err)

	//resp := "execute"
	b, err := json.Marshal(out)
	lib.CheckError(err)

	return string(b), out
}

type Fx struct {
}

func (p *Fx) ToString(value float64) string {
	var r int
	r = int(math.Round(value))
	log.Println(r)

	return fmt.Sprint(r)
}
func (p *Fx) SetCoverage(value float64) string {

	return fmt.Sprintf("%f", value)
}

func (p *Fx) Tax(value float64) string {

	return fmt.Sprintf("%f", value)
}
func (p *Fx) GetContentValue(buildingType string) float64 {
	buildingType = strings.ToUpper(buildingType)
	if buildingType == "SERVIZI MANUALI" {
		return 0.10
	}
	if buildingType == "COMMERCIALE" {
		return 0.20
	}
	if buildingType == "PRODUZIONE" {
		return 0.30
	}
	if buildingType == "EDILI" {
		return 0.15
	}
	if buildingType == "SERVIZI INTELLETTUALI" {
		return 0.15
	}
	return 0
}
func (p *Fx) GetMachineryvalue(buildingType string) float64 {
	buildingType = strings.ToUpper(buildingType)
	if buildingType == "SERVIZI MANUALI" {
		return 0.10
	}
	if buildingType == "COMMERCIALE" {
		return 0.20
	}
	if buildingType == "EDILI" {
		return 0.15
	}
	if buildingType == "SERVIZI INTELLETTUALI" {
		return 0.10
	}
	if buildingType == "PRODUZIONE" {
		return 0.15
	}
	return 0
}
func (p *Fx) GetTheftValue(buildingType string) float64 {
	buildingType = strings.ToUpper(buildingType)
	if buildingType == "SERVIZI MANUALI" {
		return 0.10
	}
	if buildingType == "COMMERCIALE" {
		return 0.15
	}
	if buildingType == "EDILI" {
		return 0.15
	}
	if buildingType == "SERVIZI INTELLETTUALI" {
		return 0.10
	}
	if buildingType == "PRODUZIONE" {
		return 0.10
	}
	return 0
}
func (p *Fx) GetElectronicValue(buildingType string) float64 {
	buildingType = strings.ToUpper(buildingType)
	if buildingType == "SERVIZI MANUALI" {
		return 0.10
	}
	if buildingType == "COMMERCIALE" {
		return 0.20
	}
	if buildingType == "EDILI" {
		return 0.15
	}
	if buildingType == "SERVIZI INTELLETTUALI" {
		return 0.10
	}
	if buildingType == "PRODUZIONE" {
		return 0.15
	}
	return 0
}
func (p *Fx) GetBuildigValue(buildingType string) int {
	buildingType = strings.ToUpper(buildingType)
	if buildingType == "INDUSTRIALE" {
		return 600
	}
	if buildingType == "COMMERCIALE" {
		return 1000
	}
	if buildingType == "CIVILE_TIPO_UFFICIO" {
		return 1400
	}
	return 0
}
func (p *Fx) Log(any interface{}) {
	log.Println(any)

}

func (p *Fx) FormatString(stringToFormat string, params ...interface{}) string {
	return fmt.Sprintf(stringToFormat, params...)
}

func (p *Fx) AppendString(aString, subString string) string {
	return fmt.Sprintf("%s%s", aString, subString)
}

func (p *Fx) Replace(input string, old string, new string) string {
	return strings.Replace(input, old, new, 1)
}

func (p *Fx) ReplaceAt(input string, replacement string, index int64) string {
	return input[:index] + string(replacement) + input[index+1:]
}

func (p *Fx) RoundNear(value float64, nearest int64) float64 {
	log.Println((math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest))

	return (math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest)
}

func (p *Fx) FloatToString(in float64, decimal int64) string {
	return fmt.Sprintf("%.*f", decimal, in)
}

func (p *Fx) DeleteOfferFromGuarantee(m map[string]*models.GuaranteValue, key string) {
	delete(m, key)
}

func (p *Fx) CalculateMonthlyCoveragePrices(guarantees map[string]*models.Guarante) {
	for _, guarantee := range guarantees {
		for _, offer := range guarantee.Offer {
			offer.PremiumGrossMonthly = offer.PremiumGrossYearly / 12
			offer.PremiumTaxAmountMonthly = offer.PremiumTaxAmountYearly / 12
			offer.PremiumNetMonthly = offer.PremiumNetYearly / 12
		}
	}
}

func (p *Fx) CalculateOfferPrices(guarantees map[string]*models.Guarante, offersPrices map[string]map[string]*models.Price) {
	for _, guarantee := range guarantees {
		for offerKey, offerValue := range guarantee.Offer {
			offersPrices[offerKey][yearly].Net += offerValue.PremiumNetYearly
			offersPrices[offerKey][yearly].Tax += offerValue.PremiumTaxAmountYearly
			offersPrices[offerKey][yearly].Gross += offerValue.PremiumGrossYearly
			offersPrices[offerKey][monthly].Net += offerValue.PremiumNetYearly / 12
			offersPrices[offerKey][monthly].Tax += offerValue.PremiumTaxAmountYearly / 12
			offersPrices[offerKey][monthly].Gross += offerValue.PremiumGrossYearly / 12
		}
	}
}

func (p *Fx) RoundMonthlyOfferPrices(out *RuleOut, roundingCoverages ...string) {
	updatePrices := func(coverage string, offerType string, priceStruct map[string]*models.Price) {
		out.Guarantees[coverage].Offer[offerType].PremiumGrossMonthly += priceStruct[monthly].Delta
		newNetPrice := out.Guarantees[coverage].Offer[offerType].PremiumGrossMonthly / (1 + (out.Guarantees[coverage].Tax / 100))
		newTax := out.Guarantees[coverage].Offer[offerType].PremiumGrossMonthly - newNetPrice
		out.OfferPrice[offerType][monthly].Net += newNetPrice - out.Guarantees[coverage].Offer[offerType].PremiumNetMonthly
		out.OfferPrice[offerType][monthly].Tax += newTax - out.Guarantees[coverage].Offer[offerType].PremiumTaxAmountMonthly
		out.Guarantees[coverage].Offer[offerType].PremiumNetMonthly = newNetPrice
		out.Guarantees[coverage].Offer[offerType].PremiumTaxAmountMonthly = newTax
	}

	for offerType, priceStruct := range out.OfferPrice {
		nonRoundedGrossPrice := priceStruct[yearly].Gross
		roundedMonthlyGrossPrice := math.Round(priceStruct[monthly].Gross)
		yearlyGrossPrice := roundedMonthlyGrossPrice * 12
		priceStruct[monthly].Delta = (yearlyGrossPrice - nonRoundedGrossPrice) / 12
		priceStruct[monthly].Gross = roundedMonthlyGrossPrice

		for _, roundingCoverage := range roundingCoverages {
			hasGuarantee := out.Guarantees[roundingCoverage].Offer[offerType].PremiumNetMonthly > 0
			if hasGuarantee {
				updatePrices(roundingCoverage, offerType, priceStruct)
				break
			}
		}
	}
}

func (p *Fx) RoundYearlyOfferPrices(out *RuleOut, roundingCoverages ...string) {
	updatePrices := func(coverage string, offerType string, priceStruct map[string]*models.Price) {
		out.Guarantees[coverage].Offer[offerType].PremiumGrossYearly += priceStruct[yearly].Delta
		newNetPrice := out.Guarantees[coverage].Offer[offerType].PremiumGrossYearly / (1 + (out.Guarantees[coverage].Tax / 100))
		newTax := out.Guarantees[coverage].Offer[offerType].PremiumGrossYearly - newNetPrice
		out.OfferPrice[offerType][yearly].Net += newNetPrice - out.Guarantees[coverage].Offer[offerType].PremiumNetYearly
		out.OfferPrice[offerType][yearly].Tax += newTax - out.Guarantees[coverage].Offer[offerType].PremiumTaxAmountYearly
		out.Guarantees[coverage].Offer[offerType].PremiumNetYearly = newNetPrice
		out.Guarantees[coverage].Offer[offerType].PremiumTaxAmountYearly = newTax
	}

	for offerType, priceStruct := range out.OfferPrice {
		ceilGrossPrice := math.Ceil(priceStruct[yearly].Gross)
		priceStruct[yearly].Delta = ceilGrossPrice - priceStruct[yearly].Gross
		priceStruct[yearly].Gross = ceilGrossPrice
		for _, roundingCoverage := range roundingCoverages {
			hasGuarantee := out.Guarantees[roundingCoverage].Offer[offerType].PremiumNetMonthly > 0
			if hasGuarantee {
				updatePrices(roundingCoverage, offerType, priceStruct)
				break
			}
		}
	}
}

func (p *Fx) RoundToTwoDecimalPlaces(guarantees map[string]*models.Guarante, offersPrices map[string]map[string]*models.Price) {
	roundFloatTwoDecimals := func(in float64) float64 {
		res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", in), 64)
		return res
	}

	for _, guarantee := range guarantees {
		for _, offerType := range guarantee.Offer {
			offerType.PremiumNetMonthly = roundFloatTwoDecimals(offerType.PremiumNetMonthly)
			offerType.PremiumTaxAmountMonthly = roundFloatTwoDecimals(offerType.PremiumTaxAmountMonthly)
			offerType.PremiumGrossMonthly = roundFloatTwoDecimals(offerType.PremiumGrossMonthly)

			offerType.PremiumNetYearly = lib.RoundFloatTwoDecimals(offerType.PremiumNetYearly)
			offerType.PremiumTaxAmountYearly = lib.RoundFloatTwoDecimals(offerType.PremiumTaxAmountYearly)
			offerType.PremiumGrossYearly = lib.RoundFloatTwoDecimals(offerType.PremiumGrossYearly)
		}
	}

	for _, offerType := range offersPrices {
		offerType[monthly].Net = roundFloatTwoDecimals(offerType[monthly].Net)
		offerType[monthly].Tax = roundFloatTwoDecimals(offerType[monthly].Tax)
		offerType[monthly].Delta = roundFloatTwoDecimals(offerType[monthly].Delta)

		offerType[yearly].Net = lib.RoundFloatTwoDecimals(offerType[yearly].Net)
		offerType[yearly].Tax = lib.RoundFloatTwoDecimals(offerType[yearly].Tax)
		offerType[yearly].Delta = lib.RoundFloatTwoDecimals(offerType[yearly].Delta)
	}
}

func (p *Fx) FilterOffersByMinimumPrice(guarantees map[string]*models.Guarante, offersPrices map[string]map[string]*models.Price, yearlyPriceMinimum float64, monthlyPriceMinimum float64) {
	for offerKey, offer := range offersPrices {
		hasNotOfferMinimumYearlyPrice := offer[yearly].Gross < yearlyPriceMinimum
		hasNotOfferMinimumMonthlyPrice := offer[monthly].Gross < monthlyPriceMinimum
		if hasNotOfferMinimumMonthlyPrice && hasNotOfferMinimumYearlyPrice {
			delete(offersPrices, offerKey)
			for _, guarantee := range guarantees {
				delete(guarantee.Offer, offerKey)
			}
			return
		}
		if hasNotOfferMinimumMonthlyPrice {
			delete(offersPrices[offerKey], monthly)
			for _, guarantee := range guarantees {
				if p.HasGuaranteePerOffer(guarantees, offerKey, guarantee.Slug) {
					guarantee.Offer[offerKey].PremiumNetMonthly = 0.0
					guarantee.Offer[offerKey].PremiumTaxAmountMonthly = 0.0
					guarantee.Offer[offerKey].PremiumGrossMonthly = 0.0
				}
			}
		}
		if hasNotOfferMinimumYearlyPrice {
			delete(offersPrices[offerKey], yearly)
			for _, guarantee := range guarantees {
				if p.HasGuaranteePerOffer(guarantees, offerKey, guarantee.Slug) {
					guarantee.Offer[offerKey].PremiumNetYearly = 0.0
					guarantee.Offer[offerKey].PremiumTaxAmountYearly = 0.0
					guarantee.Offer[offerKey].PremiumGrossYearly = 0.0
				}
			}
		}

	}
}

func (p *Fx) HasGuarantee(guarantees map[string]*models.Guarante, guaranteeKey string) bool {
	for _, guarantee := range guarantees {
		if guarantee.Slug == guaranteeKey {
			return true
		}
	}
	return false
}

func (p *Fx) HasGuaranteePerOffer(guarantees map[string]*models.Guarante, offerSlug string, guaranteeKey string) bool {
	for _, guarantee := range guarantees {
		if guarantee.Slug == guaranteeKey {
			for offerKey, _ := range guarantee.Offer {
				if offerKey == offerSlug {
					return true
				}
			}
		}
	}
	return false
}

func (p *Fx) RemoveGuaranteeIfCondition(guarantees map[string]*models.Guarante, guaranteeKey string, condition bool) {
	if condition {
		p.RemoveGuarantee(guarantees, guaranteeKey)
	}
}

func (p *Fx) RemoveOfferFromGuaranteeIfCondition(guaranteeOffer map[string]*models.GuaranteValue, offerKey string, condition bool) {
	if condition {
		p.DeleteOfferFromGuarantee(guaranteeOffer, offerKey)
	}
}

func (p *Fx) RemoveGuarantee(guarantees map[string]*models.Guarante, guaranteeKey string) {
	delete(guarantees, guaranteeKey)
}

func (p *Fx) RemoveOfferFromGuarantee(guaranteeOffer map[string]*models.GuaranteValue, offerKey string) {
	delete(guaranteeOffer, offerKey)
}

func (p *Fx) RemoveGuaranteesPriceZero(guarantees map[string]*models.Guarante) {
	for _, guarantee := range guarantees {
		for offerKey, _ := range guarantee.Offer {
			if guarantee.Offer[offerKey].PremiumGrossYearly == 0.0 {
				delete(guarantee.Offer, offerKey)
			}
		}
	}
}

func (p *Fx) RemoveOfferPrice(offerPrice map[string]map[string]*models.Price, offerKey string) {
	delete(offerPrice, offerKey)
}
