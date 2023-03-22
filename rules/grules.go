package rules

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func rulesFromJson(groule []byte, out interface{}, in []byte, data []byte) (string, interface{}) {

	log.Println("RulesFromJson")
	//rules := lib.CheckEbyte(ioutil.ReadFile("pmi-allrisk.json"))

	fx := &Fx{}
	// create new instance of DataContext
	dataContext := ast.NewDataContext()
	// add your JSON Fact into data context using AddJSON() function.
	err := dataContext.AddJSON("in", in)
	log.Println("RulesFromJson in")
	lib.CheckError(err)
	err = dataContext.Add("out", out)
	//err = dataContext.AddJSON("out", out)
	log.Println("RulesFromJson out")
	lib.CheckError(err)
	err = dataContext.AddJSON("data", data)
	log.Println("RulesFromJson data loaded")
	lib.CheckError(err)
	err = dataContext.Add("fx", fx)
	log.Println("RulesFromJson fx loaded")
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

func (p *Fx) FormatString(stringToFormat string, params ...string) string {
	return fmt.Sprintf(stringToFormat, params)
}

func (p *Fx) AppendString(aString, subString string) string {
	return fmt.Sprintf("%s%s", aString, subString)
}

func (p *Fx) ReplaceAt(input string, replacement string, index int64) string {
	return input[:index] + string(replacement) + input[index+1:]
}

func (p *Fx) RoundNear(value float64, nearest int64) float64 {
	log.Println((math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest))

	return (math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest)
}

func (p *Fx) CalculateMonthlyCoveragePrices(out *models.RuleOut) {
	for _, coverageValue := range out.Coverages {
		for _, offer := range coverageValue.Offer {
			offer.PremiumGrossMonthly = offer.PremiumGrossYearly / 12
			offer.PremiumTaxAmountMonthly = offer.PremiumTaxAmountYearly / 12
			offer.PremiumNetMonthly = offer.PremiumNetYearly / 12
		}
	}
}

func (p *Fx) CalculateOfferPrices(out *models.RuleOut) {
	for _, coverage := range out.Coverages {
		for offerKey, offerValue := range coverage.Offer {
			out.OfferPrice[offerKey][yearly].Net += offerValue.PremiumNetYearly
			out.OfferPrice[offerKey][yearly].Tax += offerValue.PremiumTaxAmountYearly
			out.OfferPrice[offerKey][yearly].Gross += offerValue.PremiumGrossYearly
			out.OfferPrice[offerKey][monthly].Net += offerValue.PremiumNetYearly / 12
			out.OfferPrice[offerKey][monthly].Tax += offerValue.PremiumTaxAmountYearly / 12
			out.OfferPrice[offerKey][monthly].Gross += offerValue.PremiumGrossYearly / 12
		}
	}
}

func (p *Fx) RoundYearlyOfferPrices(out *models.RuleOut) {
	for offerType, priceStruct := range out.OfferPrice {
		ceilGrossPrice := math.Ceil(priceStruct[yearly].Gross)
		priceStruct[yearly].Delta = ceilGrossPrice - priceStruct[yearly].Gross
		priceStruct[yearly].Gross = ceilGrossPrice
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly > 0
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly += priceStruct[yearly].Delta
			newNetPrice := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly / (1 + (out.Coverages["IPI"].Tax / 100))
			newTax := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly - newNetPrice
			out.OfferPrice[offerType][yearly].Net += newNetPrice - out.Coverages["IPI"].Offer[offerType].PremiumNetYearly
			out.OfferPrice[offerType][yearly].Tax += newTax - out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountYearly
			out.Coverages["IPI"].Offer[offerType].PremiumNetYearly = newNetPrice
			out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountYearly = newTax
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly += priceStruct[yearly].Delta
			newNetPrice := out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly / (1 + (out.Coverages["DRG"].Tax / 100))
			newTax := out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly - newNetPrice
			out.OfferPrice[offerType][yearly].Net += newNetPrice - out.Coverages["DRG"].Offer[offerType].PremiumNetYearly
			out.OfferPrice[offerType][yearly].Tax += newTax - out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountYearly
			out.Coverages["DRG"].Offer[offerType].PremiumNetYearly = newNetPrice
			out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountYearly = newTax
		}
	}
}

