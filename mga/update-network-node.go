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
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	body = lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &inputNode)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	if err := network.TestNetworkNodeUniqueness(inputNode.Code); err != nil {
		log.Printf("error validating node code: %s", err)
		return "", "", err
	}

	inputNode.Normalize()

	err = network.UpdateNode(inputNode)
	if err != nil {
		log.Printf("error updating network node %s", inputNode.Uid)
		return "", nil, err
	}

	log.Printf("network node %s updated successfully", inputNode.Uid)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
