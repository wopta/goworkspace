package network

import (
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func GetNodeAncestors(nodeUid string) ([]models.NetworkTreeElement, error) {
	var (
		err error
	)

	query := fmt.Sprintf("SELECT rootUid, ntr.parentUid, nodeUid, COALESCE(nnv.name, '') AS name, relativeLevel, "+
		"ntr.creationDate  FROM `%s.%s` ntr INNER JOIN `%s.%s` nnv ON ntr.nodeUid = nnv.uid  "+
		"WHERE nodeUid = @nodeUid ORDER BY relativeLevel", models.WoptaDataset,
		models.NetworkTreeStructureTable, models.WoptaDataset, models.NetworkNodesView)
	params := map[string]interface{}{
		"nodeUid": nodeUid,
	}

	ancestors, err := lib.QueryParametrizedRowsBigQuery[models.NetworkTreeElement](query, params)
	if err != nil {
		log.ErrorF("error fetching ancestors from BigQuery for node %s: %s", nodeUid, err.Error())
		return nil, err
	}

	return ancestors, nil
}

func GetNodeChildren(nodeUid string) ([]models.NetworkTreeElement, error) {
	var (
		err error
	)

	query := fmt.Sprintf("SELECT rootUid, ntr.parentUid, nodeUid, COALESCE(nnv.name, '') AS name, relativeLevel, "+
		"ntr.creationDate  FROM `%s.%s` ntr INNER JOIN `%s.%s` nnv ON ntr.nodeUid = nnv.uid  "+
		"WHERE rootUid = @rootUid ORDER BY relativeLevel", models.WoptaDataset,
		models.NetworkTreeStructureTable, models.WoptaDataset, models.NetworkNodesView)
	params := map[string]interface{}{
		"rootUid": nodeUid,
	}

	children, err := lib.QueryParametrizedRowsBigQuery[models.NetworkTreeElement](query, params)
	if err != nil {
		log.ErrorF("error fetching children from BigQuery for node %s: %s", nodeUid, err.Error())
		return nil, err
	}

	return children, nil
}

func IsParentOf(parentUid, childUid string) bool {
	children, _ := GetNodeChildren(parentUid)
	return len(lib.SliceFilter(children, func(child models.NetworkTreeElement) bool {
		return child.NodeUid == childUid
	})) == 1
}

func IsChildOf(parentUid, childUid string) bool {
	ancestors, _ := GetNodeAncestors(childUid)
	return len(lib.SliceFilter(ancestors, func(ancestor models.NetworkTreeElement) bool {
		return ancestor.NodeUid == parentUid
	})) == 1
}
