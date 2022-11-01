package models

type User struct {
	Name         string `firestore:"name" json:"name,omitempty"`
	Type         string `firestore:"type" json:"type,omitempty"`
	Cluster      string `firestore:"cluster" json:"cluster,omitempty"`
	Surname      string `firestore:"surname" json:"surname,omitempty"`
	Address      string `firestore:"address" json:"address,omitempty"`
	PostalCode   string `firestore:"postalCode" json:"postalCode,omitempty"`
	Role         string `firestore:"role" json:"role,omitempty"`
	Work         string `firestore:"work" json:"work,omitempty"`
	WorkType     string `firestore:"workType" json:"workType,omitempty"`
	Mail         string `firestore:"mail" json:"mail,omitempty"`
	Phone        string `firestore:"phone" json:"phone,omitempty"`
	FiscalCode   string `firestore:"fiscalCode" json:"fiscalCode,omitempty"`
	VatCode      string `firestore:"vatCode" json:"vatCode,omitempty"`
	RiskClass    string `firestore:"riskClass" json:"riskClass,omitempty"`
	CreationDate string `firestore:"creationDate" json:"creationDate,omitempty"`
	UpdatedDate  string `firestore:"updatedDate" json:"updatedDate,omitempty"`
}
