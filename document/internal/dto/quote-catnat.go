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
type QuoteCatnatDTO struct {
	Sede                buildingCatnatDto
	EarthquakeGuarantee guaranteeCatnatDto
	FloodGuarantee      guaranteeCatnatDto
	LandslideGuarantee  guaranteeCatnatDto
	PaymentSplit        string
	Prize               priceDTO
}

func NewCatnatDto() QuoteCatnatDTO {
	return QuoteCatnatDTO{}
}

func (dto *QuoteCatnatDTO) FromPolicy(policy *models.Policy) {
	dto.Sede = buildingCatnatDto{}
	dto.Sede.buildingDTO = *newBuildingDTO()
	dto.Sede.buildingDTO.fromPolicy(*policy.Assets[0].Building, policy.Assets[0].Guarantees)
	dto.Sede.Type = useTypeMap[policy.Assets[0].Building.UseType]
	dto.Sede.BuildingMaterial = buildingMaterialMap[policy.Assets[0].Building.BuildingMaterial]
	dto.Sede.BuildingYear = buildingMaterialMap[policy.Assets[0].Building.BuildingYear]
	dto.Sede.LowestFloor = buildingMaterialMap[policy.Assets[0].Building.LowestFloor]
	dto.Sede.Floor = buildingMaterialMap[policy.Assets[0].Building.Floor]

	dto.EarthquakeGuarantee = newGuaranteeCatnatDto(policy, "EARTHQUAKE")
	dto.FloodGuarantee = newGuaranteeCatnatDto(policy, "FLOOD")
	dto.LandslideGuarantee = newGuaranteeCatnatDto(policy, "LANDSLIDE")

	dto.PaymentSplit = constants.PaymentSplitMap[policy.PaymentSplit]

	dto.Prize.Split = getSplit(policy.PaymentSplit)
	dto.Prize.Gross.ValueFloat = policy.PriceGross
	dto.Prize.Gross.Text = lib.HumanaizePriceEuro(policy.PriceGross)
}

func newGuaranteeCatnatDto(p *models.Policy, guarantee string) (res guaranteeCatnatDto) {
	var total float64
	for _, g := range p.Assets[0].Guarantees {
		if g.Group != guarantee {
			continue
		}
		if strings.HasSuffix(g.Name, "building") {
			res.Building = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.SumInsuredLimitOfIndemnity
		} else if strings.HasSuffix(g.Name, "content") {
			res.Content = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.SumInsuredLimitOfIndemnity
		} else if strings.HasSuffix(g.Name, "stock") {
			res.Stock = lib.HumanaizePriceEuro(float64(g.Value.SumInsuredLimitOfIndemnity))
			total += g.Value.SumInsuredLimitOfIndemnity
		}
	}
	res.Total = lib.HumanaizePriceEuro(total)

	return res
}
