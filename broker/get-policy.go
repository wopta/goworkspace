package broker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetPolicyFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")
	log.Printf("policyUid: %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error fetching authToken: %s", err.Error())
		return "", nil, err
	}

	policy, err := plc.GetPolicy(policyUid, lib.PolicyCollection)
	if err != nil {
		log.Printf("error fetching policy: %s", err.Error())
		return "", nil, err
	}

	switch authToken.Role {
	case lib.UserRoleCustomer:
		if policy.Contractor.Uid != authToken.UserID {
			log.Printf("error fetching policy: invalid user id: %s", authToken.UserID)
			return "", nil, errors.New("invalid user id")
		}
	case lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent:
		if policy.ProducerUid != authToken.UserID && !network.IsChildOf(authToken.UserID, policy.ProducerUid) {
			log.Printf("error fetching policy: invalid producer uid: %s", authToken.UserID)
			return "", nil, errors.New("invalid producer uid")
		}
	}

	res, err := json.Marshal(policy)

	log.Println("Handler end -------------------------------------------------")

	return string(res), policy, err
}
