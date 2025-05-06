package broker

import (
	"encoding/json"
	"errors"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/policy/utils"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("GetPolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")
	log.Printf("policyUid: %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.ErrorF("error fetching authToken: %s", err.Error())
		return "", nil, err
	}

	policy, err := plc.GetPolicy(policyUid, lib.PolicyCollection)
	if err != nil {
		log.ErrorF("error fetching policy: %s", err.Error())
		return "", nil, err
	}

	if authToken.IsNetworkNode && !utils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, authToken.UserID) {
		log.ErrorF("error fetching policy invalid producer uid: %s", authToken.UserID)
		return "", nil, errors.New("invalid producer uid")
	}

	res, err := json.Marshal(policy)

	log.Println("Handler end -------------------------------------------------")

	return string(res), policy, err
}
