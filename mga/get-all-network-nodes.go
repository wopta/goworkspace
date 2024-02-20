package mga

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type GetAllNetworkNodesResponse struct {
	NetworkNodes []models.NetworkNode `json:"networkNodes"`
}

func GetAllNetworkNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetAllNetworkNodesResponse

	log.SetPrefix("[GetAllNetworkNodesFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	networkNodes, err := network.GetAllNetworkNodes()
	if err != nil {
		log.Printf("error getting network nodes: %s", err.Error())
		return "", nil, err
	}

	response.NetworkNodes = networkNodes
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		log.Printf("error marshalling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}
