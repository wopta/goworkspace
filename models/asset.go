package models

type Asset struct {
	Name      string
	Address   string
	Type      string
	Building  Building `json:"buildingType"`
	Coverages []Coverage
}
type Building struct {
	Name             string
	Address          string
	Type             string
	BuildingType     string `json:"buildingType"`
	BuildingMaterial string `json:"buildingMaterial"`
	BuildingYear     string `json:"buildingYear"`
	Employer         int64  `json:"employer"`
	IsAllarm         bool   `json:"isAllarm"`
	Floor            int64  `json:"floor"`
	IsPRA            bool   `json:"isPra"`
	Costruction      string `json:"costruction"`
	IsHolder         bool   `json:"isHolder"`
	Coverages        []Coverage
}
