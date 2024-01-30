package models

import (
	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
	"time"
)

type AgentNode struct {
	Name               string                `json:"name" firestore:"name" bigquery:"name"`
	Surname            string                `json:"surname,omitempty" firestore:"surname,omitempty" bigquery:"surname"`
	FiscalCode         string                `json:"fiscalCode,omitempty" firestore:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	VatCode            string                `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	Phone              string                `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	BirthDate          string                `json:"birthDate,omitempty" firestore:"birthDate,omitempty" bigquery:"-"`
	BigBirthDate       bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"birthDate"`
	BirthCity          string                `json:"birthCity,omitempty" firestore:"birthCity,omitempty" bigquery:"birthCity"`
	BirthProvince      string                `json:"birthProvince,omitempty" firestore:"birthProvince,omitempty" bigquery:"birthProvince"`
	Residence          *NodeAddress          `json:"residence,omitempty" firestore:"residence,omitempty" bigquery:"residence"`
	Domicile           *NodeAddress          `json:"domicile,omitempty" firestore:"domicile,omitempty" bigquery:"domicile"`
	RuiCode            string                `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection         string                `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration    time.Time             `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"ruiRegistration"`
}

func (an *AgentNode) Normalize() {
	an.Name = lib.ToUpper(an.Name)
	an.Surname = lib.ToUpper(an.Surname)
	an.FiscalCode = lib.ToUpper(an.FiscalCode)
	an.VatCode = lib.ToUpper(an.VatCode)
	an.BirthCity = lib.ToUpper(an.BirthCity)
	an.BirthProvince = lib.ToUpper(an.BirthProvince)
	if an.Residence != nil {
		an.Residence.Normalize()
	}
	if an.Domicile != nil {
		an.Domicile.Normalize()
	}
	an.Phone = lib.TrimSpace(an.Phone)
	an.RuiCode = lib.ToUpper(an.RuiCode)
	an.RuiSection = lib.ToUpper(an.RuiSection)
}

func parseBigQueryAgentNode(agent *AgentNode) *AgentNode {
	if agent == nil {
		return nil
	}

	if agent.BirthDate != "" {
		birthDate, _ := time.Parse(time.RFC3339, agent.BirthDate)
		agent.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}
	if agent.Residence != nil {
		agent.Residence.BigLocation = lib.GetBigQueryNullGeography(
			agent.Residence.Location.Lng,
			agent.Residence.Location.Lat,
		)
	}
	if agent.Domicile != nil {
		agent.Domicile.BigLocation = lib.GetBigQueryNullGeography(
			agent.Domicile.Location.Lng,
			agent.Domicile.Location.Lat,
		)
	}
	agent.BigRuiRegistration = lib.GetBigQueryNullDateTime(agent.RuiRegistration)

	return agent
}
