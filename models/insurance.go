package models

type Insurance struct {
	Name      string
	Address   string
	Type      string
	Building  Building `json:"buildingType"`
	Coverages []Coverage
}
