package network

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetNodeByUid(uid string) (*models.NetworkNode, error) {
	var node *models.NetworkNode
	docSnapshot, err := lib.GetFirestoreErr(models.NetworkNodesCollection, uid)

	if err != nil {
		return nil, fmt.Errorf("could not fetch node: %s", err.Error())
	}
	err = docSnapshot.DataTo(&node)

	if node == nil || err != nil {
		return nil, fmt.Errorf("could not parse node: %s", err.Error())
	}
	return node, err
}

func initNode(node *models.NetworkNode) error {
	if len(node.Uid) == 0 {
		node.Uid = lib.NewDoc(models.NetworkNodesCollection)
	}
	now := time.Now().UTC()
	node.CreationDate, node.UpdatedDate = now, now
	node.NetworkUid = node.NetworkCode
	node.Role = node.Type
	node.IsActive = true

	if node.Type != models.PartnershipNetworkNodeType && node.Type != models.AreaManagerNetworkNodeType {
		if node.ExternalNetworkCode == "" {
			node.ExternalNetworkCode = node.Code
		}

		warrant := node.GetWarrant()
		if warrant == nil {
			return fmt.Errorf("could not find warrant for node with value '%s'", node.Warrant)
		}

		if node.Products == nil {
			node.Products = make([]models.Product, 0)
			for _, product := range warrant.Products {
				companies := make([]models.Company, 0)
				for _, company := range product.Companies {
					companies = append(companies, models.Company{
						Name:         company.Name,
						ProducerCode: node.Code,
					})
				}
				node.Products = append(node.Products, models.Product{
					Name:      product.Name,
					Companies: companies,
				})
			}
		} else {
			for prodIndex, product := range node.Products {
				for companyIndex, company := range product.Companies {
					if company.ProducerCode == "" {
						node.Products[prodIndex].Companies[companyIndex].ProducerCode = node.Code
					}
				}
			}
		}
	}

	if node.IsMgaProponent {
		node.HasAnnex = true
	}

	return addWorksForUid(node, node)
}

func CreateNode(node models.NetworkNode) (*models.NetworkNode, error) {
	if err := initNode(&node); err != nil {
		return nil, err
	}
	return &node, lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
}

func UpdateNode(node models.NetworkNode) error {
	var (
		err          error
		originalNode *models.NetworkNode
	)

	log.Println("[UpdateNode] function start ----------------------------------")

	log.Printf("[UpdateNode] fetching network node %s from Firestore...", node.Uid)

	originalNode, err = GetNodeByUid(node.Uid)
	if err != nil {
		log.Printf("[UpdateNode] error fetching network node from firestore: %s", err.Error())
		return err
	}

	if originalNode.AuthId != "" && originalNode.Mail != node.Mail {
		_, err := lib.UpdateUserEmail(node.Uid, node.Mail)
		if err != nil {
			log.Printf("[UpdateNode] error updating network node mail on Firebase Auth: %s", err.Error())
			return err
		}
	}
	originalNode.Mail = node.Mail
	originalNode.Warrant = node.Warrant
	originalNode.Products = node.Products
	originalNode.ParentUid = node.ParentUid
	originalNode.IsActive = node.IsActive
	originalNode.Designation = node.Designation
	originalNode.HasAnnex = node.HasAnnex
	originalNode.UpdatedDate = time.Now().UTC()

	originalNode.IsMgaProponent = node.IsMgaProponent
	originalNode.HasAnnex = node.HasAnnex
	if originalNode.IsMgaProponent {
		originalNode.HasAnnex = true
	}
	err = addWorksForUid(originalNode, &node)
	if err != nil {
		log.Printf("[UpdateNode] error updating WorksForUid '%s' in network node '%s': %s", originalNode.WorksForUid, originalNode.Uid, err.Error())
		return err
	}

	switch node.Type {
	case models.AgentNetworkNodeType:
		originalNode.Agent = node.Agent
	case models.AgencyNetworkNodeType:
		originalNode.Agency = node.Agency
	case models.BrokerNetworkNodeType:
		originalNode.Broker = node.Broker
	case models.AreaManagerNetworkNodeType:
		originalNode.AreaManager = node.AreaManager
	}

	if originalNode.AuthId == "" {
		originalNode.Code = node.Code
		originalNode.Type = node.Type
		originalNode.Role = node.Type
	}

	log.Printf("[UpdateNode] writing network node %s in Firestore...", originalNode.Uid)

	err = lib.SetFirestoreErr(models.NetworkNodesCollection, originalNode.Uid, originalNode)
	if err != nil {
		log.Printf("[UpdateNode] error updating network node %s in Firestore", originalNode.Uid)
		return err
	}

	log.Printf("[UpdateNode] writing network node %s in BigQuery...", originalNode.Uid)

	return originalNode.SaveBigQuery("")
}

func GetNetworkNodeByUid(nodeUid string) *models.NetworkNode {
	if nodeUid == "" {
		log.Println("[GetNetworkNodeByUid] nodeUid empty")
		return nil
	}

	networkNode, err := GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("[GetNetworkNodeByUid] error getting producer %s from Firestore", nodeUid)
	}

	return networkNode
}

