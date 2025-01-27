package dto

import (
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/models"
)

type buildingDTO struct {
	StreetName       string
	StreetNumber     string
	City             string
	PostalCode       string
	CityCode         string
	BuildingMaterial string
	HasSandwichPanel string
	HasAlarm         string
	HasSprinkler     string
	Naics            string
	NaicsDetail      string
	Guarantees       map[string]*guaranteeDTO
}

func newBuildingDTO() *buildingDTO {
	return &buildingDTO{
		StreetName:       constants.EmptyField,
		StreetNumber:     constants.EmptyField,
		City:             constants.EmptyField,
		PostalCode:       constants.EmptyField,
		CityCode:         constants.EmptyField,
		BuildingMaterial: constants.EmptyField,
		HasSandwichPanel: constants.EmptyField,
		HasAlarm:         constants.EmptyField,
		HasSprinkler:     constants.EmptyField,
		Naics:            constants.EmptyField,
		NaicsDetail:      constants.EmptyField,
		Guarantees:       make(map[string]*guaranteeDTO),
	}
}

func (b *buildingDTO) fromPolicy(building models.Building, guarantees []models.Guarante) {
	if building.BuildingAddress.StreetName != "" {
		b.StreetName = building.BuildingAddress.StreetName
	}
	if building.BuildingAddress.StreetNumber != "" {
		b.StreetNumber = building.BuildingAddress.StreetNumber
	}
	if building.BuildingAddress.City != "" {
		b.City = building.BuildingAddress.City
	}
	if building.BuildingAddress.PostalCode != "" {
		b.PostalCode = building.BuildingAddress.PostalCode
	}
	if building.BuildingAddress.CityCode != "" {
		b.CityCode = building.BuildingAddress.CityCode
	}
	if building.BuildingMaterial != "" {
		b.BuildingMaterial = building.BuildingMaterial
	}
	if building.HasSandwichPanel {
		b.HasSandwichPanel = yes
	} else {
		b.HasSandwichPanel = no
	}
	if building.HasAlarm {
		b.HasAlarm = yes
	} else {
		b.HasAlarm = no
	}
	if building.HasSprinkler {
		b.HasSprinkler = yes
	} else {
		b.HasSprinkler = no
	}
	if building.Naics != "" {
		b.Naics = building.Naics
	}
	if building.NaicsDetail != "" {
		b.NaicsDetail = building.NaicsDetail
	}

	for _, guarantee := range guarantees {
		b.Guarantees[guarantee.Slug].fromPolicy(guarantee)
	}
}
