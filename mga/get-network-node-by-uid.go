package mga

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/network"
)

func GetNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetNetworkNodeByUidFx] ")

	log.Println("Handler start -----------------------------------------------")

	nodeUid := r.Header.Get("uid")
	log.Printf("Uid %s", nodeUid)

	networkNode := network.GetNetworkNodeByUid(nodeUid)

	jsonOut, err := networkNode.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), networkNode, err
}
