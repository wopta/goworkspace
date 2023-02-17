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
	Uid               string             `json:"uid,omitempty"`
	Template          string             `json:"template,omitempty"`
	GcsFilename       string             `json:"gcsFilename,omitempty"`
	Policy            *Policy            `json:"policy,omitempty"`
	Contractor        *User              `json:"contractor,omitempty"`
	Coverages         []Guarantee        `json:"coverages,omitempty"`
	Statements        []SpecialCondition `json:"statements,omitempty"`
	SpecialConditions []SpecialCondition `json:"specialConditions,omitempty"`
}

type SpecialCondition struct {
	A *string `json:"a,omitempty"`
	B *string `json:"b,omitempty"`
}
