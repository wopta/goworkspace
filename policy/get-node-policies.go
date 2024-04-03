package policy

import (
	"encoding/json"
	"errors"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"io"
	"log"
	"net/http"
)

type GetNodePoliciesResp struct {
	Policies []PolicyInfo `json:"policies"`
}

func GetNodePoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		req      GetPoliciesReq
		resp     GetNodePoliciesResp
		policies = make([]PolicyInfo, 0)
	)

	log.SetPrefix("[GetNodePoliciesFx] ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error fetching authToken: %s", err.Error())
		return "", nil, err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("request: %s", string(body))
	defer r.Body.Close()
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	nodeUid := r.Header.Get("nodeUid")
	reqNode, err := network.GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("error fetching node %s from Firestore: %s", reqNode.Uid, err.Error())
		return "", nil, err
	}

	if !CanBeAccessedBy(authToken.Role, nodeUid, authToken.UserID) {
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

	log.Printf("response: %s", string(rawPolices))
	log.Println("Handler End -------------------------------------------------")

	return string(rawPolices), policies, err
}
