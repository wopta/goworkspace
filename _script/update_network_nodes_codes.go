package _script

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"time"
)

func UpdateNetworkNodesCodes() {
	networkNodes := make([]models.NetworkNode, 0)
	updatedNodes := 0

	ctx := context.Background()
	client, _ := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	iter := client.Collection(models.NetworkNodesCollection).Documents(ctx)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err.Error())
		}
		var networkNode models.NetworkNode
		if err := doc.DataTo(&networkNode); err != nil {
			log.Println(err.Error())
		}
		networkNodes = append(networkNodes, networkNode)
	}

	notUpdatedNodes := 0

	for _, nn := range networkNodes {
		if nn.Type == models.PartnershipNetworkNodeType {
			notUpdatedNodes++
			log.Printf("Partnership Uid: %s", nn.Uid)
			nn.SaveBigQuery("")
			continue
		} else if nn.Type != models.AgencyNetworkNodeType && nn.Type != models.AgentNetworkNodeType {
			log.Printf("Network Node %s Not Updated", nn.Code)
			notUpdatedNodes++
			continue
		}
		if nn.Warrant == "mga_life_agent" {
			nn.Products = []models.Product{
				{
					Name: models.LifeProduct,
					Companies: []models.Company{
						{
							Name:         models.AxaCompany,
							ProducerCode: nn.Code,
						},
					},
				},
			}
		} else if nn.Warrant == "finaip_gap_agency" || nn.Warrant == "mga_life_gap_agent" {
			nn.Products = []models.Product{
				{
					Name: models.GapProduct,
					Companies: []models.Company{
						{
							Name:         models.SogessurCompany,
							ProducerCode: nn.Code,
						},
					},
				},
			}
		} else if nn.Warrant == "mga_multi-product_multi-node" {
			nn.Products = []models.Product{
				{
					Name: models.LifeProduct,
					Companies: []models.Company{
						{
							Name:         models.AxaCompany,
							ProducerCode: nn.Code,
						},
					},
				},
				{
					Name: models.GapProduct,
					Companies: []models.Company{
						{
							Name:         models.SogessurCompany,
							ProducerCode: nn.Code,
						},
					},
				},
			}
		} else {
			log.Printf("Network Node %s Not Updated", nn.Code)
			notUpdatedNodes++
			continue
		}
		nn.ExternalNetworkCode = nn.Code
		nn.UpdatedDate = time.Now().UTC()
		err := lib.SetFirestoreErr(models.NetworkNodesCollection, nn.Uid, nn)
		if err != nil {
			log.Println(err.Error())
		}
		err = nn.SaveBigQuery("")
		if err != nil {
			log.Println(nn.Uid + " error: " + err.Error())
		}
		updatedNodes++
		log.Printf("Network Node %s Updated Succesfully", nn.Code)
	}

	log.Printf("Total Nodes: %d", len(networkNodes))
	log.Printf("Updated Nodes: %d", updatedNodes)
	log.Printf("Not Updated Nodes: %d", notUpdatedNodes)
}
