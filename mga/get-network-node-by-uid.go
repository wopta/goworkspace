package mga

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
)

func GetNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetNetworkNodeByUidFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	nodeUid := chi.URLParam(r, "uid")
	log.Printf("Uid %s", nodeUid)

	networkNode := network.GetNetworkNodeByUid(nodeUid)

	// DO NOT EXPOSE CONFIGS
	networkNode.JwtConfig = lib.JwtConfig{}
	networkNode.CallbackConfig = nil

	jsonOut, err := networkNode.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), networkNode, err
}
