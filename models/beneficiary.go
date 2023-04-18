package models

type Beneficiary struct {
	User
	IsFamilyMember         bool `json:"isFamilyMember" firestore:"isFamilyMember"`
	IsContactable          bool `json:"isContactable" firestore:"isContactable"`
	IsLegitimateSuccessors bool `json:"isLegitimateSuccessors" firestore:"isLegitimateSuccessors"`
}
