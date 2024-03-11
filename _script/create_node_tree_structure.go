package _script

import (
	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"log"
	"time"
)

type nodeInfo struct {
	Uid        string
	Code       string
	Level      int
	BreadCrumb string
	RootUid    string
	ParentUid  string
	Name       string
	Ancestors  []nodeInfo
}

type bigQueryNodeInfo struct {
	RootUid       string                `bigquery:"rootUid"`
	ParentUid     string                `bigquery:"parentUid"`
	ChildUid      string                `bigquery:"childUid"`
	Name          string                `bigquery:"name"`
	AbsoluteLevel int                   `bigquery:"absoluteLevel"`
	RelativeLevel int                   `bigquery:"relativeLevel"`
	CreationDate  bigquery.NullDateTime `bigquery:"creationDate"`
}

func CreateNodeTreeStructure() {
	log.Println("function start ----------------------------------------------")

	nodesList, err := getAllNodes()
	if err != nil {
		panic(err)
	}

	nodesList = lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.Type != models.PartnershipNetworkNodeType
	})

	dbNodes := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.ParentUid == ""
	})

	firstLevelNodes := lib.SliceMap(dbNodes, func(node models.NetworkNode) models.NetworkNode {
		return node
	})

	toBeVisitedNodes := make([]nodeInfo, 0)
	visitedNodes := make([]nodeInfo, 0)

	for _, node := range firstLevelNodes {
		toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
			Uid:       node.Uid,
			Code:      node.Code,
			Level:     1,
			RootUid:   node.Uid,
			ParentUid: "",
			Name:      node.GetName(),
			Ancestors: []nodeInfo{},
		})
	}

	for len(toBeVisitedNodes) > 0 {
		index := len(toBeVisitedNodes) - 1
		currentNode := toBeVisitedNodes[index]
		children := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
			return node.ParentUid == currentNode.Uid
		})
		toBeVisitedNodes = toBeVisitedNodes[:index]
		parents := append(currentNode.Ancestors, currentNode)
		for _, child := range children {
			toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
				Uid:       child.Uid,
				Code:      child.Code,
				Level:     currentNode.Level + 1,
				RootUid:   currentNode.RootUid,
				ParentUid: currentNode.Uid,
				Name:      child.GetName(),
				Ancestors: parents,
			})
		}
		visitedNodes = append(visitedNodes, currentNode)

	}

	for _, nn := range visitedNodes {
		if len(nn.Ancestors) > 0 {
			dbNode := bigQueryNodeInfo{
				ChildUid:      nn.Uid,
				AbsoluteLevel: nn.Level,
				RootUid:       "",
				ParentUid:     nn.ParentUid,
				Name:          nn.Name,
				CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			}

			for _, p := range nn.Ancestors {
				dbNode.RootUid = p.Uid
				dbNode.RelativeLevel = nn.Level - p.Level

				log.Printf("rootUid: %s\tparentUid: %s\tchildUid: %s\tchildLevel: %02d\tparentLevel: %02d\trelativeLevel: %02d\t\n",
					dbNode.RootUid, dbNode.ParentUid, nn.Uid, dbNode.AbsoluteLevel, p.Level, dbNode.RelativeLevel)

				err = lib.InsertRowsBigQuery("wopta", "node-tree-structure", dbNode)
				if err != nil {
					log.Printf("insert error: %s", err.Error())
					return
				}
			}
		}
	}

	log.Println("function end ------------------------------------------------")
}

func getAllNodes() ([]models.NetworkNode, error) {
	log.Println("fetching all nodes from Firestore...")
	allNodes, err := network.GetAllNetworkNodes()
	if err != nil {
		log.Printf("error fetching all nodes from Firestore: %s", err.Error())
		return nil, err
	}
	filteredNodes := lib.SliceFilter(allNodes, func(node models.NetworkNode) bool {
		return node.Type != models.PartnershipNetworkNodeType
	})
	log.Printf("found %02d nodes", len(filteredNodes))
	return filteredNodes, nil
}
