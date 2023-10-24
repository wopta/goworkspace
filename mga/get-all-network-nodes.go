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
	log.Println("[GetAllNetworkNodesFx] Handler start -----------------------")
	var response GetAllNetworkNodesResponse

	networkNodes, err := network.GetAllNetworkNodes()
	if err != nil {
		log.Printf("[GetAllNetworkNodesFx] error getting network nodes: %s", err.Error())
		return "", nil, err
	}

	response.NetworkNodes = networkNodes
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		log.Printf("[GetAllNetworkNodesFx] error marshalling response: %s", err.Error())
		return "", nil, err
	}

	return string(responseBytes), response, nil
}
