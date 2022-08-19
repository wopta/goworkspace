package models

type ProfileAllriskJson struct {
	Vat              int64  `json:"vat"`
	SquareMeters     int64  `json:"squareMeters"`
	IsBuildingOwner  bool   `json:"isBuildingOwner"`
	Revenue          int64  `json:"revenue"`
	Address          string `json:"address"`
	Ateco            string `json:"ateco"`
	BusinessSector   string `json:"businessSector"`
	BuildingType     string `json:"buildingType"`
	BuildingMaterial string `json:"buildingMaterial"`
	BuildingYear     string `json:"buildingYear"`
	Employer         int64  `json:"employer"`
	IsAllarm         bool   `json:"isAllarm"`
	Floor            int64  `json:"floor"`
	IsPRA            bool   `json:"isPra"`
	Costruction      string `json:"costruction"`
	IsHolder         bool   `json:"isHolder"`
	Result           string `json:"result"`
}
type ProfileAllrisk struct {
	Vat              int64
	SquareMeters     int64
	IsBuildingOwner  bool
	Revenue          int64
	Address          string
	Ateco            string
	BusinessSector   string
	BuildingType     string
	BuildingMaterial string
	BuildingYear     string
	Employer         int64
	IsAllarm         bool
	Floor            int64
	IsPRA            bool
	Costruction      string
	IsHolder         bool
	Result           string
}
