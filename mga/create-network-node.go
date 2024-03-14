package mga

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"io"
	"log"
	"net/http"
	"time"
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
	log.Printf("request body: %s", string(body))
	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", "", err
	}
	// TODO: check node.Type in warrant.AllowedTypes
	// TODO: check unique node.Code
	// TODO: check unique companyCode for company

	log.Println("creating network node into Firestore...")

	inputNode.Normalize()

	node, err := network.CreateNode(*inputNode)
	if err != nil {
		log.Println("error creating network node into Firestore...")
		return "", "", err
	}
	log.Printf("network node created with uid %s", node.Uid)

	node.SaveBigQuery(origin)

	//node := inputNode
	if node.ParentUid != "" {
		parentNode := network.GetNetworkNodeByUid(node.ParentUid)
		if parentNode == nil {
			log.Printf("parent node with uid %s not found", node.ParentUid)
			return "", nil, errors.New("parent node not found")
		}

		ancestors, err := parentNode.GetAncestors()
		if err != nil {
			log.Printf("error getting node %s ancestors: %s", node.Uid, err.Error())
			return "", nil, err
		}

		if len(ancestors) == 0 {
			treeRelation := models.NetworkTreeElement{
				RootUid:       parentNode.Uid,
				ParentUid:     parentNode.Uid,
				NodeUid:       node.Uid,
				Name:          node.GetName(),
				AbsoluteLevel: 2,
				RelativeLevel: 1,
				CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			}
			err = lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, treeRelation)
			if err != nil {
				log.Printf("insert error: %s", err.Error())
				return "", nil, err
			}
		} else {
			parentAbsoluteLevel := ancestors[0].AbsoluteLevel
			parentRelation := models.NetworkTreeElement{
				RootUid:       parentNode.Uid,
				ParentUid:     parentNode.Uid,
				NodeUid:       node.Uid,
				Name:          node.GetName(),
				AbsoluteLevel: parentAbsoluteLevel + 1,
				RelativeLevel: 1,
				CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
			}
			err = lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, parentRelation)
			if err != nil {
				log.Printf("insert error: %s", err.Error())
				return "", nil, err
			}

			for _, ancestor := range ancestors {
				treeRelation := models.NetworkTreeElement{
					RootUid:       ancestor.RootUid,
					ParentUid:     node.ParentUid,
					NodeUid:       node.Uid,
					Name:          node.GetName(),
					AbsoluteLevel: ancestor.AbsoluteLevel + 1,
					RelativeLevel: ancestor.RelativeLevel + 1,
					CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
				}
				err = lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, treeRelation)
				if err != nil {
					log.Printf("insert error: %s", err.Error())
					return "", nil, err
				}
			}
		}
	}

	log.Println("network node successfully created!")

	models.CreateAuditLog(r, string(body))

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
}
