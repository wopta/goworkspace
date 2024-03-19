package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

func NodeSubTreeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err  error
		root models.NetworkTreeElement
	)

	log.SetPrefix("NodeSubTreeFx ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	nodeUid := r.Header.Get("nodeUid")
	log.Printf("Node Uid: %s", nodeUid)

	log.Println("loading authToken from idToken...")

	idToken := r.Header.Get("Authorization")
	err = checkAccess(idToken, nodeUid)
	if err != nil {
		return "", nil, err
	}

	root, err = GetNodeSubTree(nodeUid)

	rawRoot, err := json.Marshal(root)

	log.Println("Handler end -------------------------------------------------")

	return string(rawRoot), root, err
}

func GetNodeSubTree(nodeUid string) (models.NetworkTreeElement, error) {
	var (
		err  error
		root models.NetworkTreeElement
	)

	log.Printf("Fetching node from Firestore...")
	node := GetNetworkNodeByUid(nodeUid)
	if node == nil {
		log.Printf("no node found with uid %s", nodeUid)
		return root, errors.New("node not found")
	}
	log.Printf("Node found")

	log.Printf("Fetching children for node %s", nodeUid)

	query := fmt.Sprintf("SELECT rootUid, ntr.parentUid, nodeUid, nnv.name AS name, relativeLevel, ntr.creationDate "+
		"FROM `%s.%s` AS ntr "+
		"INNER JOIN `%s.%s` nnv ON ntr.nodeUid = nnv.uid "+
		"WHERE ntr.rootUid = @rootUid AND ntr.nodeUid != @rootUid",
		models.WoptaDataset, models.NetworkTreeStructureTable, models.WoptaDataset, models.NetworkNodesView)

	params := map[string]interface{}{
		"rootUid": nodeUid,
	}

	subNetwork, err := lib.QueryParametrizedRowsBigQuery[models.NetworkTreeElement](query, params)
	if err != nil {
		log.Printf("error fetching children from BigQuery for node %s: %s", nodeUid, err.Error())
		return root, err
	}

	log.Printf("Found %02d children for node %s", len(subNetwork), nodeUid)

	root = models.NetworkTreeElement{
		RootUid:   node.Uid,
		NodeUid:   node.Uid,
		ParentUid: node.ParentUid,
		Name:      node.GetName(),
	}

	root = getNodeChildren(root, subNetwork)

	return root, err
}

func checkAccess(idToken, nodeUid string) error {
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error getting authToken")
		return err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)
	if authToken.Role != models.UserRoleAdmin && authToken.UserID != nodeUid {
		query := fmt.Sprintf("SELECT rootUid, parentUid, nodeUid, relativeLevel, creationDate "+
			"FROM `%s.%s`"+
			"WHERE rootUid = @rootUid AND nodeUid = @nodeUid", models.WoptaDataset, models.NetworkTreeStructureTable)

		params := map[string]interface{}{
			"rootUid": authToken.UserID,
			"nodeUid": nodeUid,
		}

		result, err := lib.QueryParametrizedRowsBigQuery[models.NetworkTreeElement](query, params)
		if err != nil {
			log.Printf("error fetching children from BigQuery for node %s: %s", nodeUid, err.Error())
			return err
		}
		if len(result) == 0 {
			log.Printf("node %s not autorized to access subtree with root uid %s", authToken.UserID, nodeUid)
			return errors.New("cannot access subtree")
		}
	}
	return nil
}

func getNodeChildren(node models.NetworkTreeElement, allNodes []models.NetworkTreeElement) models.NetworkTreeElement {
	children := lib.SliceFilter(allNodes, func(structure models.NetworkTreeElement) bool {
		return structure.ParentUid == node.NodeUid
	})
	if len(children) == 0 {
		return node
	}

	node.Children = make([]models.NetworkTreeElement, 0)
	for _, child := range children {
		res := getNodeChildren(child, allNodes)
		node.Children = append(node.Children, res)
	}
	return node
}
