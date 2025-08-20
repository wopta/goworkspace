package dto

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
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

type BuildingCatnatDto struct {
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
	Sede         BuildingCatnatDto
	PaymentSplit string
	Prize        priceDTO
	Questions    QuestionsCatnatDto
	Guarantees   CatnatGuaranteeDTO
}

func NewCatnatDto() QuoteCatnatDTO {
	return QuoteCatnatDTO{}
}

func (b *BuildingCatnatDto) fromPolicy(policy *models.Policy) {
	b.buildingDTO = *newBuildingDTO()
	b.buildingDTO.fromPolicy(*policy.Assets[0].Building, policy.Assets[0].Guarantees)
	b.Type = useTypeMap[policy.Assets[0].Building.UseType]
	b.BuildingMaterial = buildingMaterialMap[policy.Assets[0].Building.BuildingMaterial]
	b.BuildingYear = buildingYearMap[policy.Assets[0].Building.BuildingYear]
	b.LowestFloor = lowestFloorMap[policy.Assets[0].Building.LowestFloor]
	b.Floor = floorMap[policy.Assets[0].Building.Floor]

}
func (dto *QuoteCatnatDTO) FromPolicy(policy *models.Policy) {
	dto.Sede = BuildingCatnatDto{}
	dto.Sede.fromPolicy(policy)

	dto.Questions = newQuestionCatnatDto(policy)

	dto.PaymentSplit = constants.PaymentSplitMap[policy.PaymentSplit]

	dto.Prize = priceDTO{}
	dto.Prize.fromPolicy(*policy)

	dto.Guarantees = CatnatGuaranteeDTO{}
	dto.Guarantees.fromPolicy(policy)
}
