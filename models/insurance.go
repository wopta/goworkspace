package models

type Insurance struct {
	Name         string
	Address      string
	Type         string
	policyholder User
	Building     Building `json:"buildingType"`
	Coverages    []Coverage
}
