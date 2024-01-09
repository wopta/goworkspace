package models

type Beneficiary struct {
	Name                   string   `firestore:"name"                        json:"name,omitempty"              bigquery:"-"`
	Surname                string   `firestore:"surname"                     json:"surname,omitempty"           bigquery:"-"`
	Mail                   string   `firestore:"mail"                        json:"mail,omitempty"              bigquery:"-"`
	Phone                  string   `firestore:"phone"                       json:"phone,omitempty"             bigquery:"-"`
	FiscalCode             string   `firestore:"fiscalCode,omitempty"                  json:"fiscalCode,omitempty"        bigquery:"-"`
	VatCode                string   `firestore:"vatCode,omitempty"                     json:"vatCode,omitempty"                     bigquery:"-"`
	Residence              *Address `firestore:"residence,omitempty"         json:"residence,omitempty"         bigquery:"-"`
	CompanyAddress         *Address `firestore:"companyAddress,omitempty" json:"companyAddress,omitempty" bigquery:"-"`
	IsFamilyMember         bool     `json:"isFamilyMember" firestore:"isFamilyMember"`
	IsContactable          bool     `json:"isContactable" firestore:"isContactable"`
	IsLegitimateSuccessors bool     `json:"isLegitimateSuccessors" firestore:"isLegitimateSuccessors"`
	BeneficiaryType        string   `json:"beneficiaryType" firestore:"beneficiaryType"`
}
