package models

import (
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
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
	Users           []string              `json:"users" firestore:"users" bigquery:"-"`
	Products        []Product             `json:"products" firestore:"products" bigquery:"-"`
	Policies        []string              `json:"policies" firestore:"policies" bigquery:"-"`
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
	Name                   string                 `json:"name" firestore:"name" bigquery:"name"`
	VatCode                string                 `json:"vatCode,omitempty" firestore:"vatCode,omitempty" bigquery:"vatCode"`
	Mail                   string                 `json:"mail" firestore:"mail" bigquery:"mail"`
	Phone                  string                 `json:"phone,omitempty" firestore:"phone,omitempty" bigquery:"phone"`
	Address                *Address               `json:"address,omitempty" firestore:"address,omitempty" bigquery:"-"`
	BigAddressStreetName   string                 `json:"-" firestore:"-" bigquery:"addressStreetName"`
	BigAddressStreetNumber string                 `json:"-" firestore:"-" bigquery:"addressStreetNumber"`
	BigAddressCity         string                 `json:"-" firestore:"-" bigquery:"addressCity"`
	BigAddressPostalCode   string                 `json:"-" firestore:"-" bigquery:"addressPostalCode"`
	BigAddressLocality     string                 `json:"-" firestore:"-" bigquery:"addressLocality"`
	BigAddressCityCode     string                 `json:"-" firestore:"-" bigquery:"addressCityCode"`
	Location               Location               `json:"location,omitempty" firestore:"location,omitempty" bigquery:"-"`
	BigLocation            bigquery.NullGeography `json:"-" firestore:"-" bigquery:"location"`
	RuiCode                string                 `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection             string                 `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration        time.Time              `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration     bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"ruiRegistration"`
	Skin                   *Skin                  `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
}

type AgentNode struct {
	Name                     string                 `json:"name" firestore:"name" bigquery:"name"`
	Surname                  string                 `json:"surname,omitempty" firestore:"surname,omitempty" bigquery:"surname"`
	FiscalCode               string                 `json:"fiscalCode,omitempty" firestore:"fiscalCode,omitempty" bigquery:"fiscalCode"`
	BirthDate                string                 `json:"birthDate,omitempty" firestore:"birthDate,omitempty" bigquery:"-"`
	BigBirthDate             bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"birthDate"`
	BirthCity                string                 `json:"birthCity,omitempty" firestore:"birthCity,omitempty" bigquery:"birthCity"`
	BirthProvince            string                 `json:"birthProvince,omitempty" firestore:"birthProvince,omitempty" bigquery:"birthProvince"`
	Residence                *Address               `json:"residence,omitempty" firestore:"residence,omitempty" bigquery:"-"`
	BigResidenceStreetName   string                 `json:"-" firestore:"-" bigquery:"residenceStreetName"`
	BigResidenceStreetNumber string                 `json:"-" firestore:"-" bigquery:"residenceStreetNumber"`
	BigResidenceCity         string                 `json:"-" firestore:"-" bigquery:"residenceCity"`
	BigResidencePostalCode   string                 `json:"-" firestore:"-" bigquery:"residencePostalCode"`
	BigResidenceLocality     string                 `json:"-" firestore:"-" bigquery:"residenceLocality"`
	BigResidenceCityCode     string                 `json:"-" firestore:"-" bigquery:"residenceCityCode"`
	Location                 Location               `json:"location,omitempty" firestore:"location,omitempty" bigquery:"-"`
	BigLocation              bigquery.NullGeography `json:"-" firestore:"-" bigquery:"location"`
	Domicile                 *Address               `json:"domicile,omitempty" firestore:"domicile,omitempty" bigquery:"-"`
	BigDomicileStreetName    string                 `json:"-" firestore:"-" bigquery:"domicileStreetName"`
	BigDomicileStreetNumber  string                 `json:"-" firestore:"-" bigquery:"domicileStreetNumber"`
	BigDomicileCity          string                 `json:"-" firestore:"-" bigquery:"domicileCity"`
	BigDomicilePostalCode    string                 `json:"-" firestore:"-" bigquery:"domicilePostalCode"`
	BigDomicileLocality      string                 `json:"-" firestore:"-" bigquery:"domicileLocality"`
	BigDomicileCityCode      string                 `json:"-" firestore:"-" bigquery:"domicileCityCode"`
	RuiCode                  string                 `json:"ruiCode" firestore:"ruiCode" bigquery:"ruiCode"`
	RuiSection               string                 `json:"ruiSection" firestore:"ruiSection" bigquery:"ruiSection"`
	RuiRegistration          time.Time              `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	BigRuiRegistration       bigquery.NullDateTime  `json:"-" firestore:"-" bigquery:"ruiRegistration"`
}

func (nn *NetworkNode) SaveBigQuery(origin string) error {
	log.Println("[NetworkNode.SaveBigQuery] TO BE IMPLEMENTED")
	return nil
}

func (nn *NetworkNode) GetName() string {
	var name string

	// use constants
	switch nn.Type {
	case "agent":
		name = nn.Agent.Name + " " + nn.Agent.Surname
	case "agency", "broker":
		name = nn.Agency.Name
	case "partnership":
		name = nn.Partnership.Name
	case "manager":
		name = "manager"
	}

	return strings.ReplaceAll(name, " ", "-")
}
