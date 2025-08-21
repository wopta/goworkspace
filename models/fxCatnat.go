package models

import (
	"encoding/json"
	"strings"
)

type FxCatnat struct{}

func (FxCatnat) GetFiscalCodeAteco(input map[string]any) (res string) {
	j, _ := json.Marshal(input)
	var p Policy
	_ = json.Unmarshal(j, &p)
	if p.Contractor.FiscalCode != "" {
		res = p.Contractor.FiscalCode + "/"
	}
	if p.Contractor.VatCode != "" {
		res += p.Contractor.VatCode
	}
	return res
}

func (FxCatnat) GetBuildingPlace(input map[string]any) string {
	j, _ := json.Marshal(input)
	var p Policy
	_ = json.Unmarshal(j, &p)
	return strings.ToUpper(p.Assets[0].Building.BuildingAddress.StreetName + ", " + p.Assets[0].Building.BuildingAddress.StreetNumber + p.Assets[0].Building.BuildingAddress.PostalCode + " " + p.Assets[0].Building.BuildingAddress.City + " (" + p.Assets[0].Building.BuildingAddress.CityCode + ")")
}
