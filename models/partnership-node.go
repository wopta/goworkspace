package models

import "github.com/wopta/goworkspace/lib"

type PartnershipNode struct {
	Name string `json:"name" firestore:"name" bigquery:"name"`
	Skin *Skin  `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
}

func (pn *PartnershipNode) Sanitize() {
	pn.Name = lib.ToLower(pn.Name)
}
