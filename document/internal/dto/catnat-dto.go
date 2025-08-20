package dto

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CatnatDTO struct {
	SedeDaAssicurare BuildingCatnatDto
	Contractor       contractorDTO
	ValidityDate     *validityDateDTO
	Questions        QuestionsCatnatDto
	Guarantee        CatnatGuaranteeDTO
	Price            priceDTO
}

type CatnatGuaranteeDTO struct {
	EarthquakeGuarantee guaranteeCatnatDto
	FloodGuarantee      guaranteeCatnatDto
	LandslideGuarantee  guaranteeCatnatDto
}

func (c *CatnatGuaranteeDTO) fromPolicy(policy *models.Policy) {
	c.EarthquakeGuarantee = newGuaranteeCatnatDto(policy, "EARTHQUAKE")
	c.FloodGuarantee = newGuaranteeCatnatDto(policy, "FLOOD")
	c.LandslideGuarantee = newGuaranteeCatnatDto(policy, "LANDSLIDE")
}
func (dto *CatnatDTO) FromPolicy(policy *models.Policy) {
	dto.SedeDaAssicurare = BuildingCatnatDto{}
	dto.SedeDaAssicurare.fromPolicy(policy)
	dto.Contractor = contractorDTO{}
	dto.Contractor.fromPolicy(policy.Contractor)
	dto.ValidityDate = &validityDateDTO{}
	dto.ValidityDate.fromPolicy(policy)
	dto.Questions = newQuestionCatnatDto(policy)
	dto.Guarantee = CatnatGuaranteeDTO{}
	dto.Guarantee.fromPolicy(policy)
	dto.Price = priceDTO{}
	dto.Price.fromPolicy(*policy)
}

func newGuaranteeCatnatDto(p *models.Policy, guarantee string) (res guaranteeCatnatDto) {
	res.Building = "===="
	res.Content = "===="
	res.Stock = "===="
	var total float64
	for _, g := range p.Assets[0].Guarantees {
		if g.Group != guarantee {
			continue
		}
		if g.Value.SumInsuredLimitOfIndemnity == 0 {
			continue
		}
		if strings.HasSuffix(g.Slug, "building") {
			res.Building = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.PremiumGrossYearly
		} else if strings.HasSuffix(g.Slug, "content") {
			res.Content = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.PremiumGrossYearly
		} else if strings.HasSuffix(g.Slug, "stock") {
			res.Stock = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.PremiumGrossYearly
		}
	}
	res.Total = lib.HumanaizePriceEuro(total)

	return res
}
func newQuestionCatnatDto(p *models.Policy) (res QuestionsCatnatDto) {
	res.AlreadyEarthquake = "===="
	res.AlreadyFlood = "===="
	res.WantFlood = "===="
	res.WantEarthquake = "===="

	var alreadyEarthquake any
	var alreadyFlood any
	var wantEarthquake any
	var wantFlood any

	if p.Assets[0].Building.UseType == "tenant" {
		alreadyEarthquake = p.QuoteQuestions["alreadyEarthquake"]
		alreadyFlood = p.QuoteQuestions["alreadyFlood"]
		wantEarthquake = p.QuoteQuestions["wantEarthquake"]
		wantFlood = p.QuoteQuestions["wantFlood"]

		res.AlreadyEarthquake = quoteQuestionMap[alreadyEarthquake.(bool)]
		res.AlreadyFlood = quoteQuestionMap[alreadyFlood.(bool)]
		if wantEarthquake != nil {
			res.WantEarthquake = quoteQuestionMap[wantEarthquake.(bool)]
		}
		if wantFlood != nil {
			res.WantFlood = quoteQuestionMap[wantFlood.(bool)]
		}
	}

	return res

}
