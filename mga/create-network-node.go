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
		request *models.NetworkNode
		err     error
	)

	log.Println("[CreateNetworkNodeFx] Handler start -------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[CreateNetworkNodeFx] request body: %s", string(body))
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("[CreateNetworkNodeFx] error unmarshaling request: %s", err.Error())
		return "", "", err
	}

	log.Println("[CreateNetworkNodeFx] creating network node into Firestore...")

	node, err := network.CreateNode(*request)
	if err != nil {
		log.Println("[CreateNetworkNodeFx] error creating network node into Firestore...")
		return "", "", err
	}
	log.Printf("[CreateNetworkNodeFx] network node created with uid %s", node.Uid)

	node.SaveBigQuery(origin)

	log.Println("[CreateNetworkNodeFx] network node successfully created!")

	models.CreateAuditLog(r, string(body))

	return "{}", "", err
}
