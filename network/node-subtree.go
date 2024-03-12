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

type NodeSubTree struct {
	NodeUid       string        `json:"nodeUid" bigquery:"childUid"`
	Name          string        `json:"name" bigquery:"name"`
	ParentUid     string        `json:"parentUid" bigquery:"parentUid"`
	RootUid       string        `json:"rootUid" bigquery:"rootUid"`
	AbsoluteLevel int           `json:"absoluteLevel" bigquery:"absoluteLevel"`
	RelativeLevel int           `json:"relativeLevel" bigquery:"relativeLevel"`
	Children      []NodeSubTree `json:"children"`
}

func NodeSubTreeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err  error
		root NodeSubTree
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

func GetNodeSubTree(nodeUid string) (NodeSubTree, error) {
	var (
		err  error
		root NodeSubTree
	)

	log.Printf("Fetching node from Firestore...")
	node := GetNetworkNodeByUid(nodeUid)
	if node == nil {
		log.Printf("no node found with uid %s", nodeUid)
		return root, errors.New("node not found")
	}
	log.Printf("Node found")

	log.Printf("Fetching children for node %s", nodeUid)

	baseQuery := fmt.Sprintf("SELECT rootUid, parentUid, childUid, absoluteLevel, relativeLevel, name "+
		"FROM `%s.%s` WHERE ", models.WoptaDataset, "node-tree-structure")
	whereClause := fmt.Sprintf("rootUid = '%s'", nodeUid)
	query := fmt.Sprintf("%s %s %s", baseQuery, whereClause, "ORDER BY absoluteLevel")
	subNetwork, err := lib.QueryRowsBigQuery[NodeSubTree](query)
	if err != nil {
		log.Printf("error fetching children from BigQuery for node %s: %s", nodeUid, err.Error())
		return root, err
	}

	root = NodeSubTree{
		NodeUid:   node.Uid,
		Name:      node.GetName(),
		ParentUid: node.ParentUid,
	}

	root = visitNode(root, subNetwork)

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
		baseQuery := fmt.Sprintf("SELECT * FROM `%s.%s` WHERE ", models.WoptaDataset, "node-tree-structure")
		whereClause := fmt.Sprintf("rootUid = '%s' AND childUid = '%s'", authToken.UserID, nodeUid)
		query := fmt.Sprintf("%s %s", baseQuery, whereClause)
		result, err := lib.QueryRowsBigQuery[NodeSubTree](query)
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

func visitNode(node NodeSubTree, allNodes []NodeSubTree) NodeSubTree {
	children := lib.SliceFilter(allNodes, func(structure NodeSubTree) bool {
		return structure.ParentUid == node.NodeUid
	})
	if len(children) == 0 {
		return node
	}

	node.Children = make([]NodeSubTree, 0)
	for _, child := range children {
		res := visitNode(child, allNodes)
		node.Children = append(node.Children, res)
	}
	return node
}
