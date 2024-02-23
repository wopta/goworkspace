package mga

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/network"
)

func DeleteNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	log.SetPrefix("[DeleteNetworkNodeFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	nodeUid := r.Header.Get("uid")

	log.Printf("deleting node %s from firestore...", nodeUid)

	err = network.DeleteNetworkNodeByUid(origin, nodeUid)
	if err != nil {
		log.Printf("error deleting node %s from firestore", nodeUid)
		return "", "", err
	}

	log.Printf("node %s deleted from firestore...", nodeUid)

	log.Println("Handler end -------------------------------------------------")

	return "", "", nil
}
