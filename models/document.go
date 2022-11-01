package models

import "encoding/json"

func UnmarshalWelcome(data []byte) (Document, error) {
	var r Document
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Document) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Document struct {
	Insurance         string             `json:"insurance"`
	CompanyRef        string             `json:"companyRef"`
	Template          string             `json:"template"`
	Contractor        User               `json:"contractor"`
	Coverages         []Coverage         `json:"coverages"`
	Statements        []SpecialCondition `json:"statements"`
	SpecialConditions []SpecialCondition `json:"specialConditions"`
}

type SpecialCondition struct {
	A string `json:"a"`
	B string `json:"b"`
}
