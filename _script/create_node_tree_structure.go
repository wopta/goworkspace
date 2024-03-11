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
	Parents    []nodeInfo
}

type bigQueryNodeInfo struct {
	ChildUid      string                `bigquery:"childUid"`
	AbsoluteLevel int                   `bigquery:"absoluteLevel"`
	RelativeLevel int                   `bigquery:"relativeLevel"`
	ParentUid     string                `bigquery:"parentUid"`
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
			Uid:     node.Uid,
			Code:    node.Code,
			Level:   1,
			Parents: []nodeInfo{},
		})
	}

	for len(toBeVisitedNodes) > 0 {
		index := len(toBeVisitedNodes) - 1
		currentNode := toBeVisitedNodes[index]
		children := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
			return node.ParentUid == currentNode.Uid
		})
		toBeVisitedNodes = toBeVisitedNodes[:index]
		parents := append(currentNode.Parents, currentNode)
		for _, child := range children {
			toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
				Uid:     child.Uid,
				Code:    child.Code,
				Level:   currentNode.Level + 1,
				Parents: parents,
			})
		}
		visitedNodes = append(visitedNodes, currentNode)

	}

	for _, nn := range visitedNodes {
		if len(nn.Parents) > 0 {
			dbNode := bigQueryNodeInfo{
				ChildUid:      nn.Uid,
				AbsoluteLevel: nn.Level,
				ParentUid:     "",
				CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			}

			for _, p := range nn.Parents {
				dbNode.ParentUid = p.Uid
				dbNode.RelativeLevel = nn.Level - p.Level

				log.Printf("child: %s\tparent: %s\t childLevel: %02d\tparentLevel: %02d\trelativeLevel: %02d\n", nn.Code, p.Code,
					dbNode.AbsoluteLevel, p.Level, dbNode.RelativeLevel)

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
