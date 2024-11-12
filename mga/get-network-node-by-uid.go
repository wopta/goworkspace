package mga

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
)

func GetNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("error: %+v", err.Error())
		}
		log.SetPrefix("")
		log.Println("Handler end -------------------------------------------------")
	}()

	log.SetPrefix("[GetNetworkNodeByUidFx] ")
	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	nodeUid := chi.URLParam(r, "uid")
	log.Printf("Uid %s", nodeUid)

	if authToken.IsNetworkNode && !network.IsParentOf(authToken.UserID, nodeUid) {
		return "", nil, errors.New("cannot access this node")
	}

	networkNode := network.GetNetworkNodeByUid(nodeUid)

	// DO NOT EXPOSE CONFIGS
	networkNode.JwtConfig = lib.JwtConfig{}
	networkNode.CallbackConfig = nil

	rawResp, err := networkNode.Marshal()

	return string(rawResp), networkNode, err
}
