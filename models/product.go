package models

type Product struct {
	ID        *string `json:"id,omitempty"`
	Company   *string `json:"company,omitempty"`
	Name      *string `json:"name,omitempty"`
	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`
}
