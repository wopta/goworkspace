package _script

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func DisableAllNodes() {
	var (
		networkNodes []models.NetworkNode
	)

	docsnap := lib.WhereFirestore(models.NetworkNodesCollection, "authId", "!=", "")
	networkNodes = models.NetworkNodeToListData(docsnap)

	for _, nn := range networkNodes {
		if nn.AuthId != "" {
			log.Printf("NetworkNode Code: %s", nn.Code)
			err := lib.HandleUserAuthenticationStatus(nn.Uid, true)
			if err != nil {
				return
			}
			log.Printf("NetworkNode %s disabled", nn.Code)
		}
	}
}

func EnableAllNodes() {
	var (
		networkNodes []models.NetworkNode
	)

	docsnap := lib.WhereFirestore(models.NetworkNodesCollection, "authId", "!=", "")
	networkNodes = models.NetworkNodeToListData(docsnap)

	for _, nn := range networkNodes {
		if nn.AuthId != "" && nn.IsActive {
			log.Printf("NetworkNode Code: %s", nn.Code)
			err := lib.HandleUserAuthenticationStatus(nn.Uid, false)
			if err != nil {
				return
			}
			log.Printf("NetworkNode %s enabled", nn.Code)
		}
	}
}
