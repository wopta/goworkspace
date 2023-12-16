package _script

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"os"
)

func UpdateNetworkNodeInfoSimplo() {
	var (
		networkNodes []models.NetworkNode
	)

	ctx := context.Background()
	client, _ := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	iter := client.Collection(models.NetworkNodesCollection).Documents(ctx)
	networkNodes = models.NetworkNodeToListData(iter)

	for _, nn := range networkNodes {
		if nn.Type == models.AgencyNetworkNodeType || nn.Type == models.AgentNetworkNodeType {
			nn.HasAnnex = true
			nn.IsMgaProponent = true
			nn.Designation = "Addetto Attività intermediazione al di fuori dei locali"
			nn.WorksForUid = models.WorksForMgaUid

			err := lib.SetFirestoreErr(models.NetworkNodesCollection, nn.Uid, nn)
			if err != nil {
				log.Printf("error updating network node %s: %s", nn.Code, err.Error())
				continue
			}

			err = nn.SaveBigQuery("")
			if err != nil {
				log.Printf("error updating network node %s in BigQuery: %s", nn.Code, err.Error())
				continue
			}

			log.Printf("Network Node %s Updated Succesfully", nn.Code)
		} else {
			log.Printf("Network Node %s Skipped", nn.Code)
		}
	}
}
