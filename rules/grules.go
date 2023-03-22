package rules

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/shopspring/decimal"
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

func (p *Fx) CalculateOfferPrices(out *models.RuleOut) {
	/*for _, coveragesPerPaymentFrequency := range out.Coverages {
		for _, coverage := range coveragesPerPaymentFrequency {
			for offerKey, offerValue := range coverage.Offer {
				out.OfferPrice[offerKey][yearly].Net = out.OfferPrice[offerKey][yearly].Net.Add(offerValue.PremiumNet)
				out.OfferPrice[offerKey][yearly].Tax = out.OfferPrice[offerKey][yearly].Tax.Add(offerValue.PremiumTaxAmount)
				out.OfferPrice[offerKey][yearly].Gross = out.OfferPrice[offerKey][yearly].Gross.Add(offerValue.PremiumGross)
				out.OfferPrice[offerKey][monthly].Net = out.OfferPrice[offerKey][monthly].Net.Add(offerValue.PremiumNet.DivRound(decimal.NewFromInt(12), 2))
				out.OfferPrice[offerKey][monthly].Tax = out.OfferPrice[offerKey][monthly].Tax.Add(offerValue.PremiumTaxAmount.DivRound(decimal.NewFromInt(12), 2))
				out.OfferPrice[offerKey][monthly].Gross = out.OfferPrice[offerKey][monthly].Gross.Add(offerValue.PremiumGross.DivRound(decimal.NewFromInt(12), 2))
			}
		}
	}*/

	for _, coverage := range out.Coverages {
		for offerKey, offerValue := range coverage.Offer {
			out.OfferPrice[offerKey][yearly].Net = out.OfferPrice[offerKey][yearly].Net.Add(offerValue.PremiumNetYearly)
			out.OfferPrice[offerKey][yearly].Tax = out.OfferPrice[offerKey][yearly].Tax.Add(offerValue.PremiumTaxAmountYearly)
			out.OfferPrice[offerKey][yearly].Gross = out.OfferPrice[offerKey][yearly].Gross.Add(offerValue.PremiumGrossYearly)
			out.OfferPrice[offerKey][monthly].Net = out.OfferPrice[offerKey][monthly].Net.Add(offerValue.PremiumNetYearly.DivRound(decimal.NewFromInt(12), 2))
			out.OfferPrice[offerKey][monthly].Tax = out.OfferPrice[offerKey][monthly].Tax.Add(offerValue.PremiumTaxAmountYearly.DivRound(decimal.NewFromInt(12), 2))
			out.OfferPrice[offerKey][monthly].Gross = out.OfferPrice[offerKey][monthly].Gross.Add(offerValue.PremiumGrossYearly.DivRound(decimal.NewFromInt(12), 2))
		}
	}
}

