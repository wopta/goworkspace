package models

type Asset struct {
	Name       string      `json:"Name"`
	Address    string      `json:"address"`
	Type       string      `json:"type"`
	Building   Building    `json:"building"`
	Person     User        `json:"person"`
	Guarantees []Guarantee `json:"guarantees,omitempty"`
}
type Building struct {
	Name             string `json:"Name"`
	Address          string `json:"address"`
	Type             string `json:"type"`
	BuildingType     string `json:"buildingType"`
	BuildingMaterial string `json:"buildingMaterial"`
	BuildingYear     string `json:"buildingYear"`
	Employer         int64  `json:"employer"`
	IsAllarm         bool   `json:"isAllarm"`
	Floor            int64  `json:"floor"`
	IsPRA            bool   `json:"isPra"`
	Costruction      string `json:"costruction"`
	IsHolder         bool   `json:"isHolder"`
}
