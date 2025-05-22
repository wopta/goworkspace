package models

import "gitlab.dev.wopta.it/goworkspace/lib"

type Address struct {
	StreetName    string `json:"streetName" firestore:"streetName" bigquery:"-"`
	StreetNumber  string `json:"streetNumber" firestore:"streetNumber" bigquery:"-"`
	City          string `json:"city" firestore:"city" bigquery:"-"`
	PostalCode    string `json:"postalCode" firestore:"postalCode" bigquery:"-"`
	Locality      string `json:"locality" firestore:"locality" bigquery:"-"`
	CityCode      string `json:"cityCode" firestore:"cityCode" bigquery:"-"`
	Area          string `json:"area" firestore:"area" bigquery:"-"`
	IsManualInput bool   `json:"isManualInput" firestore:"isManualInput" bigquery:"-"`
}

func (a *Address) Normalize() {
	a.StreetName = lib.ToUpper(a.StreetName)
	a.StreetNumber = lib.ToUpper(a.StreetNumber)
	a.City = lib.ToUpper(a.City)
	a.PostalCode = lib.ToUpper(a.PostalCode)
	a.Locality = lib.ToUpper(a.Locality)
	a.CityCode = lib.ToUpper(a.CityCode)
	a.Area = lib.ToUpper(a.Area)
}

func (a Address) IsEqual(address Address) bool {
	if a.StreetName != address.StreetName {
		return false
	}
	if a.StreetNumber != address.StreetNumber {
		return false
	}
	if a.City != address.City {
		return false
	}
	if a.PostalCode != address.PostalCode {
		return false
	}
	if a.Locality != address.Locality {
		return false
	}
	if a.CityCode != address.CityCode {
		return false
	}
	if a.Area != address.Area {
		return false
	}
	if a.IsManualInput != address.IsManualInput {
		return false
	}
	return true
}