func (p *Fx) RoundMonthlyOfferPrices(out *models.RuleOut) {
	for offerType, priceStruct := range out.OfferPrice {
		nonRoundedGrossPrice := priceStruct[yearly].Gross
		roundedMonthlyGrossPrice := math.Round(priceStruct[monthly].Gross)
		yearlyGrossPrice := roundedMonthlyGrossPrice * 12
		priceStruct[monthly].Delta = (yearlyGrossPrice - nonRoundedGrossPrice) / 12
		priceStruct[monthly].Gross = roundedMonthlyGrossPrice
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly > 0
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly += priceStruct[monthly].Delta
			newNetPrice := out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly / (1 + (out.Coverages["IPI"].Tax / 100))
			newTax := out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly - newNetPrice
			out.OfferPrice[offerType][monthly].Net += newNetPrice - out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly
			out.OfferPrice[offerType][monthly].Tax += newTax - out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountMonthly
			out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly = newNetPrice
			out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountMonthly = newTax
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly += priceStruct[monthly].Delta
			newNetPrice := out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly / (1 + (out.Coverages["DRG"].Tax / 100))
			newTax := out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly - newNetPrice
			out.OfferPrice[offerType][monthly].Net += newNetPrice - out.Coverages["DRG"].Offer[offerType].PremiumNetMonthly
			out.OfferPrice[offerType][monthly].Tax += newTax - out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountMonthly
			out.Coverages["DRG"].Offer[offerType].PremiumNetMonthly = newNetPrice
			out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountMonthly = newTax
		}
	}
}

func (p *Fx) RoundToTwoDecimalPlaces(out *models.RuleOut) {
	for _, coverage := range out.Coverages {
		for _, offerType := range coverage.Offer {
			offerType.PremiumNetMonthly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumNetMonthly), 64)
			offerType.PremiumTaxAmountMonthly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumTaxAmountMonthly), 64)
			offerType.PremiumGrossMonthly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumGrossMonthly), 64)

			offerType.PremiumNetYearly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumNetYearly), 64)
			offerType.PremiumTaxAmountYearly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumTaxAmountYearly), 64)
			offerType.PremiumGrossYearly, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType.PremiumGrossYearly), 64)
		}
	}

	for _, offerType := range out.OfferPrice {
		offerType[monthly].Net, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[monthly].Net), 64)
		offerType[monthly].Tax, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[monthly].Tax), 64)
		offerType[monthly].Delta, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[monthly].Delta), 64)

		offerType[yearly].Net, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[yearly].Net), 64)
		offerType[yearly].Tax, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[yearly].Tax), 64)
		offerType[yearly].Delta, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", offerType[yearly].Delta), 64)
	}
}

func (p *Fx) FilterOffersByMinimumPrice(out *models.RuleOut, yearlyPriceMinimum float64, monthlyPriceMinimum float64) {
	toBeDeleted := make([]string, 0)
	for offerType, priceStruct := range out.OfferPrice {
		hasNotOfferMinimumYearlyPrice := priceStruct[yearly].Gross < yearlyPriceMinimum
		hasNotOfferMinimumMonthlyPrice := priceStruct[monthly].Gross < monthlyPriceMinimum
		if hasNotOfferMinimumYearlyPrice || hasNotOfferMinimumMonthlyPrice {
			log.Println(offerType)
			toBeDeleted = append(toBeDeleted, offerType)
		}
	}

	for _, offerType := range toBeDeleted {
		delete(out.OfferPrice, offerType)
		for _, coverage := range out.Coverages {
			delete(coverage.Offer, offerType)
		}
	}
}
