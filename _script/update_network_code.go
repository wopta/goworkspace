package _script

import (
	"github.com/wopta/goworkspace/network"
	"log"
)

func UpdateAgentRuiCode(nodeUid, ruiCode string) {
	networkNode, err := network.GetNodeByUid(nodeUid)
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