func GetAllNetworkNodes() ([]models.NetworkNode, error) {
	var nodes []models.NetworkNode
	docIterator := lib.OrderFirestore(models.NetworkNodesCollection, "code", firestore.Asc)

	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Printf("[GetAllNetworkNodes] error getting nodes from Firestore: %s", err.Error())
		return nodes, err
	}

	for _, snapshot := range snapshots {
		var node models.NetworkNode
		err = snapshot.DataTo(&node)
		if err != nil {
			log.Printf("[GetAllNetworkNodes] error parsing node %s: %s", snapshot.Ref.ID, err.Error())
		} else {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func DeleteNetworkNodeByUid(origin, nodeUid string) error {
	if nodeUid == "" {
		log.Println("[DeleteNetworkNodeByUid] no nodeUid specified")
		return fmt.Errorf("no nodeUid specified")
	}

	fireNetwork := lib.GetDatasetByEnv(origin, models.NetworkNodesCollection)
	_, err := lib.DeleteFirestoreErr(fireNetwork, nodeUid)
	return err
}

func UpdateNetworkNodePortfolio(origin string, policy *models.Policy, networkNode *models.NetworkNode) error {
	if networkNode == nil {
		log.Printf("[UpdateNetworkNodePortfolio] no networkNode specified")
		return nil
	}

	log.Printf("[UpdateNetworkNodePortfolio] adding policy %s to networkNode %s portfolio", policy.Uid, networkNode.Uid)

	networkNode.Policies = append(networkNode.Policies, policy.Uid)

	if !lib.SliceContains(networkNode.Users, policy.Contractor.Uid) {
		log.Printf("[UpdateNetworkNodePortfolio] adding user %s to networkNode %s users list", policy.Contractor.Uid, networkNode.Uid)
		networkNode.Users = append(networkNode.Users, policy.Contractor.Uid)
	}

	networkNode.UpdatedDate = time.Now().UTC()

	log.Printf("[UpdateNetworkNodePortfolio] saving networkNode %s to Firestore...", networkNode.Uid)
	fireNetwork := lib.GetDatasetByEnv(origin, models.NetworkNodesCollection)
	err := lib.SetFirestoreErr(fireNetwork, networkNode.Uid, networkNode)
	if err != nil {
		log.Printf("[UpdateNetworkNodePortfolio] error saving networkNode %s to Firestore: %s", networkNode.Uid, err.Error())
		return err
	}

	log.Printf("[UpdateNetworkNodePortfolio] saving networkNode %s to BigQuery...", networkNode.Uid)
	err = networkNode.SaveBigQuery(origin)

	return err
}

func GetNodeByUidBigQuery(uid string) (models.NetworkNode, error) {
	query := "select * from `%s.%s` where uid = @uid limit 1"
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes[0], err
}

func CreateNodeBigQuery(node models.NetworkNode) error {
	initNode(&node)
	return lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkNodesCollection, node)
}

func GetAllSubNodesFromNodeBigQuery(uid string) ([]models.NetworkNode, error) {
	query := `WITH
	RECURSIVE network AS (
	SELECT
	  *
	FROM
	  ` + "`%s.%s`" + `
	WHERE
	  uid = @uid
	UNION ALL
	SELECT
	  child.*
	FROM
	  ` + "`%s.%s`" + ` child
	JOIN
	  network n
	ON
	  n.uid = child.parentUid )
  SELECT
	*
  FROM
	network n
  WHERE
	uid <> @uid`
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesCollection, models.WoptaDataset, models.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}

func GetAllParentNodesFromNodeBigQuery(uid string) ([]models.NetworkNode, error) {
	query := `WITH
	RECURSIVE network AS (
	SELECT
	  *
	FROM
	  ` + "`%s.%s`" + `
	WHERE
	  uid = @uid
	UNION ALL
	SELECT
	  child.*
	FROM
	  ` + "`%s.%s`" + ` child
	JOIN
	  network n
	ON
	  n.parentUid = child.uid )
  SELECT
	*
  FROM
	network n
  WHERE
	uid <> @uid`
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesCollection, models.WoptaDataset, models.NetworkNodesCollection)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}

func addWorksForUid(originalNode, inputNode *models.NetworkNode) error {
	if inputNode.Type == models.AgentNetworkNodeType && inputNode.WorksForUid != "" && inputNode.WorksForUid != models.WorksForMgaUid {
		worksForNode := GetNetworkNodeByUid(inputNode.WorksForUid)
		// TODO: Check also for broker?
		if worksForNode.Type != models.AgencyNetworkNodeType {
			return fmt.Errorf("worksForUid must reference an agency, got %s", worksForNode.Type)
		}
		if inputNode.IsMgaProponent && !lib.SliceContains(models.GetProponentRuiSections(), worksForNode.Agency.RuiSection) {
			return fmt.Errorf(
				"worksForUid must reference an agency of RuiSection %v when IsMgaProponent is %t",
				models.GetProponentRuiSections(),
				inputNode.IsMgaProponent,
			)
		}
		if !inputNode.IsMgaProponent && !lib.SliceContains(models.GetIssuerRuiSections(), worksForNode.Agency.RuiSection) {
			return fmt.Errorf(
				"worksForUid must reference an agency of RuiSection %v when IsMgaProponent is %t",
				models.GetIssuerRuiSections(),
				inputNode.IsMgaProponent,
			)
		}
	}

	originalNode.WorksForUid = inputNode.WorksForUid
	return nil
}
