package models

import (
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"google.golang.org/api/iterator"
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
	CallbackConfig      *CallbackConfig       `json:"callbackConfig,omitempty" firestore:"callbackConfig,omitempty" bigquery:"-"`
	JwtConfig           lib.JwtConfig         `json:"jwtConfig,omitempty" firestore:"jwtConfig,omitempty" bigquery:"-"`
	Consens             []NodeConsens         `json:"consens" firestore:"consens" bigquery:"-"`
}

type NodeProduct struct {
	Name      string        `json:"-" firestore:"-" bigquery:"name"`
	Companies []NodeCompany `json:"-" firestore:"-" bigquery:"companies"`
}

type NodeCompany struct {
	Name         string `json:"-" firestore:"-" bigquery:"name"`
	ProducerCode string `json:"-" firestore:"-" bigquery:"producerCode"`
}

type CallbackConfig struct {
	Name string `json:"name" firestore:"name" bigquery:"-"`
}

type NodeConsens struct {
	Slug     string            `json:"slug" firestore:"slug" bigquery:"-"`
	ExpireAt time.Time         `json:"expireAt" firestore:"expireAt" bigquery:"-"`
	StartAt  time.Time         `json:"startAt" firestore:"startAt" bigquery:"-"`
	Title    string            `json:"title" firestore:"title" bigquery:"-"`
	Subtitle string            `json:"subtitle" firestore:"subtitle" bigquery:"-"`
	Content  string            `json:"content" firestore:"content" bigquery:"-"`
	Answers  map[string]string `json:"answers" firestore:"answers" bigquery:"-"`
	GivenAt  time.Time         `json:"givenAt" firestore:"givenAt" bigquery:"-"`
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
	log.AddPrefix("NetworkNode.SaveFirestore")
	defer log.PopPrefix()

	if err != nil {
		log.Error(err)
	}
	return err
}

func (nn *NetworkNode) SaveBigQuery(origin string) error {
	log.AddPrefix("NetworkNode.SaveBigQuery")
	defer log.PopPrefix()

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
	log.AddPrefix("GetWarrant")
	defer log.PopPrefix()
	if nn.Warrant == "" {
		log.Printf("warrant not set for node %s", nn.Uid)
		return nil
	}

	log.Printf("requesting warrant %s", nn.Warrant)

	warrantBytes := lib.GetFilesByEnv(fmt.Sprintf(WarrantFormat, nn.Warrant))

	err := json.Unmarshal(warrantBytes, &warrant)
	if err != nil {
		log.ErrorF("error unmarshaling warrant %s: %s", nn.Warrant, err.Error())
		return nil
	}

	return warrant
}

func (nn *NetworkNode) HasAccessToProduct(productName string, warrant *Warrant) bool {
	log.AddPrefix("HasAccessToProduct")
	defer log.PopPrefix()
	log.Println("method start -----------------")

	needCheckTypes := []string{AgencyNetworkNodeType, AgentNetworkNodeType, BrokerNetworkNodeType}

	if !lib.SliceContains(needCheckTypes, nn.Type) {
		return true
	}

	if warrant == nil {
		warrant = nn.GetWarrant()
	}
	if warrant == nil {
		log.ErrorF("no %s warrant found", nn.Warrant)
		return false
	}

	log.Printf("checking if network node %s has access product %s", nn.Uid, productName)

	for _, product := range warrant.Products {
		if product.Name == productName {
			return true
		}
	}

	return false
}

