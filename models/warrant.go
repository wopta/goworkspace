package models

import "time"

type Warrant struct {
	Name         string    `json:"name"`         // the name of the file saved in the bucket
	AllowedTypes []string  `json:"allowedTypes"` // the allowed NetworkNode types that can use the Warrant
	CreateDate   time.Time `json:"createDate"`   // when the Warrant was created ex.: "2023-10-10T00:00:00Z"
	StartDate    time.Time `json:"startDate"`    // the date when the Warrant becomes active
	EndDate      time.Time `json:"endDate"`      // the date when the Warrant becomes inactive
	Products     []Product `json:"products"`     // the list of product with their commisionsSettings for the given Warrant
}
