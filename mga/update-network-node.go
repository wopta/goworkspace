package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"io"
	"log"
	"net/http"
)

func UpdateNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		body      []byte
		inputNode models.NetworkNode
	)

	log.Println("[UpdateNetworkNodeFx] Handler start --------------------------")

	body = lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[UpdateNetworkNodeFx] request body: %s", string(body))
	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.Printf("[UpdateNetworkNodeFx] error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	err = network.UpdateNode(inputNode)
	if err != nil {
		log.Printf("[UpdateNetworkNodeFx] error updating network node %s", inputNode.Uid)
		return "", nil, err
	}

	log.Printf("[UpdateNetworkNodeFx] network node %s updated successfully", inputNode.Uid)

	models.CreateAuditLog(r, string(body))

	return "{}", nil, nil
}
