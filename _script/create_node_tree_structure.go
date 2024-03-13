package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"log"
	"time"
)

type nodeInfo struct {
	models.NetworkTreeElement
	Ancestors []nodeInfo
}

func CreateNodeTreeStructure() {
	log.Println("function start ----------------------------------------------")

	creationDate := lib.GetBigQueryNullDateTime(time.Now().UTC())

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
			NetworkTreeElement: models.NetworkTreeElement{
				RootUid:       node.Uid,
				ParentUid:     "",
				NodeUid:       node.Uid,
				Name:          node.GetName(),
				AbsoluteLevel: 1,
				RelativeLevel: 0,
				CreationDate:  creationDate,
			},
			Ancestors: []nodeInfo{},
		})
	}

	for len(toBeVisitedNodes) > 0 {
		index := len(toBeVisitedNodes) - 1
		currentNode := toBeVisitedNodes[index]
		children := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
			return node.ParentUid == currentNode.NodeUid
		})
		toBeVisitedNodes = toBeVisitedNodes[:index]
		parents := append(currentNode.Ancestors, currentNode)
		for _, child := range children {
			toBeVisitedNodes = append(toBeVisitedNodes, nodeInfo{
				NetworkTreeElement: models.NetworkTreeElement{
					RootUid:       currentNode.RootUid,
					ParentUid:     currentNode.NodeUid,
					NodeUid:       child.Uid,
					Name:          child.GetName(),
					AbsoluteLevel: currentNode.AbsoluteLevel + 1,
					CreationDate:  creationDate,
				},
				Ancestors: parents,
			})
		}
		visitedNodes = append(visitedNodes, currentNode)

	}

	for _, nn := range visitedNodes {
		err = writeNodeToBigQuery(nn)
		if err != nil {
			log.Printf("error writing node %s to BigQuery: %s", nn.NodeUid, err.Error())
			return
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

func writeNodeToBigQuery(node nodeInfo) error {
	if len(node.Ancestors) > 0 {
		dbNode := models.NetworkTreeElement{
			NodeUid:       node.NodeUid,
			AbsoluteLevel: node.AbsoluteLevel,
			ParentUid:     node.ParentUid,
			Name:          node.Name,
			CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
		}

		for _, p := range node.Ancestors {
			dbNode.RootUid = p.NodeUid
			dbNode.RelativeLevel = node.AbsoluteLevel - p.AbsoluteLevel

			log.Printf("rootUid: %s\tparentUid: %s\tnodeUid: %s\tchildLevel: %02d\tparentLevel: %02d\trelativeLevel: %02d\t\n",
				dbNode.RootUid, dbNode.ParentUid, node.NodeUid, dbNode.AbsoluteLevel, p.AbsoluteLevel, dbNode.RelativeLevel)

			err := lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, dbNode)
			if err != nil {
				log.Printf("insert error: %s", err.Error())
				return err
			}
		}
	}
	return nil
}
