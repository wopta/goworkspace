package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"io"
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

	log.AddPrefix("UpdateNetworkNodeFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	body = lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.ErrorF("error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	inputNode.Normalize()

	err = network.UpdateNode(inputNode)
	if err != nil {
		log.ErrorF("error updating network node %s", inputNode.Uid)
		return "", nil, err
	}

	log.Printf("network node %s updated successfully", inputNode.Uid)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
