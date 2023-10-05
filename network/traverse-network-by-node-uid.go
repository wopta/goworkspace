package network

import (
	"log"
	"strings"

	"github.com/wopta/goworkspace/models"
)

func TraverseWithCallbackNetworkByNodeUid(
	node *models.NetworkNode,
	lastName string,
	callback func(n *models.NetworkNode, currentName string) string,
) {
	log.Printf("[TraverseNetworkByNodeUid] executing callback for node %s", node.Code)
	name := callback(node, lastName)

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
		parentNode, err := GetNodeByUid(node.ParentUid)
		if err != nil {
			log.Printf("[TraverseNetworkByNodeUid] error retrieving node from firestore: %s", err.Error())
			return
		}
		log.Printf("[TraverseNetworkByNodeUid] recursive call for node %s", parentNode.Code)

		TraverseWithCallbackNetworkByNodeUid(parentNode, name, func(n *models.NetworkNode, currentName string) string {
			var baseName string
			if currentName != "" {
				baseName = currentName
			} else if lastName != "" {
				baseName = strings.ToUpper(strings.Join([]string{lastName, name}, "__"))
			} else {
				baseName = name
			}
			return callback(n, baseName)
		})
	} else {
		log.Println("[TraverseNetworkByNodeUid] traverse completed")
	}
}
