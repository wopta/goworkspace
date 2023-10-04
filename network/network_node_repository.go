package network

import (
	"fmt"
	"log"
	"time"

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
	node.IsActive = true
}

func GetNetworkNodeByUid(producerUid string) *models.NetworkNode {
	if producerUid == "" {
		log.Println("[GetNetworkNodeByPolicy] producerUid empty")
		return nil
	}

	networkNode, err := GetNodeByUid(producerUid)
	if err != nil {
		log.Printf("[GetNetworkNodeByPolicy] error getting producer %s from Firestore", producerUid)
	}

	return networkNode
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
