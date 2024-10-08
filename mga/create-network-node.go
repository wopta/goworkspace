package mga

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func CreateNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		inputNode *models.NetworkNode
		err       error
	)

	log.SetPrefix("[CreateNetworkNodeFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", "", err
	}
	// TODO: check node.Type in warrant.AllowedTypes
	// TODO: check unique companyCode for company
	if err := network.TestNetworkNodeUniqueness(inputNode.Code); err != nil {
		log.Printf("error validating node code: %s", err)
		return "", "", err
	}

	log.Println("creating network node into Firestore...")

	inputNode.Normalize()

	node, err := network.CreateNode(*inputNode)
	if err != nil {
		log.Println("error creating network node into Firestore...")
		return "", "", err
	}
	log.Printf("network node created with uid %s", node.Uid)

	node.SaveBigQuery(origin)

	err = createNodeRelation(*node)
	if err != nil {
		log.Printf("error creating node %s network relations: %s", node.Uid, err.Error())
		return "", nil, err
	}

	log.Println("network node successfully created!")
	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
}

func createNodeRelation(node models.NetworkNode) error {
	relations := []models.NetworkTreeRelation{
		{
			RootUid:       node.Uid,
			ParentUid:     node.Uid,
			NodeUid:       node.Uid,
			RelativeLevel: 0,
			CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
		},
	}

	if node.ParentUid != "" {
		ancestorsTreeRelation, err := network.GetNodeAncestors(node.ParentUid)
		if err != nil {
			log.Printf("error getting node %s ancestors: %s", node.Uid, err.Error())
			return err
		}

		for _, relation := range ancestorsTreeRelation {
			relations = append(relations, models.NetworkTreeRelation{
				RootUid:       relation.RootUid,
				ParentUid:     node.ParentUid,
				NodeUid:       node.Uid,
				RelativeLevel: relation.RelativeLevel + 1,
				CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			})
		}
	}

	err := lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, relations)
	if err != nil {
		log.Printf("insert error: %s", err.Error())
		return err
	}

	return nil
}
