package models

import (
	"github.com/wopta/goworkspace/lib"
)

type PartnershipNode struct {
	Name      string        `json:"name" firestore:"name" bigquery:"name"`
	Skin      *Skin         `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
	JwtConfig lib.JwtConfig `json:"jwtConfig,omitempty" firestore:"jwtConfig,omitempty" bigquery:"-"` // DEPRECATED use root node jwtConfig
}

func (pn *PartnershipNode) Normalize() {
	pn.Name = lib.ToLower(pn.Name)
}
