package mga

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

type GetAllNetworkNodesResponse struct {
	NetworkNodes []models.NetworkNode `json:"networkNodes"`
}

func getAllNetworkNodesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
