package network

import (
	"log"

	"github.com/wopta/goworkspace/models"
)

func TraverseNetworkByNodeUid(nodeUid string, callback func(n *models.NetworkNode)) {
	log.Printf("[TraverseNetworkByNodeUid] networkNodeUid %s", nodeUid)

	node, err := GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("[TraverseNetworkByNodeUid] error retrieving node from firestore: %s", err.Error())
		return
	}

	log.Printf("[TraverseNetworkByNodeUid] executing callback for node %s", node.Uid)
	callback(&node)

	if node.ManagerUid != "" {
		manager, err := GetNodeByUid(node.ManagerUid)
		if err != nil {
			log.Printf("[TraverseNetworkByNodeUid] error retrieving manager node from firestore: %s", err.Error())
			return
		}
		log.Printf("[TraverseNetworkByNodeUid] executing callback for node %s", manager.Uid)
		callback(&manager)
	}

	if node.ParentUid != "" {
		log.Printf("[TraverseNetworkByNodeUid] recursive call for node %s", node.ParentUid)
		TraverseNetworkByNodeUid(node.ParentUid, callback)
	} else {
		log.Println("[TraverseNetworkByNodeUid] traverse completed")
	}
}
