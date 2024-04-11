package models

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
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

func (pn PartnershipNode) DecryptJwtClaims2(jwtData, key string, unmarshaler func([]byte) (LifeClaims, error)) (LifeClaims, error) {
	bytes, err := pn.DecryptJwt(jwtData, key)
	if err != nil {
		return LifeClaims{}, err
	}
	return unmarshaler(bytes)
}

type LifeClaims struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	BirthDate  string `json:"birthDate"`
	Gender     string `json:"gender"`
	FiscalCode string `json:"fiscalCode"`
	VatCode    string `json:"vatCode"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Postalcode string `json:"postalCode"`
	City       string `json:"city"`
	CityCode   string `json:"cityCode"`
	Work       string `json:"work"`
	Guarantees map[string]struct {
		Duration                   int
		SumInsuredLimitOfIndemnity float64
	} `json:"guarantees"`
	Data map[string]any `json:"data"`
	jwt.RegisteredClaims
}
