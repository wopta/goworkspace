package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
	Mail            string                `json:"mail" firestore:"mail" bigquery:"mail"`
	Warrant         string                `json:"warrant" firestore:"warrant" bigquery:"warrant"` // the name of the warrant file attached to the node
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
	AreaManager     *AgentNode            `json:"areaManager,omitempty" firestore:"areaManager,omitempty" bigquery:"areaManager"`
	Partnership     *PartnershipNode      `json:"partnership,omitempty" firestore:"partnership,omitempty" bigquery:"partnership"`
	NodeSetting     *NodeSetting          `json:"nodeSetting,omitempty" firestore:"nodeSetting,omitempty" bigquery:"-"` // Not implemented
	CreationDate    time.Time             `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	BigCreationDate bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"creationDate"`
	UpdatedDate     time.Time             `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
	BigUpdatedDate  bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"updatedDate"`
	Data            string                `json:"-" firestore:"-" bigquery:"data"`
	HasAnnex        bool                  `json:"hasAnnex" firestore:"hasAnnex" bigquery:"hasAnnex"`
	Designation     string                `json:"designation" firestore:"designation" bigquery:"designation"`
}

type PartnershipNode struct {
	Name string `json:"name" firestore:"name" bigquery:"name"`
	Skin *Skin  `json:"skin,omitempty" firestore:"skin,omitempty" bigquery:"-"`
}

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
}

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

func (nn *NetworkNode) Marshal() ([]byte, error) {
	return json.Marshal(nn)
}

func (nn *NetworkNode) SaveBigQuery(origin string) error {
	log.Println("[NetworkNode.SaveBigQuery]")

	nnJson, _ := json.Marshal(nn)

	nn.Data = string(nnJson)
	nn.BigCreationDate = lib.GetBigQueryNullDateTime(nn.CreationDate)
	nn.BigUpdatedDate = lib.GetBigQueryNullDateTime(nn.UpdatedDate)
	nn.Agent = parseBigQueryAgentNode(nn.Agent)
	nn.AreaManager = parseBigQueryAgentNode(nn.AreaManager)
	nn.Agency = parseBigQueryAgencyNode(nn.Agency)
	nn.Broker = parseBigQueryAgencyNode(nn.Broker)

	for _, p := range nn.Products {
		nn.BigProducts = append(nn.BigProducts, NodeProduct{
			Name:    p.Name,
			Version: p.Version,
		})
	}

	err := lib.InsertRowsBigQuery(WoptaDataset, NetworkNodesCollection, nn)
	return err
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

func (nn *NetworkNode) GetName() string {
	var name string

	// use constants
	switch nn.Type {
	case AgentNetworkNodeType:
		name = nn.Agent.Name + " " + nn.Agent.Surname
	case AgencyNetworkNodeType:
		name = nn.Agency.Name
	case BrokerNetworkNodeType:
		name = nn.Broker.Name
	case PartnershipNetworkNodeType:
		name = nn.Partnership.Name
	case AreaManagerNetworkNodeType:
		name = nn.AreaManager.Name + " " + nn.AreaManager.Surname
	}

	return strings.ReplaceAll(name, " ", "-")
}

func (nn *NetworkNode) GetWarrant() *Warrant {
	var (
		warrant *Warrant
	)

	if nn.Warrant == "" {
		log.Printf("[GetWarrant] warrant not set for node %s", nn.Uid)
		return nil
	}

	log.Printf("[GetWarrant] requesting warrant %s", nn.Warrant)

	warrantBytes := lib.GetFilesByEnv(fmt.Sprintf(WarrantFormat, nn.Warrant))

	err := json.Unmarshal(warrantBytes, &warrant)
	if err != nil {
		log.Printf("[GetWarrant] error unmarshaling warrant %s: %s", nn.Warrant, err.Error())
		return nil
	}

	return warrant
}

func (nn *NetworkNode) HasAccessToProduct(productName string, warrant *Warrant) bool {
	log.Println("[HasAccessToProduct] method start -----------------")

	needCheckTypes := []string{AgencyNetworkNodeType, AgentNetworkNodeType, BrokerNetworkNodeType}

	if !lib.SliceContains(needCheckTypes, nn.Type) {
		return true
	}

	if warrant == nil {
		warrant = nn.GetWarrant()
	}
	if warrant == nil {
		log.Printf("[HasAccessToProduct] no %s warrant found", nn.Warrant)
		return false
	}

	log.Printf("[HasAccessToProduct] checking if network node %s has access product %s", nn.Uid, productName)

	for _, product := range warrant.Products {
		if product.Name == productName {
			return true
		}
	}

	return false
}

func (nn *NetworkNode) GetNetworkNodeFlow(productName string, warrant *Warrant) (string, []byte) {
	if warrant == nil {
		log.Printf("[getNetworkNodeFlow] error warrant not set for node %s", nn.Uid)
		return "", []byte{}
	}

	product := warrant.GetProduct(productName)
	if product == nil {
		log.Printf("[getNetworkNodeFlow] error product not set for warrant %s", warrant.Name)
		return "", []byte{}
	}

	log.Printf("[getNetworkNodeFlow] getting flow '%s' file for product '%s'", product.Flow, productName)

	return product.Flow, lib.GetFilesByEnv(fmt.Sprintf(FlowFileFormat, product.Flow))
}
