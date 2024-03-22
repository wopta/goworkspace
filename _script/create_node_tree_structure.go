package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"log"
	"time"
)

func CreateNodeTreeStructure() {
	var (
		visitedNodes = make([]models.NetworkTreeElement, 0)
	)

	log.Println("function start ----------------------------------------------")

	nodesList, err := getAllNodes()
	if err != nil {
		panic(err)
	}

	rootNodes := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.ParentUid == ""
	})

	toBeVisitedNodes := lib.SliceMap(rootNodes, func(node models.NetworkNode) models.NetworkTreeElement {
		return models.NetworkTreeElement{
			RootUid:       node.Uid,
			ParentUid:     "",
			NodeUid:       node.Uid,
			AbsoluteLevel: 1,
			RelativeLevel: 0,
			CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			Ancestors:     []models.NetworkTreeElement{},
		}
	})

	for len(toBeVisitedNodes) > 0 {
		index := len(toBeVisitedNodes) - 1
		currentNode := toBeVisitedNodes[index]
		toBeVisitedNodes = toBeVisitedNodes[:index]
		toBeVisitedNodes = append(toBeVisitedNodes, visitNode(currentNode, nodesList)...)
		visitedNodes = append(visitedNodes, currentNode)
	}

	treeElementRelations := make([]models.NetworkTreeRelation, 0)
	for _, nn := range visitedNodes {
		treeElementRelations = append(treeElementRelations, getNodeTreeRelations(nn)...)
	}

	for _, nn := range treeElementRelations {
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

func visitNode(currentNode models.NetworkTreeElement, nodesList []models.NetworkNode) []models.NetworkTreeElement {
	toBeVisitedNodes := make([]models.NetworkTreeElement, 0)

	children := lib.SliceFilter(nodesList, func(node models.NetworkNode) bool {
		return node.ParentUid == currentNode.NodeUid
	})
	if len(children) == 0 {
		return toBeVisitedNodes
	}

	parents := append(currentNode.Ancestors, currentNode)
	for _, child := range children {
		toBeVisitedNodes = append(toBeVisitedNodes, models.NetworkTreeElement{
			RootUid:       currentNode.RootUid,
			ParentUid:     currentNode.NodeUid,
			NodeUid:       child.Uid,
			AbsoluteLevel: currentNode.AbsoluteLevel + 1,
			CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			Ancestors:     parents,
		})
	}
	return toBeVisitedNodes
}

func getNodeTreeRelations(node models.NetworkTreeElement) []models.NetworkTreeRelation {
	ancestors := make([]models.NetworkTreeRelation, 0)
	ancestors = append(ancestors, node.ToNetworkTreeRelation())
	for _, a := range node.Ancestors {
		node.RootUid = a.NodeUid
		node.RelativeLevel = node.AbsoluteLevel - a.AbsoluteLevel
		log.Printf("rootUid: %s\tparentUid: %s\tnodeUid: %s\tchildLevel: %02d\tparentLevel: %02d\trelativeLevel: %02d\t\n",
			node.RootUid, node.ParentUid, node.NodeUid, node.AbsoluteLevel, a.AbsoluteLevel, node.RelativeLevel)
		ancestors = append(ancestors, node.ToNetworkTreeRelation())
	}
	return ancestors
}

func writeNodeToBigQuery(node models.NetworkTreeRelation) error {
	err := lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, node)
	if err != nil {
		log.Printf("insert error: %s", err.Error())
		return err
	}
	return nil
}
