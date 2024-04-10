package models

import (
	"encoding/json"

	"github.com/wopta/goworkspace/lib"
)

type PartnershipNode struct {
	Name      string        `json:"name" firestore:"name" bigquery:"name"`
	Skin      *Skin         `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
	JwtConfig lib.JwtConfig `json:"jwtConfig,omitempty" firestore:"jwtConfig,omitempty" bigquery:"-"`
}

func (pn *PartnershipNode) Normalize() {
	pn.Name = lib.ToLower(pn.Name)
}

func (pn *PartnershipNode) isJwtProtected() bool {
	c := pn.JwtConfig
	return (c.KeyAlgorithm != "" && c.ContentEncryption != "") || c.SignatureAlgorithm != ""
}

func (pn *PartnershipNode) DecryptJwt(jwtData, key string) ([]byte, error) {
	if !pn.isJwtProtected() {
		return nil, nil
	}

	return lib.ParseJwt(jwtData, key, pn.JwtConfig)
}

func (pn PartnershipNode) DecryptJwtClaims(jwtData, key string, claims any) error {
	bytes, err := pn.DecryptJwt(jwtData, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &claims)
}
