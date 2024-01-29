package mga

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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

	inputNode.Sanitize()

	node, err := network.CreateNode(*inputNode)
	if err != nil {
		log.Println("error creating network node into Firestore...")
		return "", "", err
	}
	log.Printf("network node created with uid %s", node.Uid)

	node.SaveBigQuery(origin)

	log.Println("network node successfully created!")

	models.CreateAuditLog(r, string(body))

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
}
