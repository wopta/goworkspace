package models

import (
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/wopta/goworkspace/lib"
)

type NetworkNode struct {
	Uid             string                `json:"uid" firestore:"uid" bigquery:"uid"`
	AuthId          string                `json:"authId,omitempty" firestore:"authId,omitempty" bigquery:"-"`
	Code            string                `json:"code" firestore:"code" bigquery:"code"`
	Type            string                `json:"type" firestore:"type" bigquery:"type"`
	Role            string                `json:"role" firestore:"role" bigquery:"role"`
	NetworkUid      string                `json:"networkUid" firestore:"networkUid" bigquery:"networkUid"`
	NetworkCode     string                `json:"networkCode" firestore:"networkCode" bigquery:"networkCode"`
	ParentUid       string                `json:"parentUid,omitempty" firestore:"parentUid,omitempty" bigquery:"parentUid"`
	ManagerUid      string                `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"managerUid"`
	IsActive        bool                  `json:"isActive" firestore:"isActive" bigquery:"isActive"`
	Users           []string              `json:"users" firestore:"users" bigquery:"users"`
	Products        []Product             `json:"products" firestore:"products" bigquery:"-"`
	BigProducts     []NodeProduct         `json:"-" firestore:"-" bigquery:"products"`
	Policies        []string              `json:"policies" firestore:"policies" bigquery:"policies"`
	Agent           *AgentNode            `json:"agent,omitempty" firestore:"agent,omitempty" bigquery:"agent"`
	Agency          *AgencyNode           `json:"agency,omitempty" firestore:"agency,omitempty" bigquery:"agency"`
	Broker          *AgencyNode           `json:"broker,omitempty" firestore:"broker,omitempty" bigquery:"broker"`
	Partnership     *PartnershipNode      `json:"partnership,omitempty" firestore:"partnership,omitempty" bigquery:"partnership"`
	NodeSetting     *NodeSetting          `json:"nodeSetting,omitempty" firestore:"nodeSetting,omitempty" bigquery:"-"`
	CreationDate    time.Time             `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	BigCreationDate bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"creationDate"`
	UpdatedDate     time.Time             `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
	BigUpdatedDate  bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"updatedDate"`
	Data            string                `json:"-" firestore:"-" bigquery:"data"`
}

type PartnershipNode struct {
	Name string `json:"name" firestore:"name" bigquery:"name"`
	Skin *Skin  `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
}

type AgencyNode struct {
	Name               string                `json:"name" firestore:"name" bigquery:"name"`
	VatCode            string                `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	Mail               string                `json:"mail" firestore:"mail" bigquery:"mail"`
	Phone              string                `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	Address            *NodeAddress          `json:"address,omitempty" firestore:"address,omitempty" bigquery:"-"`
	RuiCode            string                `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection         string                `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration    time.Time             `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"ruiRegistration"`
	Skin               *Skin                 `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
}

type AgentNode struct {
	Name               string                `json:"name" firestore:"name" bigquery:"name"`
	Surname            string                `json:"surname,omitempty" firestore:"surname,omitempty" bigquery:"surname"`
	FiscalCode         string                `json:"fiscalCode,omitempty" firestore:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	Mail               string                `json:"mail" firestore:"mail" bigquery:"mail"`
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

// Check if it's worth updating the Address model used by User
type NodeAddress struct {
	StreetName   string                 `json:"streetName,omitempty" firestore:"streetName" bigquery:"streetName"`
	StreetNumber string                 `json:"streetNumber,omitempty" firestore:"streetNumber,omitempty" bigquery:"streetNumber"`
	City         string                 `json:"city,omitempty" firestore:"city" bigquery:"city"`
	PostalCode   string                 `json:"postalCode,omitempty" firestore:"postalCode" bigquery:"postalCode"`
	Locality     string                 `json:"locality,omitempty" firestore:"locality" bigquery:"locality"`
	CityCode     string                 `json:"cityCode,omitempty" firestore:"cityCode" bigquery:"cityCode"`
	Area         string                 `json:"area,omitempty" firestore:"area,omitempty" bigquery:"area"`
	Location     Location               `json:"location,omitempty" firestore:"location,omitempty" bigquery:"-"`
	BigLocation  bigquery.NullGeography `json:"-" firestore:"-" bigquery:"location"`
}

type NodeProduct struct {
	Name    string `json:"-" firestore:"-" bigquery:"name"`
	Version string `json:"-" firestore:"-" bigquery:"version"`
}

func (nn *NetworkNode) SaveBigQuery(origin string) error {
	log.Println("[NetworkNode.SaveBigQuery]")

	nnJson, _ := json.Marshal(nn)

	nn.Data = string(nnJson)
	nn.BigCreationDate = lib.GetBigQueryNullDateTime(nn.CreationDate)
	nn.BigUpdatedDate = lib.GetBigQueryNullDateTime(nn.UpdatedDate)
	nn.parseBigQueryAgent()
	nn.parseBigQueryAgency()

	for _, p := range nn.Products {
		nn.BigProducts = append(nn.BigProducts, NodeProduct{
			Name:    p.Name,
			Version: p.Version,
		})
	}

	err := lib.InsertRowsBigQuery(WoptaDataset, NetworkNodeCollection, nn)
	return err
}

func (nn *NetworkNode) parseBigQueryAgent() {
	if nn.Agent == nil {
		return
	}

	if nn.Agent.BirthDate != "" {
		birthDate, _ := time.Parse(time.RFC3339, nn.Agent.BirthDate)
		nn.Agent.BigBirthDate = lib.GetBigQueryNullDateTime(birthDate)
	}
	if nn.Agent.Residence != nil {
		nn.Agent.Residence.BigLocation = lib.GetBigQueryNullGeography(
			nn.Agent.Residence.Location.Lng,
			nn.Agent.Residence.Location.Lat,
		)
	}
	if nn.Agent.Domicile != nil {
		nn.Agent.Domicile.BigLocation = lib.GetBigQueryNullGeography(
			nn.Agent.Domicile.Location.Lng,
			nn.Agent.Domicile.Location.Lat,
		)
	}
	nn.Agent.BigRuiRegistration = lib.GetBigQueryNullDateTime(nn.Agent.RuiRegistration)
}

func (nn *NetworkNode) parseBigQueryAgency() {
	if nn.Agency == nil {
		return
	}

	if nn.Agency.Address != nil {
		nn.Agency.Address.BigLocation = lib.GetBigQueryNullGeography(
			nn.Agency.Address.Location.Lng,
			nn.Agency.Address.Location.Lat,
		)
	}
	nn.Agency.BigRuiRegistration = lib.GetBigQueryNullDateTime(nn.Agency.RuiRegistration)
}