func (nn *NetworkNode) GetNetworkNodeFlow(productName string, warrant *Warrant) (string, []byte) {
	if warrant == nil {
		log.ErrorF("error warrant not set for node %s", nn.Uid)
		return "", []byte{}
	}

	product := warrant.GetProduct(productName)
	if product == nil {
		log.ErrorF("error product not set for warrant %s", warrant.Name)
		return "", []byte{}
	}

	log.Printf("getting flow '%s' file for product '%s'", product.Flow, productName)

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

func (nn *NetworkNode) GetRuiCode() string {
	var ruiCode string

	switch nn.Type {
	case AgentNetworkNodeType:
		ruiCode = nn.Agent.RuiCode
	case AgencyNetworkNodeType:
		ruiCode = nn.Agency.RuiCode
	case BrokerNetworkNodeType:
		ruiCode = nn.Broker.RuiCode
	case AreaManagerNetworkNodeType:
		ruiCode = nn.AreaManager.RuiCode
	}

	return ruiCode
}

func (nn *NetworkNode) GetRuiSection() string {
	var ruiSection string

	switch nn.Type {
	case AgentNetworkNodeType:
		ruiSection = nn.Agent.RuiSection
	case AgencyNetworkNodeType:
		ruiSection = nn.Agency.RuiSection
	case BrokerNetworkNodeType:
		ruiSection = nn.Broker.RuiSection
	case AreaManagerNetworkNodeType:
		ruiSection = nn.AreaManager.RuiSection
	}

	return ruiSection
}

func (nn *NetworkNode) GetRuiRegistration() time.Time {
	var ruiRegistration time.Time

	switch nn.Type {
	case AgentNetworkNodeType:
		ruiRegistration = nn.Agent.RuiRegistration
	case AgencyNetworkNodeType:
		ruiRegistration = nn.Agency.RuiRegistration
	case BrokerNetworkNodeType:
		ruiRegistration = nn.Broker.RuiRegistration
	case AreaManagerNetworkNodeType:
		ruiRegistration = nn.AreaManager.RuiRegistration
	}

	return ruiRegistration
}

func (nn *NetworkNode) GetVatCode() string {
	var vatCode string

	switch nn.Type {
	case AgentNetworkNodeType:
		vatCode = nn.Agent.VatCode
	case AgencyNetworkNodeType:
		vatCode = nn.Agency.VatCode
	case BrokerNetworkNodeType:
		vatCode = nn.Broker.VatCode
	case AreaManagerNetworkNodeType:
		vatCode = nn.AreaManager.VatCode
	}

	return vatCode
}

func (nn *NetworkNode) GetFiscalCode() string {
	var fiscalCode string

	switch nn.Type {
	case AgentNetworkNodeType:
		fiscalCode = nn.Agent.FiscalCode
	case AreaManagerNetworkNodeType:
		fiscalCode = nn.AreaManager.FiscalCode
	}

	return fiscalCode
}

func (nn *NetworkNode) GetManagerName() string {
	var name string

	switch nn.Type {
	case AgencyNetworkNodeType:
		name = fmt.Sprintf("%s %s", lib.Capitalize(nn.Agency.Manager.Name), lib.Capitalize(nn.Agency.Manager.Surname))
	case BrokerNetworkNodeType:
		name = fmt.Sprintf("%s %s", lib.Capitalize(nn.Broker.Manager.Name), lib.Capitalize(nn.Broker.Manager.Surname))
	}

	return name
}

func (nn *NetworkNode) GetManagerFiscalCode() string {
	var fiscalCode string

	switch nn.Type {
	case AgencyNetworkNodeType:
		fiscalCode = nn.Agency.Manager.FiscalCode
	case BrokerNetworkNodeType:
		fiscalCode = nn.Broker.Manager.FiscalCode
	}

	return fiscalCode
}

func (nn *NetworkNode) GetManagerRuiSection() string {
	var ruiSection string

	switch nn.Type {
	case AgencyNetworkNodeType:
		ruiSection = nn.Agency.Manager.RuiSection
	case BrokerNetworkNodeType:
		ruiSection = nn.Broker.Manager.RuiSection
	}

	return ruiSection
}

func (nn *NetworkNode) GetManagerRuiCode() string {
	var ruiCode string

	switch nn.Type {
	case AgencyNetworkNodeType:
		ruiCode = nn.Agency.Manager.RuiCode
	case BrokerNetworkNodeType:
		ruiCode = nn.Broker.Manager.RuiCode
	}

	return ruiCode
}

func (nn *NetworkNode) GetManagerRuiResgistration() time.Time {
	var ruiRegistration time.Time

	switch nn.Type {
	case AgencyNetworkNodeType:
		ruiRegistration = nn.Agency.Manager.RuiRegistration
	case BrokerNetworkNodeType:
		ruiRegistration = nn.Broker.Manager.RuiRegistration
	}

	return ruiRegistration
}

func (nn *NetworkNode) GetManagerPhone() string {
	var phone string

	switch nn.Type {
	case AgencyNetworkNodeType:
		phone = nn.Agency.Manager.Phone
	case BrokerNetworkNodeType:
		phone = nn.Broker.Manager.Phone
	}

	return phone
}

func (nn *NetworkNode) GetWebsite() string {
	var website string

	switch nn.Type {
	case AgencyNetworkNodeType:
		website = nn.Agency.Website
	case BrokerNetworkNodeType:
		website = nn.Broker.Website
	}

	return website
}

func (nn *NetworkNode) GetPhone() string {
	var phone string

	switch nn.Type {
	case AgentNetworkNodeType:
		phone = nn.Agent.Phone
	case AgencyNetworkNodeType:
		phone = nn.Agency.Phone
	case BrokerNetworkNodeType:
		phone = nn.Broker.Phone
	case AreaManagerNetworkNodeType:
		phone = nn.AreaManager.Phone
	}

	return phone
}

func (nn *NetworkNode) IsJwtProtected() bool {
	c := nn.JwtConfig
	return (c.KeyAlgorithm != "" && c.ContentEncryption != "") || c.SignatureAlgorithm != ""
}

func (nn *NetworkNode) DecryptJwt(jwtData string) ([]byte, error) {
	if !nn.IsJwtProtected() {
		return nil, nil
	}

	return lib.ParseJwt(jwtData, nn.JwtConfig)
}

func (nn NetworkNode) DecryptJwtClaims(jwtData string, unmarshaler func([]byte) (LifeClaims, error)) (LifeClaims, error) {
	bytes, err := nn.DecryptJwt(jwtData)
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
