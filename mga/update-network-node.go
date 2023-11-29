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

func UpdateNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		body      []byte
		inputNode models.NetworkNode
	)

	log.SetPrefix("[UpdateNetworkNodeFx] ")

	log.Println("Handler start -----------------------------------------------")

	body = lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("request body: %s", string(body))
	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	err = network.UpdateNode(inputNode)
	if err != nil {
		log.Printf("error updating network node %s", inputNode.Uid)
		return "", nil, err
	}

	log.Printf("network node %s updated successfully", inputNode.Uid)

	models.CreateAuditLog(r, string(body))

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
