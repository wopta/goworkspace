package models

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/wopta/goworkspace/lib"
)

type NetworkNode struct {
	Uid                 string                `json:"uid" firestore:"uid" bigquery:"uid"`
	AuthId              string                `json:"authId,omitempty" firestore:"authId,omitempty" bigquery:"-"`
	Code                string                `json:"code" firestore:"code" bigquery:"code"`
	ExternalNetworkCode string                `json:"externalNetworkCode" firestore:"externalNetworkCode" bigquery:"externalNetworkCode"`
	Type                string                `json:"type" firestore:"type" bigquery:"type"`
	Role                string                `json:"role" firestore:"role" bigquery:"role"`
	Mail                string                `json:"mail" firestore:"mail" bigquery:"mail"`
	Warrant             string                `json:"warrant" firestore:"warrant" bigquery:"warrant"`             // the name of the warrant file attached to the node
	NetworkUid          string                `json:"networkUid" firestore:"networkUid" bigquery:"networkUid"`    // DEPRECATED
	NetworkCode         string                `json:"networkCode" firestore:"networkCode" bigquery:"networkCode"` // DEPRECATED
	ParentUid           string                `json:"parentUid,omitempty" firestore:"parentUid,omitempty" bigquery:"parentUid"`
	ManagerUid          string                `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"managerUid"` // DEPRECATED
	IsActive            bool                  `json:"isActive" firestore:"isActive" bigquery:"isActive"`
	Users               []string              `json:"users" firestore:"users" bigquery:"users"`
	Products            []Product             `json:"products" firestore:"products" bigquery:"-"`
	BigProducts         []NodeProduct         `json:"-" firestore:"-" bigquery:"products"`
	Policies            []string              `json:"policies" firestore:"policies" bigquery:"policies"`
	Agent               *AgentNode            `json:"agent,omitempty" firestore:"agent,omitempty" bigquery:"agent"`
	Agency              *AgencyNode           `json:"agency,omitempty" firestore:"agency,omitempty" bigquery:"agency"`
	Broker              *AgencyNode           `json:"broker,omitempty" firestore:"broker,omitempty" bigquery:"broker"`
	AreaManager         *AgentNode            `json:"areaManager,omitempty" firestore:"areaManager,omitempty" bigquery:"areaManager"`
	Partnership         *PartnershipNode      `json:"partnership,omitempty" firestore:"partnership,omitempty" bigquery:"partnership"`
	NodeSetting         *NodeSetting          `json:"nodeSetting,omitempty" firestore:"nodeSetting,omitempty" bigquery:"-"` // Not implemented
	CreationDate        time.Time             `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	BigCreationDate     bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"creationDate"`
	UpdatedDate         time.Time             `json:"updatedDate" firestore:"updatedDate" bigquery:"-"`
	BigUpdatedDate      bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"updatedDate"`
	Data                string                `json:"-" firestore:"-" bigquery:"data"`
	HasAnnex            bool                  `json:"hasAnnex" firestore:"hasAnnex" bigquery:"hasAnnex"`
	Designation         string                `json:"designation" firestore:"designation" bigquery:"designation"`
	IsMgaProponent      bool                  `json:"isMgaProponent" firestore:"isMgaProponent" bigquery:"-"`
	WorksForUid         string                `json:"worksForUid" firestore:"worksForUid" bigquery:"-"`
}

type NodeProduct struct {
	Name      string        `json:"-" firestore:"-" bigquery:"name"`
	Companies []NodeCompany `json:"-" firestore:"-" bigquery:"companies"`
}

type NodeCompany struct {
	Name         string `json:"-" firestore:"-" bigquery:"name"`
	ProducerCode string `json:"-" firestore:"-" bigquery:"producerCode"`
}

func NetworkNodeToListData(query *firestore.DocumentIterator) []NetworkNode {
	result := make([]NetworkNode, 0)
	for {
		d, err := query.Next()
		if err != nil {
		}
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value NetworkNode
		e := d.DataTo(&value)
		value.Uid = d.Ref.ID
		lib.CheckError(e)
		result = append(result, value)
	}
	return result
}

func (nn *NetworkNode) Marshal() ([]byte, error) {
	return json.Marshal(nn)
}

func (nn *NetworkNode) Normalize() {
	nn.Code = lib.TrimSpace(nn.Code)
	nn.ExternalNetworkCode = lib.TrimSpace(nn.ExternalNetworkCode)
	nn.Type = lib.TrimSpace(nn.Type)
	nn.Role = lib.TrimSpace(nn.Role)
	nn.Mail = lib.ToUpper(nn.Mail)
	nn.Warrant = lib.TrimSpace(nn.Warrant)
	nn.ParentUid = lib.TrimSpace(nn.ParentUid)
	nn.Designation = lib.TrimSpace(nn.Designation)
	nn.WorksForUid = lib.TrimSpace(nn.WorksForUid)

	switch nn.Type {
	case AgentNetworkNodeType:
		nn.Agent.Normalize()
	case AgencyNetworkNodeType:
		nn.Agency.Normalize()
	case BrokerNetworkNodeType:
		nn.Broker.Normalize()
	case AreaManagerNetworkNodeType:
		nn.AreaManager.Normalize()
	}
}

func (nn *NetworkNode) SaveFirestore() error {
	err := lib.SetFirestoreErr(NetworkNodesCollection, nn.Uid, nn)
	if err != nil {
		log.Printf("[NetworkNode.SaveFirestore] error: %s", err.Error())
	}
	return err
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
		companies := make([]NodeCompany, 0)
		for _, c := range p.Companies {
			companies = append(companies, NodeCompany{
				Name:         c.Name,
				ProducerCode: c.ProducerCode,
			})
		}
		nn.BigProducts = append(nn.BigProducts, NodeProduct{
			Name:      p.Name,
			Companies: companies,
		})
	}

	err := lib.InsertRowsBigQuery(WoptaDataset, NetworkNodesCollection, nn)
	return err
}

func (nn *NetworkNode) GetName() string {
	var name string

	switch nn.Type {
	case AgentNetworkNodeType:
		name = fmt.Sprintf("%s %s", lib.Capitalize(nn.Agent.Name), lib.Capitalize(nn.Agent.Surname))
	case AgencyNetworkNodeType:
		name = lib.Capitalize(nn.Agency.Name)
	case BrokerNetworkNodeType:
		name = lib.Capitalize(nn.Broker.Name)
	case PartnershipNetworkNodeType:
		name = lib.Capitalize(nn.Partnership.Name)
	case AreaManagerNetworkNodeType:
		name = fmt.Sprintf("%s %s", lib.Capitalize(nn.AreaManager.Name), lib.Capitalize(nn.AreaManager.Surname))
	}

	return name
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

func (nn *NetworkNode) GetAddress() string {
	var (
		address       string
		addressFormat = "%s, %s - %s %s (%s)"
	)

	switch nn.Type {
	case AgencyNetworkNodeType:
		return fmt.Sprintf(
			addressFormat,
			nn.Agency.Address.StreetName,
			nn.Agency.Address.StreetNumber,
			nn.Agency.Address.PostalCode,
			nn.Agency.Address.City,
			nn.Agency.Address.CityCode,
		)
	}

	return address
}

func (nn *NetworkNode) GetAncestors() ([]NetworkTreeElement, error) {
	var (
		err error
	)

	query := fmt.Sprintf("SELECT rootUid, ntr.parentUid, nodeUid, nnv.name AS name, relativeLevel, "+
		"ntr.creationDate  FROM `%s.%s` ntr INNER JOIN `%s.%s` nnv ON ntr.nodeUid = nnv.uid  "+
		"WHERE nodeUid = @nodeUid ORDER BY relativeLevel", WoptaDataset,
		NetworkTreeStructureTable, WoptaDataset, NetworkNodesView)
	params := map[string]interface{}{
		"nodeUid": nn.Uid,
	}

	ancestors, err := lib.QueryParametrizedRowsBigQuery[NetworkTreeElement](query, params)
	if err != nil {
		log.Printf("error fetching ancestors from BigQuery for node %s: %s", nn.Uid, err.Error())
		return nil, err
	}

	return ancestors, nil
}

func (nn *NetworkNode) GetChildren() ([]NetworkTreeElement, error) {
	var (
		err error
	)

	query := fmt.Sprintf("SELECT rootUid, ntr.parentUid, nodeUid, nnv.name AS name, relativeLevel, "+
		"ntr.creationDate  FROM `%s.%s` ntr INNER JOIN `%s.%s` nnv ON ntr.nodeUid = nnv.uid  "+
		"WHERE rootUid = @rootUid ORDER BY relativeLevel", WoptaDataset,
		NetworkTreeStructureTable, WoptaDataset, NetworkNodesView)
	params := map[string]interface{}{
		"rootUid": nn.Uid,
	}

	children, err := lib.QueryParametrizedRowsBigQuery[NetworkTreeElement](query, params)
	if err != nil {
		log.Printf("error fetching children from BigQuery for node %s: %s", nn.Uid, err.Error())
		return nil, err
	}

	return children, nil
}

func (nn *NetworkNode) IsParentOf(nodeUid string) bool {
	children, _ := nn.GetChildren()
	return len(lib.SliceFilter(children, func(child NetworkTreeElement) bool {
		return child.NodeUid == nodeUid
	})) == 1
}
