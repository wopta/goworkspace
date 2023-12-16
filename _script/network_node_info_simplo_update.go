package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
)

func UpdateNetworkNodeInfoSimplo() {
	var (
		networkNodes []models.NetworkNode
	)

	docsnap := lib.WhereFirestore(models.NetworkNodesCollection, "authId", "!=", "")
	networkNodes = models.NetworkNodeToListData(docsnap)

	for _, nn := range networkNodes {
		nn.HasAnnex = true
		nn.IsMgaProponent = true
		nn.Designation = "Addetto Attività intermediazione al di fuori dei locali"
		nn.WorksForUid = "__wopta__"

		if nn.Type == models.AgencyNetworkNodeType || nn.Type == models.AgentNetworkNodeType {
			err := lib.SetFirestoreErr(models.NetworkNodesCollection, nn.Uid, nn)
			if err != nil {
				log.Printf("error updating network node %s: %s", nn.Code, err.Error())
				continue
			}

			nn.SaveBigQuery("")
		}
	}
}
