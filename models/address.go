package models

import "github.com/wopta/goworkspace/lib"

type Address struct {
	StreetName   string `json:"streetName" firestore:"streetName" bigquery:"-"`
	StreetNumber string `json:"streetNumber" firestore:"streetNumber" bigquery:"-"`
	City         string `json:"city" firestore:"city" bigquery:"-"`
	PostalCode   string `json:"postalCode" firestore:"postalCode" bigquery:"-"`
	Locality     string `json:"locality" firestore:"locality" bigquery:"-"`
	CityCode     string `json:"cityCode" firestore:"cityCode" bigquery:"-"`
	Area         string `json:"area" firestore:"area" bigquery:"-"`
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
