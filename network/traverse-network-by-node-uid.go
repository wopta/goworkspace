package network

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/models"
)

func TraverseWithCallbackNetworkByNodeUid(
	node *models.NetworkNode,
	lastName string,
	callback func(n *models.NetworkNode, currentName string) string,
) {
	log.AddPrefix("TraverseNetworkByNodeUid")
	defer log.PopPrefix()

	log.Printf("executing callback for node %s", node.Code)
	name := callback(node, lastName)

	if node.ManagerUid != "" {
		manager, err := GetNodeByUidErr(node.ManagerUid)
		if err != nil {
			log.Printf("error retrieving manager node from firestore: %s", err.Error())
			return
		}
		if manager == nil {
			return
		}
		log.Printf("executing callback for node %s", manager.Uid)
		callback(manager, name)
	}

	if node.ParentUid != "" {
		parentNode, err := GetNodeByUidErr(node.ParentUid)
		if err != nil {
			log.ErrorF("error retrieving node from firestore: %s", err.Error())
			return
		}
		if parentNode == nil {
			log.ErrorF("error no node found: %s", node.ParentUid)
			return
		}
		log.Printf("recursive call for node %s", parentNode.Code)

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
		log.Println("traverse completed")
	}
}
