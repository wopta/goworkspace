package mga

import (
	"errors"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func getNetworkNodeByUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.ErrorF("error: %+v", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
	}()

	log.AddPrefix("[GetNetworkNodeByUidFx] ")
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
	if networkNode == nil {
		return "", nil, errors.New("node not found")
	}

	// DO NOT EXPOSE CONFIGS
	networkNode.JwtConfig = lib.JwtConfig{}
	networkNode.CallbackConfig = nil

	rawResp, err := networkNode.Marshal()

	return string(rawResp), networkNode, err
}
