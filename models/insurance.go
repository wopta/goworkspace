package models

type Insurance struct {
	Name       string
	Address    string
	Type       string
	Contractor User
	Assets     []Asset
	Coverages  []Coverage
}
