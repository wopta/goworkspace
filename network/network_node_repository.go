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

func initNode(node *models.NetworkNode) {
	if len(node.Uid) == 0 {
		node.Uid = lib.NewDoc(models.NetworkNodesCollection)
	}
	now := time.Now().UTC()
	node.CreationDate, node.UpdatedDate = now, now
	node.NetworkUid = node.NetworkCode
	node.Role = node.Type
	node.IsActive = true
}

func CreateNode(node models.NetworkNode) (*models.NetworkNode, error) {
	initNode(&node)
	return &node, lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
}

func UpdateNode(node models.NetworkNode) error {
	var originalNode models.NetworkNode
	//updatedNode := make(map[string]interface{}, 0)

	log.Println("[UpdateNode] function start ----------------------------------")

	log.Printf("[UpdateNode] fetching network node %s from Firestore...", node.Uid)

	docSnap, err := lib.GetFirestoreErr(models.NetworkNodesCollection, node.Uid)
	if err != nil {
		log.Printf("[UpdateNode] error fetching network node from firestore: %s", err.Error())
		return err
	}
	err = docSnap.DataTo(&originalNode)
	if err != nil {
		log.Printf("[UpdateNode] error unmarshaling network node %s: %s", node.Uid, err.Error())
		return err
	}

	originalNode.Mail = node.Mail
	originalNode.Warrant = node.Warrant
	originalNode.ParentUid = node.ParentUid
	originalNode.IsActive = node.IsActive
	originalNode.Designation = node.Designation
	originalNode.HasAnnex = node.HasAnnex
	originalNode.UpdatedDate = time.Now().UTC()

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
