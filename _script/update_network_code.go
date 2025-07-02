package _script

import (
	"log"

	"gitlab.dev.wopta.it/goworkspace/network"
)

func UpdateAgentRuiCode(nodeUid, ruiCode string) {
	networkNode, err := network.GetNodeByUidErr(nodeUid)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

	log.Printf("old ruiCode: %s", networkNode.Agent.RuiCode)

	networkNode.Agent.RuiCode = ruiCode

	log.Printf("new ruiCode: %s", networkNode.Agent.RuiCode)

	err = networkNode.SaveFirestore()
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

	err = networkNode.SaveBigQuery("")
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}
}
