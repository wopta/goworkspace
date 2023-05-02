package models

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	monthly = "monthly"
	yearly  = "yearly"
)

type Fx struct{}

func (fx *Fx) ToString(value float64) string {
	var r int
	r = int(math.Round(value))
	log.Println(r)

	return fmt.Sprint(r)
}
func (fx *Fx) SetCoverage(value float64) string {

	return fmt.Sprintf("%f", value)
}

func (fx *Fx) Tax(value float64) string {

	return fmt.Sprintf("%f", value)
}
func (fx *Fx) GetContentValue(buildingType string) float64 {
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
func (fx *Fx) GetMachineryvalue(buildingType string) float64 {
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
func (fx *Fx) GetTheftValue(buildingType string) float64 {
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
func (fx *Fx) GetElectronicValue(buildingType string) float64 {
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
func (fx *Fx) GetBuildigValue(buildingType string) int {
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
func (fx *Fx) Log(any interface{}) {
	log.Println(any)

}

func (fx *Fx) FormatString(stringToFormat string, params ...interface{}) string {
	return fmt.Sprintf(stringToFormat, params...)
}

func (fx *Fx) AppendString(aString, subString string) string {
	return fmt.Sprintf("%s%s", aString, subString)
}

func (fx *Fx) Replace(input string, old string, new string) string {
	return strings.Replace(input, old, new, 1)
}

func (fx *Fx) ReplaceAt(input string, replacement string, index int64) string {
	return input[:index] + string(replacement) + input[index+1:]
}

func (fx *Fx) RoundNear(value float64, nearest int64) float64 {
	log.Println((math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest))

	return (math.Round(float64(value)/float64(nearest)) * float64(nearest)) - float64(nearest)
}

func (fx *Fx) FloatToString(in float64, decimal int64) string {
	return fmt.Sprintf("%.*f", decimal, in)
}

/*
	SURVEY AND STATEMENT CUSTOM FUNCTIONS
*/

func (fx *Fx) AppendStatement(statements []*Statement, title string, subtitle string, hasMultipleAnswers bool, hasAnswer bool, expectedAnswer bool) []*Statement {
	statement := &Statement{
		Title:              title,
		Subtitle:           subtitle,
		HasMultipleAnswers: nil,
		Questions:          make([]*Question, 0),
		Answer:             nil,
		HasAnswer:          hasAnswer,
		ExpectedAnswer:     nil,
	}
	if hasAnswer {
		statement.ExpectedAnswer = &expectedAnswer
	}
	if hasMultipleAnswers {
		statement.HasMultipleAnswers = &hasMultipleAnswers
	}
	return append(statements, statement)
}

func (fx *Fx) AppendSurvey(surveys []*Survey, title string, subtitle string, hasMultipleAnswers bool, hasAnswer bool, expectedAnswer bool) []*Survey {
	survey := &Survey{
		Title:              title,
		Subtitle:           subtitle,
		HasMultipleAnswers: nil,
		Questions:          make([]*Question, 0),
		Answer:             nil,
		HasAnswer:          hasAnswer,
		ExpectedAnswer:     nil,
	}
	if hasAnswer {
		survey.ExpectedAnswer = &expectedAnswer
	}
	if hasMultipleAnswers {
		survey.HasMultipleAnswers = &hasMultipleAnswers
	}
	return append(surveys, survey)
}

func (fx *Fx) AppendQuestion(questions []*Question, text string, isBold bool, indent bool, hasAnswer bool, expectedAnswer bool) []*Question {
	question := &Question{
		Question:       text,
		IsBold:         isBold,
		Indent:         indent,
		Answer:         nil,
		HasAnswer:      hasAnswer,
		ExpectedAnswer: nil,
	}
	if hasAnswer {
		question.ExpectedAnswer = &expectedAnswer
	}

	return append(questions, question)
}

func (fx *Fx) HasGuaranteePolicy(input map[string]interface{}, guaranteeSlug string) bool {
	j, err := json.Marshal(input)
	lib.CheckError(err)
	var policy Policy
	err = json.Unmarshal(j, &policy)
	lib.CheckError(err)
	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == guaranteeSlug {
				return true
			}
		}
	}
	return false
}

func (fx *Fx) GetGuaranteeIndex(input map[string]interface{}, guranteeSlug string) int {
	j, _ := json.Marshal(input)
	var policy Policy
	_ = json.Unmarshal(j, &policy)
	for _, asset := range policy.Assets {
		for i, guarantee := range asset.Guarantees {
			if guarantee.Slug == guranteeSlug {
				return i
			}
		}
	}
	return -1
}

func (fx *Fx) DeleteOfferFromGuarantee(m map[string]*GuaranteValue, key string) {
	delete(m, key)
}

func (fx *Fx) CalculateMonthlyCoveragePrices(guarantees map[string]*Guarante) {
	for _, guarantee := range guarantees {
		for _, offer := range guarantee.Offer {
			offer.PremiumGrossMonthly = offer.PremiumGrossYearly / 12
			offer.PremiumTaxAmountMonthly = offer.PremiumTaxAmountYearly / 12
			offer.PremiumNetMonthly = offer.PremiumNetYearly / 12
		}
	}
}

func (fx *Fx) CalculateOfferPrices(guarantees map[string]*Guarante, offersPrices map[string]map[string]*Price) {
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

func (fx *Fx) RoundMonthlyOfferPrices(guarantees map[string]*Guarante, offerPrice map[string]map[string]*Price, roundingGuarantees ...string) {
	updatePrices := func(coverage string, offerType string, priceStruct map[string]*Price) {
		guarantees[coverage].Offer[offerType].PremiumGrossMonthly += priceStruct[monthly].Delta
		newNetPrice := guarantees[coverage].Offer[offerType].PremiumGrossMonthly / (1 + (guarantees[coverage].Tax / 100))
		newTax := guarantees[coverage].Offer[offerType].PremiumGrossMonthly - newNetPrice
		offerPrice[offerType][monthly].Net += newNetPrice - guarantees[coverage].Offer[offerType].PremiumNetMonthly
		offerPrice[offerType][monthly].Tax += newTax - guarantees[coverage].Offer[offerType].PremiumTaxAmountMonthly
		guarantees[coverage].Offer[offerType].PremiumNetMonthly = newNetPrice
		guarantees[coverage].Offer[offerType].PremiumTaxAmountMonthly = newTax
	}

	for offerType, priceStruct := range offerPrice {
		nonRoundedGrossPrice := priceStruct[yearly].Gross
		roundedMonthlyGrossPrice := math.Round(priceStruct[monthly].Gross)
		yearlyGrossPrice := roundedMonthlyGrossPrice * 12
		priceStruct[monthly].Delta = (yearlyGrossPrice - nonRoundedGrossPrice) / 12
		priceStruct[monthly].Gross = roundedMonthlyGrossPrice

		for _, roundingCoverage := range roundingGuarantees {
			hasGuarantee := guarantees[roundingCoverage].Offer[offerType].PremiumNetMonthly > 0
			if hasGuarantee {
				updatePrices(roundingCoverage, offerType, priceStruct)
				break
			}
		}
	}
}

func (fx *Fx) RoundYearlyOfferPrices(guarantees map[string]*Guarante, offerPrice map[string]map[string]*Price, roundingGuarantees ...string) {
	updatePrices := func(coverage string, offerType string, priceStruct map[string]*Price) {
		guarantees[coverage].Offer[offerType].PremiumGrossYearly += priceStruct[yearly].Delta
		newNetPrice := guarantees[coverage].Offer[offerType].PremiumGrossYearly / (1 + (guarantees[coverage].Tax / 100))
		newTax := guarantees[coverage].Offer[offerType].PremiumGrossYearly - newNetPrice
		offerPrice[offerType][yearly].Net += newNetPrice - guarantees[coverage].Offer[offerType].PremiumNetYearly
		offerPrice[offerType][yearly].Tax += newTax - guarantees[coverage].Offer[offerType].PremiumTaxAmountYearly
		guarantees[coverage].Offer[offerType].PremiumNetYearly = newNetPrice
		guarantees[coverage].Offer[offerType].PremiumTaxAmountYearly = newTax
	}

	for offerType, priceStruct := range offerPrice {
		ceilGrossPrice := math.Ceil(priceStruct[yearly].Gross)
		priceStruct[yearly].Delta = ceilGrossPrice - priceStruct[yearly].Gross
		priceStruct[yearly].Gross = ceilGrossPrice
		for _, roundingCoverage := range roundingGuarantees {
			hasGuarantee := guarantees[roundingCoverage].Offer[offerType].PremiumNetMonthly > 0
			if hasGuarantee {
				updatePrices(roundingCoverage, offerType, priceStruct)
				break
			}
		}
	}
}

func (fx *Fx) RoundToTwoDecimalPlaces(guarantees map[string]*Guarante, offersPrices map[string]map[string]*Price) {
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

func (fx *Fx) FilterOffersByMinimumPrice(guarantees map[string]*Guarante, offersPrices map[string]map[string]*Price, yearlyPriceMinimum float64, monthlyPriceMinimum float64) {
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
				if fx.HasGuaranteePerOffer(guarantees, offerKey, guarantee.Slug) {
					guarantee.Offer[offerKey].PremiumNetMonthly = 0.0
					guarantee.Offer[offerKey].PremiumTaxAmountMonthly = 0.0
					guarantee.Offer[offerKey].PremiumGrossMonthly = 0.0
				}
			}
		}
		if hasNotOfferMinimumYearlyPrice {
			delete(offersPrices[offerKey], yearly)
			for _, guarantee := range guarantees {
				if fx.HasGuaranteePerOffer(guarantees, offerKey, guarantee.Slug) {
					guarantee.Offer[offerKey].PremiumNetYearly = 0.0
					guarantee.Offer[offerKey].PremiumTaxAmountYearly = 0.0
					guarantee.Offer[offerKey].PremiumGrossYearly = 0.0
				}
			}
		}

	}
}

func (fx *Fx) HasGuarantee(guarantees map[string]*Guarante, guaranteeKey string) bool {
	for _, guarantee := range guarantees {
		if guarantee.Slug == guaranteeKey {
			return true
		}
	}
	return false
}

func (fx *Fx) HasGuaranteePerOffer(guarantees map[string]*Guarante, offerSlug string, guaranteeKey string) bool {
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

func (fx *Fx) RemoveGuaranteeIfCondition(guarantees map[string]*Guarante, guaranteeKey string, condition bool) {
	if condition {
		fx.RemoveGuarantee(guarantees, guaranteeKey)
	}
}

func (fx *Fx) RemoveOfferFromGuaranteeIfCondition(guaranteeOffer map[string]*GuaranteValue, offerKey string, condition bool) {
	if condition {
		fx.DeleteOfferFromGuarantee(guaranteeOffer, offerKey)
	}
}

func (fx *Fx) RemoveGuarantee(guarantees map[string]*Guarante, guaranteeKey string) {
	delete(guarantees, guaranteeKey)
}

func (fx *Fx) RemoveOfferFromGuarantee(guaranteeOffer map[string]*GuaranteValue, offerKey string) {
	delete(guaranteeOffer, offerKey)
}

func (fx *Fx) RemoveGuaranteesPriceZero(guarantees map[string]*Guarante) {
	for _, guarantee := range guarantees {
		for offerKey, _ := range guarantee.Offer {
			if guarantee.Offer[offerKey].PremiumGrossYearly == 0.0 {
				delete(guarantee.Offer, offerKey)
			}
		}
	}
}

func (fx *Fx) RemoveOfferPrice(offerPrice map[string]map[string]*Price, offerKey string) {
	delete(offerPrice, offerKey)
}