func (p *Fx) RoundYearlyOfferPrices(out *models.RuleOut) {
	for offerType, priceStruct := range out.OfferPrice {
		ceilGrossPrice := priceStruct[yearly].Gross.Ceil()
		priceStruct[yearly].Delta = ceilGrossPrice.Sub(priceStruct[yearly].Gross)
		priceStruct[yearly].Gross = ceilGrossPrice
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly.GreaterThan(decimal.NewFromInt(0))
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly = out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly.Add(priceStruct[yearly].Delta)
			newNetPrice := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly.
				DivRound(decimal.NewFromInt(1).
					Add(out.Coverages["IPI"].Tax.Div(decimal.NewFromInt(100))), 2)
			newTax := out.Coverages["IPI"].Offer[offerType].PremiumGrossYearly.Sub(newNetPrice)
			out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.
				Add(newNetPrice.Sub(out.Coverages["IPI"].Offer[offerType].PremiumNetYearly))
			out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.
				Add(newTax.Sub(out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountYearly))
			out.Coverages["IPI"].Offer[offerType].PremiumNetYearly = newNetPrice
			out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountYearly = newTax
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly = out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly.Add(priceStruct[yearly].Delta)
			newNetPrice := out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly.
				DivRound(decimal.NewFromInt(1).
					Add(out.Coverages["DRG"].Tax.Div(decimal.NewFromInt(100))), 2)
			newTax := out.Coverages["DRG"].Offer[offerType].PremiumGrossYearly.Sub(newNetPrice)
			out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.
				Add(newNetPrice.Sub(out.Coverages["DRG"].Offer[offerType].PremiumNetYearly))
			out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.
				Add(newTax.Sub(out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountYearly))
			out.Coverages["DRG"].Offer[offerType].PremiumNetYearly = newNetPrice
			out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountYearly = newTax
		}
		/*
			ceilGrossPrice := priceStruct[yearly].Gross.Ceil()
			priceStruct[yearly].Delta = ceilGrossPrice.Sub(priceStruct[yearly].Gross)
			priceStruct[yearly].Gross = ceilGrossPrice
			hasIPIGuarantee := out.Coverages[yearly]["IPI"].Offer[offerType].PremiumGross.GreaterThan(decimal.NewFromInt(0))
			if hasIPIGuarantee {
				out.Coverages[yearly]["IPI"].Offer[offerType].PremiumGross = out.Coverages[yearly]["IPI"].Offer[offerType].PremiumGross.Add(priceStruct[yearly].Delta)
				newNetPrice := out.Coverages[yearly]["IPI"].Offer[offerType].PremiumGross.
					DivRound(decimal.NewFromInt(1).
						Add(out.Coverages[yearly]["IPI"].Tax.Div(decimal.NewFromInt(100))), 2)
				newTax := out.Coverages[yearly]["IPI"].Offer[offerType].PremiumGross.Sub(newNetPrice)
				out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.
					Add(newNetPrice.Sub(out.Coverages[yearly]["IPI"].Offer[offerType].PremiumNet))
				out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.
					Add(newTax.Sub(out.Coverages[yearly]["IPI"].Offer[offerType].PremiumTaxAmount))
				out.Coverages[yearly]["IPI"].Offer[offerType].PremiumNet = newNetPrice
				out.Coverages[yearly]["IPI"].Offer[offerType].PremiumTaxAmount = newTax
			} else {
				out.Coverages[yearly]["DRG"].Offer[offerType].PremiumGross = out.Coverages[yearly]["DRG"].Offer[offerType].PremiumGross.Add(priceStruct[yearly].Delta)
				newNetPrice := out.Coverages[yearly]["DRG"].Offer[offerType].PremiumGross.
					DivRound(decimal.NewFromInt(1).
						Add(out.Coverages[yearly]["DRG"].Tax.Div(decimal.NewFromInt(100))), 2)
				newTax := out.Coverages[yearly]["DRG"].Offer[offerType].PremiumGross.Sub(newNetPrice)
				out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.
					Add(newNetPrice.Sub(out.Coverages[yearly]["DRG"].Offer[offerType].PremiumNet))
				out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.
					Add(newTax.Sub(out.Coverages[yearly]["DRG"].Offer[offerType].PremiumTaxAmount))
				out.Coverages[yearly]["DRG"].Offer[offerType].PremiumNet = newNetPrice
				out.Coverages[yearly]["DRG"].Offer[offerType].PremiumTaxAmount = newTax
			}
		*/
	}
}

func (p *Fx) RoundMonthlyOfferPrices(out *models.RuleOut) {
	for offerType, priceStruct := range out.OfferPrice {
		nonRoundedGrossPrice := priceStruct[yearly].Gross
		roundedMonthlyGrossPrice := priceStruct[monthly].Gross.Round(0)
		yearlyGrossPrice := roundedMonthlyGrossPrice.Mul(decimal.NewFromInt(12))
		priceStruct[monthly].Delta = (yearlyGrossPrice.Sub(nonRoundedGrossPrice)).DivRound(decimal.NewFromInt(12), 2)
		priceStruct[monthly].Gross = roundedMonthlyGrossPrice
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly.GreaterThan(decimal.NewFromInt(0))
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly = out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly.Add(priceStruct[monthly].Delta)
			newNetPrice := out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly.
				DivRound(decimal.NewFromInt(1).
					Add(out.Coverages["IPI"].Tax.Div(decimal.NewFromInt(100))), 2)
			newTax := out.Coverages["IPI"].Offer[offerType].PremiumGrossMonthly.Sub(newNetPrice)
			out.OfferPrice[offerType][monthly].Net = out.OfferPrice[offerType][monthly].Net.
				Add(newNetPrice.Sub(out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly))
			out.OfferPrice[offerType][monthly].Tax = out.OfferPrice[offerType][monthly].Tax.
				Add(newTax.Sub(out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountMonthly))
			out.Coverages["IPI"].Offer[offerType].PremiumNetMonthly = newNetPrice
			out.Coverages["IPI"].Offer[offerType].PremiumTaxAmountMonthly = newTax
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly = out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly.Add(priceStruct[monthly].Delta)
			newNetPrice := out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly.
				DivRound(decimal.NewFromInt(1).
					Add(out.Coverages["DRG"].Tax.Div(decimal.NewFromInt(100))), 2)
			newTax := out.Coverages["DRG"].Offer[offerType].PremiumGrossMonthly.Sub(newNetPrice)
			out.OfferPrice[offerType][monthly].Net = out.OfferPrice[offerType][monthly].Net.
				Add(newNetPrice.Sub(out.Coverages["DRG"].Offer[offerType].PremiumNetMonthly))
			out.OfferPrice[offerType][monthly].Tax = out.OfferPrice[offerType][monthly].Tax.
				Add(newTax.Sub(out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountMonthly))
			out.Coverages["DRG"].Offer[offerType].PremiumNetMonthly = newNetPrice
			out.Coverages["DRG"].Offer[offerType].PremiumTaxAmountMonthly = newTax
		}
	}
}

