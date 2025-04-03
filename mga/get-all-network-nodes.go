package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type GetAllNetworkNodesResponse struct {
	NetworkNodes []models.NetworkNode `json:"networkNodes"`
}

func GetAllNetworkNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetAllNetworkNodesResponse

	log.AddPrefix("[GetAllNetworkNodesFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	networkNodes, err := network.GetAllNetworkNodes()
	if err != nil {
		log.ErrorF("error getting network nodes: %s", err.Error())
		return "", nil, err
	}

	// DO NOT EXPOSE CONFIGS
	for i := range networkNodes {
		networkNodes[i].JwtConfig = lib.JwtConfig{}
		networkNodes[i].CallbackConfig = nil
	}

	response.NetworkNodes = networkNodes
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		log.ErrorF("error marshalling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}
