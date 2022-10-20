package models

type Product struct {
	Vat             int64  `json:"vat"`
	SquareMeters    int64  `json:"squareMeters"`
	IsBuildingOwner bool   `json:"isBuildingOwner"`
	Revenue         int64  `json:"revenue"`
	Address         string `json:"address"`
	Ateco           string `json:"ateco"`
	BusinessSector  string `json:"businessSector"`

	Result    string               `json:"result"`
	Coverages map[string]*Coverage `json:"coverages"`
}