/*func (p *Fx) RoundOfferPrices(out *models.RuleOut) {
	for offerType, priceStruct := range out.OfferPrice {
		oldYearlyPriceGross := priceStruct[yearly].Gross
		ceilPriceGrossYear := priceStruct[yearly].Gross.Ceil()
		priceStruct[yearly].Delta = ceilPriceGrossYear.Sub(priceStruct[yearly].Gross)
		priceStruct[yearly].Gross = ceilPriceGrossYear
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumGross.GreaterThan(decimal.NewFromInt(0))
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGross = out.Coverages["IPI"].Offer[offerType].PremiumGross.Add(priceStruct[yearly].Delta)
			oldPremiumNet := out.Coverages["IPI"].Offer[offerType].PremiumNet
			oldPremiumTaxAmount := out.Coverages["IPI"].Offer[offerType].PremiumTaxAmount
			out.Coverages["IPI"].Offer[offerType].PremiumNet = out.Coverages["IPI"].Offer[offerType].PremiumGross.DivRound(decimal.NewFromInt(1).Add(out.Coverages["IPI"].Tax.Div(decimal.NewFromInt(100))), 2)
			out.Coverages["IPI"].Offer[offerType].PremiumTaxAmount = out.Coverages["IPI"].Offer[offerType].PremiumGross.Sub(out.Coverages["IPI"].Offer[offerType].PremiumNet)
			out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.Add(out.Coverages["IPI"].Offer[offerType].PremiumNet.Sub(oldPremiumNet))
			out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.Add(out.Coverages["IPI"].Offer[offerType].PremiumTaxAmount.Sub(oldPremiumTaxAmount))
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGross = out.Coverages["DRG"].Offer[offerType].PremiumGross.Add(priceStruct[yearly].Delta)
			oldPremiumNet := out.Coverages["DRG"].Offer[offerType].PremiumNet
			oldPremiumTaxAmount := out.Coverages["DRG"].Offer[offerType].PremiumTaxAmount
			out.Coverages["DRG"].Offer[offerType].PremiumNet = out.Coverages["DRG"].Offer[offerType].PremiumGross.DivRound(decimal.NewFromInt(1).Add(out.Coverages["DRG"].Tax.Div(decimal.NewFromInt(100))), 2)
			out.Coverages["DRG"].Offer[offerType].PremiumTaxAmount = out.Coverages["DRG"].Offer[offerType].PremiumGross.Sub(out.Coverages["DRG"].Offer[offerType].PremiumNet)
			out.OfferPrice[offerType][yearly].Net = out.OfferPrice[offerType][yearly].Net.Add(out.Coverages["DRG"].Offer[offerType].PremiumNet.Sub(oldPremiumNet))
			out.OfferPrice[offerType][yearly].Tax = out.OfferPrice[offerType][yearly].Tax.Add(out.Coverages["DRG"].Offer[offerType].PremiumTaxAmount.Sub(oldPremiumTaxAmount))
		}

		roundPriceGrossMonth := oldYearlyPriceGross.Div(decimal.NewFromInt(12)).Round(0)                               //priceStruct[monthly].Gross.Round(0)
		priceStruct[monthly].Delta = oldYearlyPriceGross.DivRound(decimal.NewFromInt(12), 2).Sub(roundPriceGrossMonth) //roundPriceGrossMonth.Sub(priceStruct[monthly].Gross)
		priceStruct[monthly].Gross = roundPriceGrossMonth
		if hasIPIGuarantee {
			oldMonthlyNet := out.OfferPrice[offerType][monthly].Net
			out.OfferPrice[offerType][monthly].Net := out.OfferPrice[offerType][monthly].Net.Add(out)
		} else {

		}
	}
}*/

