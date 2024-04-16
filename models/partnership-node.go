package models

import (
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

func (pn *PartnershipNode) IsJwtProtected() bool {
	c := pn.JwtConfig
	return (c.KeyAlgorithm != "" && c.ContentEncryption != "") || c.SignatureAlgorithm != ""
}

func (pn *PartnershipNode) DecryptJwt(jwtData, key string) ([]byte, error) {
	if !pn.IsJwtProtected() {
		return nil, nil
	}

	return lib.ParseJwt(jwtData, key, pn.JwtConfig)
}

func (pn PartnershipNode) DecryptJwtClaims(jwtData, key string, unmarshaler func([]byte) (LifeClaims, error)) (LifeClaims, error) {
	bytes, err := pn.DecryptJwt(jwtData, key)
	if err != nil {
		return LifeClaims{}, err
	}
	return unmarshaler(bytes)
}

type ClaimsGuarantee struct {
	Duration                   int     `json:"duration"`
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
}

type LifeClaims struct {
	Name       string                     `json:"name"`
	Surname    string                     `json:"surname"`
	BirthDate  string                     `json:"birthDate"`
	Gender     string                     `json:"gender"`
	FiscalCode string                     `json:"fiscalCode"`
	VatCode    string                     `json:"vatCode"`
	Email      string                     `json:"email"`
	Phone      string                     `json:"phone"`
	Address    string                     `json:"address"`
	Postalcode string                     `json:"postalCode"`
	City       string                     `json:"city"`
	CityCode   string                     `json:"cityCode"`
	Work       string                     `json:"work"`
	Guarantees map[string]ClaimsGuarantee `json:"guarantees"`
	Data       map[string]any             `json:"data"`
	jwt.RegisteredClaims
}

func (c *LifeClaims) IsEmpty() bool {
	return c.Data == nil
}
