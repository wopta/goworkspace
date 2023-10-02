package mga

import (
	"log"
	"net/http"
)

func GetNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetNetworkNodeByUidFx] Handler start -----------------------")

	nodeUid := r.Header.Get("uid")
	log.Printf("[GetNetworkNodeByUidFx] Uid %s", nodeUid)

	// Call network domain

	return "", "", nil
}
