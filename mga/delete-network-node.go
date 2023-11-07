package mga

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func DeleteNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	log.Println("[DeleteNetworkNodeFx] Handler start -------------------------")

	origin := r.Header.Get("Origin")
	nodeUid := r.Header.Get("uid")

	log.Printf("[DeleteNetworkNodeFx] deleting node %s from firestore...", nodeUid)

	err = network.DeleteNetworkNodeByUid(origin, nodeUid)
	if err != nil {
		log.Printf("[DeleteNetworkNodeFx] error deleting node %s from firestore", nodeUid)
		return "", "", err
	}

	log.Printf("[DeleteNetworkNodeFx] node %s deleted from firestore...", nodeUid)

	models.CreateAuditLog(r, "")

	return "", "", nil
}
