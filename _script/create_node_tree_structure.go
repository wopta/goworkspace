package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"log"
)

func CreateNodeTreeStructure() {

	log.Println("function start ----------------------------------------------")

	allNodes, err := getAllNodes()
	if err != nil {
		panic(err)
	}

	log.Println(len(allNodes))

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
		return node.Type != models.PartnershipNetworkNodeType && node.Type != models.AreaManagerNetworkNodeType
	})
	log.Printf("found %02d nodes", len(filteredNodes))
	return filteredNodes, nil
}
