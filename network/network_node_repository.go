package network

import (
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetNodeByUid(uid string) (models.NetworkNode, error) {
	var node *models.NetworkNode
	docSnapshot := lib.WhereLimitFirestore(models.NetworkNodesCollection, "uid", "==", uid, 1)

	documents, err := docSnapshot.GetAll()

	if err != nil {
		return models.NetworkNode{}, fmt.Errorf("could not fetch node: %s", err.Error())
	}

	if len(documents) == 0 {
		return models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	err = documents[0].DataTo(node)

	if node == nil || err != nil {
		return models.NetworkNode{}, fmt.Errorf("could not parse node: %s", err.Error())
	}
	return *node, err
}

func CreateNode(node models.NetworkNode) (string, error) {
	if len(node.Uid) == 0 {
		node.Uid = lib.NewDoc(models.NetworkNodesCollection)
	}
	return node.Uid, lib.SetFirestoreErr(models.NetworkNodesCollection, node.Uid, node)
}

func GetNodeByUidBigQuery(uid string) (models.NetworkNode, error) {
	query := "select * from `%s.%s` where uid = @uid limit 1"
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesTable)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes[0], err
}

func CreateNodeBigQuery(node models.NetworkNode) error {
	return lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkNodesTable, node)
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
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesTable, models.WoptaDataset, models.NetworkNodesTable)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}

func GetAllParentNodesFromNode(uid string) ([]models.NetworkNode, error) {
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
	query = fmt.Sprintf(query, models.WoptaDataset, models.NetworkNodesTable, models.WoptaDataset, models.NetworkNodesTable)
	params := map[string]interface{}{"uid": uid}
	nodes, err := lib.QueryParametrizedRowsBigQuery[models.NetworkNode](query, params)

	if len(nodes) == 0 {
		return []models.NetworkNode{}, fmt.Errorf("could not find node with uid %s", uid)
	}
	return nodes, err
}
