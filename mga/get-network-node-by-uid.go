package mga

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/network"
)

func GetNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetNetworkNodeByUidFx] Handler start -----------------------")

	nodeUid := r.Header.Get("uid")
	log.Printf("[GetNetworkNodeByUidFx] Uid %s", nodeUid)

	networkNode := network.GetNetworkNodeByUid(nodeUid)

	jsonOut, err := networkNode.Marshal()

	return string(jsonOut), networkNode, err
}
