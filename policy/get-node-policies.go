package policy

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/policy/models"
	"github.com/wopta/goworkspace/policy/utils"
)

type GetNodePoliciesResp struct {
	Policies []models.PolicyInfo `json:"policies"`
}

func GetNodePoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		req      GetPoliciesReq
		resp     GetNodePoliciesResp
		policies = make([]models.PolicyInfo, 0)
	)

	log.SetPrefix("[GetNodePoliciesFx] ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error fetching authToken: %s", err.Error())
		return "", nil, err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	nodeUid := chi.URLParam(r, "nodeUid")
	reqNode, err := network.GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("error fetching node %s from Firestore: %s", reqNode.Uid, err.Error())
		return "", nil, err
	}

	if !utils.CanBeAccessedBy(authToken.Role, nodeUid, authToken.UserID) {
		log.Printf("cannot access to node %s policies", nodeUid)
		return "", nil, errors.New("cannot access to node policies")
	}

	resp.Policies, err = getPortfolioPolicies([]string{nodeUid}, req.Queries, req.Limit)
	if err != nil {
		log.Printf("error query: %s", err.Error())
		return "", nil, err
	}
	log.Printf("found %02d policies", len(resp.Policies))

	rawPolices, err := json.Marshal(resp)

	log.Println("Handler End -------------------------------------------------")

	return string(rawPolices), policies, err
}
