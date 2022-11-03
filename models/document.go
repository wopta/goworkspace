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

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()

type Document struct {
	Template          string             `json:"template"`
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
