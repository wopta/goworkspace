package models

import (
	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

// Check if it's worth updating the Address model used by User
type NodeAddress struct {
	StreetName   string                 `json:"streetName" firestore:"streetName" bigquery:"streetName"`
	StreetNumber string                 `json:"streetNumber" firestore:"streetNumber" bigquery:"streetNumber"`
	City         string                 `json:"city" firestore:"city" bigquery:"city"`
	PostalCode   string                 `json:"postalCode" firestore:"postalCode" bigquery:"postalCode"`
	Locality     string                 `json:"locality" firestore:"locality" bigquery:"locality"`
	CityCode     string                 `json:"cityCode" firestore:"cityCode" bigquery:"cityCode"`
	Area         string                 `json:"area" firestore:"area" bigquery:"area"`
	Location     Location               `json:"location" firestore:"location" bigquery:"-"`
	BigLocation  bigquery.NullGeography `json:"-" firestore:"-" bigquery:"location"`
}

func (na *NodeAddress) Sanitize() {
	na.StreetName = lib.ToUpper(na.StreetName)
	na.StreetNumber = lib.ToUpper(na.StreetNumber)
	na.City = lib.ToUpper(na.City)
	na.PostalCode = lib.ToUpper(na.PostalCode)
	na.Locality = lib.ToUpper(na.Locality)
	na.CityCode = lib.ToUpper(na.CityCode)
	na.Area = lib.ToUpper(na.Area)
}
