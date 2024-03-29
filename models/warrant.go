package models

import "time"

type Warrant struct {
	Name         string    `json:"name"`         // the name of the file saved in the bucket
	Description  string    `json:"description"`  // the description for the mandate use: which types of nodes are allowed, network which is assigned, products included etc
	AllowedTypes []string  `json:"allowedTypes"` // the allowed NetworkNode types that can use the Warrant
	CreateDate   time.Time `json:"createDate"`   // when the Warrant was created ex.: "2023-10-10T00:00:00Z"
	StartDate    time.Time `json:"startDate"`    // the date when the Warrant becomes active
	EndDate      time.Time `json:"endDate"`      // the date when the Warrant becomes inactive
	Products     []Product `json:"products"`     // the list of product with their commisionsSettings for the given Warrant
}

func (w *Warrant) GetProduct(productName string) *Product {
	for _, product := range w.Products {
		if product.Name == productName {
			return &product
		}
	}
	return nil
}

func (w *Warrant) GetFlowName(productName string) string {
	var flowName string
	product := w.GetProduct(productName)

	if product != nil {
		flowName = product.Flow
	}

	return flowName
}

func (w *Warrant) HasProductByName(productName string) bool {
	return w.GetProduct(productName) != nil
}
