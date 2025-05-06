package mga

import (
	"github.com/wopta/goworkspace/lib/log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/network"
)

func DeleteNetworkNodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	log.AddPrefix("[DeleteNetworkNodeFx] ")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	nodeUid := chi.URLParam(r, "uid")

	log.Printf("deleting node %s from firestore...", nodeUid)

	err = network.DeleteNetworkNodeByUid(origin, nodeUid)
	if err != nil {
		log.ErrorF("error deleting node %s from firestore", nodeUid)
		return "", "", err
	}

	log.Printf("node %s deleted from firestore...", nodeUid)

	log.Println("Handler end -------------------------------------------------")

	return "", "", nil
}
