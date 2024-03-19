package mga

import (
	"encoding/json"
	"errors"
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

	if node.ParentUid != "" {
		err = createNodeRelation(*node)
		if err != nil {
			log.Printf("error creating node %s network relations: %s", node.Uid, err.Error())
			return "", nil, err
		}
	}

	log.Println("network node successfully created!")

	models.CreateAuditLog(r, string(body))

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
}

func createNodeRelation(node models.NetworkNode) error {
	parentNode := network.GetNetworkNodeByUid(node.ParentUid)
	if parentNode == nil {
		log.Printf("parent node with uid %s not found", node.ParentUid)
		return errors.New("parent node not found")
	}

	ancestorsTreeRelation, err := parentNode.GetAncestors()
	if err != nil {
		log.Printf("error getting node %s ancestors: %s", node.Uid, err.Error())
		return err
	}

	parentRelation := models.NetworkTreeRelation{
		RootUid:       parentNode.Uid,
		ParentUid:     parentNode.Uid,
		NodeUid:       node.Uid,
		RelativeLevel: 1,
		CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
	}
	err = lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, parentRelation)
	if err != nil {
		log.Printf("insert error: %s", err.Error())
		return err
	}

	for _, relation := range ancestorsTreeRelation {
		treeRelation := models.NetworkTreeRelation{
			RootUid:       relation.RootUid,
			ParentUid:     node.ParentUid,
			NodeUid:       node.Uid,
			RelativeLevel: relation.RelativeLevel + 1,
			CreationDate:  lib.GetBigQueryNullDateTime(time.Now().UTC()),
		}
		err = lib.InsertRowsBigQuery(models.WoptaDataset, models.NetworkTreeStructureTable, treeRelation)
		if err != nil {
			log.Printf("insert error: %s", err.Error())
			return err
		}
	}
	return nil
}