//TODO: implement FilterOfferByMinimumPrice function for both monthly and yearly

func (p *Fx) FilterOffersByMinimumPrice(out *models.RuleOut, yearlyPriceMinimum float64, monthlyPriceMinimum float64) {
	toBeDeleted := make([]string, 0)
	for offerType, priceStruct := range out.OfferPrice {
		hasNotOfferMinimumYearlyPrice := priceStruct[yearly].Gross.LessThan(decimal.NewFromFloat(yearlyPriceMinimum))
		//hasNotOfferMinimumMonthlyPrice := priceStruct[monthly].Gross.LessThan(decimal.NewFromFloat(monthlyPriceMinimum))
		if hasNotOfferMinimumYearlyPrice { // || hasNotOfferMinimumMonthlyPrice {
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

func (p *Fx) NewDecimalFromInt(in int64) decimal.Decimal {
	return decimal.NewFromInt(in)
}

func (p *Fx) NewDecimalFromFloat(in float64) decimal.Decimal {
	return decimal.NewFromFloat(in)
}

func (p *Fx) DecimalToString(d decimal.Decimal) string {
	return d.String()
}

func (p *Fx) AddDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Add(d2)
}

func (p *Fx) SubDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Sub(d2)
}

func (p *Fx) MulDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Mul(d2)
}

func (p *Fx) DivDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Div(d2)
}

func (p *Fx) DivRoundDecimal(d1 decimal.Decimal, d2 decimal.Decimal, precision int64) decimal.Decimal {
	return d1.DivRound(d2, int32(precision))
}

func (p *Fx) AddDecimalWithInt(d1 decimal.Decimal, in int64) decimal.Decimal {
	return d1.Add(decimal.NewFromInt(in))
}

func (p *Fx) SubDecimalWithInt(d1 decimal.Decimal, in int64) decimal.Decimal {
	return d1.Sub(decimal.NewFromInt(in))
}

func (p *Fx) MulDecimalWithInt(d1 decimal.Decimal, in int64) decimal.Decimal {
	return d1.Mul(decimal.NewFromInt(in))
}

func (p *Fx) DivDecimalWithInt(d1 decimal.Decimal, in int64) decimal.Decimal {
	return d1.Div(decimal.NewFromInt(in))
}

func (p *Fx) DivRoundDecimalWithInt(d1 decimal.Decimal, in int64, precision int64) decimal.Decimal {
	return d1.DivRound(decimal.NewFromInt(in), int32(precision))
}

func (p *Fx) AddDecimalWithFloat(d1 decimal.Decimal, in float64) decimal.Decimal {
	return d1.Add(decimal.NewFromFloat(in))
}

func (p *Fx) SubDecimalWithFloat(d1 decimal.Decimal, in float64) decimal.Decimal {
	return d1.Sub(decimal.NewFromFloat(in))
}

func (p *Fx) MulDecimalWithFloat(d1 decimal.Decimal, in float64) decimal.Decimal {
	return d1.Mul(decimal.NewFromFloat(in))
}

func (p *Fx) DivDecimalWithFloat(d1 decimal.Decimal, in float64) decimal.Decimal {
	return d1.Div(decimal.NewFromFloat(in))
}

func (p *Fx) DivRoundDecimalWithFloat(d1 decimal.Decimal, in float64, precision int64) decimal.Decimal {
	return d1.DivRound(decimal.NewFromFloat(in), int32(precision))
}

func (p *Fx) AddDecimalWithString(d1 decimal.Decimal, in string) decimal.Decimal {
	d2, _ := decimal.NewFromString(in)
	return d1.Add(d2)
}

func (p *Fx) SubDecimalWithString(d1 decimal.Decimal, in string) decimal.Decimal {
	d2, _ := decimal.NewFromString(in)
	return d1.Sub(d2)
}

func (p *Fx) MulDecimalWithString(d1 decimal.Decimal, in string) decimal.Decimal {
	d2, _ := decimal.NewFromString(in)
	return d1.Mul(d2)
}

func (p *Fx) DivDecimalWithString(d1 decimal.Decimal, in string) decimal.Decimal {
	d2, _ := decimal.NewFromString(in)
	return d1.Div(d2)
}

func (p *Fx) DivRoundDecimalWithString(d1 decimal.Decimal, in string, precision int64) decimal.Decimal {
	d2, _ := decimal.NewFromString(in)
	return d1.DivRound(d2, int32(precision))
}

func (p *Fx) RoundDecimal(d decimal.Decimal, places int64) decimal.Decimal {
	return d.Round(int32(places))
}
