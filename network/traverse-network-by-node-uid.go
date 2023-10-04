package network

import (
	"log"

	"github.com/wopta/goworkspace/models"
)

func TraverseNetworkByNodeUid(
	node *models.NetworkNode,
	callback func(n *models.NetworkNode, currentName string) string,
) {
	log.Printf("[TraverseNetworkByNodeUid] executing callback for node %s", node.Uid)
	name := callback(node, node.Code)

	if node.ManagerUid != "" {
		manager, err := GetNodeByUid(node.ManagerUid)
		if err != nil {
			log.Printf("[TraverseNetworkByNodeUid] error retrieving manager node from firestore: %s", err.Error())
			return
		}
		log.Printf("[TraverseNetworkByNodeUid] executing callback for node %s", manager.Uid)
		callback(manager, name)
	}

	if node.ParentUid != "" {
		node, err := GetNodeByUid(node.ParentUid)
		if err != nil {
			log.Printf("[TraverseNetworkByNodeUid] error retrieving node from firestore: %s", err.Error())
			return
		}
		log.Printf("[TraverseNetworkByNodeUid] recursive call for node %s", node.Uid)
		TraverseNetworkByNodeUid(node, callback)
	} else {
		log.Println("[TraverseNetworkByNodeUid] traverse completed")
	}
}
