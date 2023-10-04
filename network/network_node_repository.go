package network

import (
	"fmt"
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetNodeByUid(uid string) (models.NetworkNode, error) {
	var node *models.NetworkNode
	docSnapshot, err := lib.GetFirestoreErr(models.NetworkNodesCollection, uid)

	if err != nil {
		return models.NetworkNode{}, fmt.Errorf("could not fetch node: %s", err.Error())
	}
	err = docSnapshot.DataTo(&node)

	if node == nil || err != nil {
		return models.NetworkNode{}, fmt.Errorf("could not parse node: %s", err.Error())
	}
	return *node, err
}

func initNode(node *models.NetworkNode) {
	if len(node.Uid) == 0 {
		node.Uid = lib.NewDoc(models.NetworkNodesCollection)
	}
	now := time.Now().UTC()
	node.CreationDate, node.UpdatedDate = now, now
	node.NetworkUid = node.NetworkCode
	node.IsActive = true
}

func CreateNode(node models.NetworkNode) (string, error) {
	initNode(&node)
	return node.Uid, lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
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

func UpdateNetworkNodePortfolio(origin string, policy *models.Policy, networkNode *models.NetworkNode) error {
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
