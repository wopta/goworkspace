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
	if na == nil {
		return
	}
	na.StreetName = lib.ToLower(na.StreetName)
	na.StreetNumber = lib.ToLower(na.StreetNumber)
	na.City = lib.ToLower(na.City)
	na.PostalCode = lib.TrimSpace(na.PostalCode)
	na.Locality = lib.ToLower(na.Locality)
	na.CityCode = lib.ToLower(na.CityCode)
	na.Area = lib.ToLower(na.Area)
}
