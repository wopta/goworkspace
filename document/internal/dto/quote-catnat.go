package dto

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

var useTypeMap = map[string]string{
	"owner-tenant": "Proprieta e conduttore",
	"tenant":       "Conduttore",
}

var quoteQuestionMap = map[bool]string{
	true:  "si",
	false: "no",
}

var buildingYearMap = map[string]string{
	"before_1950":       "Fino al 1950",
	"from_1950_to_1990": "1950-1990",
	"after_1990":        "Post 1990",
	"unknown":           "Non conosciuto",
}
var floorMap = map[string]string{
	"up_to_2":     "Da 0 a 2",
	"more_than_3": "3 o più",
}
var lowestFloorMap = map[string]string{
	"first_floor":  "Primo piano",
	"upper_floor":  "Superiori al primo",
	"ground_floor": "Piano terra / Piano strada",
	"underground":  "Scantinato",
}
var buildingMaterialMap = map[string]string{
	"brick":    "Muratura",
	"concrete": "Cemento armato",
	"steel":    "Acciao",
	"unknown":  "Non conosciuto / Altro",
}

type buildingCatnatDto struct {
	Type             string
	BuildingYear     string
	BuildingMaterial string
	Floor            string
	LowestFloor      string
	buildingDTO
}
type guaranteeCatnatDto struct {
	Content  string
	Building string
	Stock    string
	Total    string
}
type QuestionsCatnatDto struct {
	AlreadyEarthquake string
	AlreadyFlood      string
	WantEarthquake    string
	WantFlood         string
}
type QuoteCatnatDTO struct {
	Sede                buildingCatnatDto
	EarthquakeGuarantee guaranteeCatnatDto
	FloodGuarantee      guaranteeCatnatDto
	LandslideGuarantee  guaranteeCatnatDto
	PaymentSplit        string
	Prize               priceDTO
	Questions           QuestionsCatnatDto
}

func NewCatnatDto() QuoteCatnatDTO {
	return QuoteCatnatDTO{}
}

func (b *buildingCatnatDto) fromPolicy(policy *models.Policy) {
	b.buildingDTO = *newBuildingDTO()
	b.buildingDTO.fromPolicy(*policy.Assets[0].Building, policy.Assets[0].Guarantees)
	b.Type = useTypeMap[policy.Assets[0].Building.UseType]
	b.BuildingMaterial = buildingMaterialMap[policy.Assets[0].Building.BuildingMaterial]
	b.BuildingYear = buildingYearMap[policy.Assets[0].Building.BuildingYear]
	b.LowestFloor = lowestFloorMap[policy.Assets[0].Building.LowestFloor]
	b.Floor = floorMap[policy.Assets[0].Building.Floor]
}
func (dto *QuoteCatnatDTO) FromPolicy(policy *models.Policy) {
	dto.Sede = buildingCatnatDto{}
	dto.Sede.fromPolicy(policy)
	dto.EarthquakeGuarantee = newGuaranteeCatnatDto(policy, "EARTHQUAKE")
	dto.FloodGuarantee = newGuaranteeCatnatDto(policy, "FLOOD")
	dto.LandslideGuarantee = newGuaranteeCatnatDto(policy, "LANDSLIDE")

	dto.Questions = newQuestionCatnatDto(policy)

	dto.PaymentSplit = constants.PaymentSplitMap[policy.PaymentSplit]

	dto.Prize.Split = getSplit(policy.PaymentSplit)
	dto.Prize.Gross.ValueFloat = policy.PriceGross
	dto.Prize.Gross.Text = lib.HumanaizePriceEuro(policy.PriceGross)
	dto.Prize.Consultancy.ValueFloat = policy.ConsultancyValue.Price
	dto.Prize.Consultancy.Text = lib.HumanaizePriceEuro(policy.ConsultancyValue.Price)
	dto.Prize.Total.ValueFloat = policy.ConsultancyValue.Price + policy.PriceGross
	dto.Prize.Total.Text = lib.HumanaizePriceEuro(policy.ConsultancyValue.Price + policy.PriceGross)
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
			total += g.Value.SumInsuredLimitOfIndemnity
		} else if strings.HasSuffix(g.Slug, "content") {
			res.Content = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.SumInsuredLimitOfIndemnity
		} else if strings.HasSuffix(g.Slug, "stock") {
			res.Stock = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.SumInsuredLimitOfIndemnity
		}
	}
	res.Total = lib.HumanaizePriceEuro(total)

	return res
}
