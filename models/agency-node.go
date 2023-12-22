package models

import (
	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
	"time"
)

type AgencyNode struct {
	Name               string                `json:"name" firestore:"name" bigquery:"name"`
	VatCode            string                `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	Phone              string                `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	Address            *NodeAddress          `json:"address,omitempty" firestore:"address,omitempty" bigquery:"-"`
	Manager            *AgentNode            `json:"manager,omitempty" firestore:"manager,omitempty" bigquery:"manager"`
	RuiCode            string                `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection         string                `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration    time.Time             `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"ruiRegistration"`
	Skin               *Skin                 `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
	Pec                string                `json:"pec,omitempty" firestore:"pec,omitempty" bigquery:"-"`
	Website            string                `json:"website,omitempty" firestore:"website,omitempty" bigquery:"-"`
}

func (an *AgencyNode) Sanitize() {
	an.Name = lib.ToLower(an.Name)
	an.VatCode = lib.ToLower(an.VatCode)
	an.Phone = lib.TrimSpace(an.Phone)
	if an.Address != nil {
		an.Address.Sanitize()
	}
	if an.Manager != nil {
		an.Manager.Sanitize()
	}
	an.RuiCode = lib.ToLower(an.RuiCode)
	an.RuiSection = lib.ToLower(an.RuiSection)
	an.Pec = lib.ToLower(an.Pec)
	an.Website = lib.ToLower(an.Website)
}

func parseBigQueryAgencyNode(agency *AgencyNode) *AgencyNode {
	if agency == nil {
		return nil
	}

	if agency.Address != nil {
		agency.Address.BigLocation = lib.GetBigQueryNullGeography(
			agency.Address.Location.Lng,
			agency.Address.Location.Lat,
		)
	}
	agency.Manager = parseBigQueryAgentNode(agency.Manager)
	agency.BigRuiRegistration = lib.GetBigQueryNullDateTime(agency.RuiRegistration)

	return agency
}
