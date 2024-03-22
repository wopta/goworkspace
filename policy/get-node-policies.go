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

func GetNodePoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		req      GetPoliciesReq
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

	node, err := network.GetNodeByUid(authToken.UserID)
	if err != nil {
		log.Printf("error fetching node %s from Firestore: %s", authToken.UserID, err.Error())
		return "", nil, err
	}

	nodeUid := r.Header.Get("nodeUid")
	reqNode, err := network.GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("error fetching node %s from Firestore: %s", reqNode.Uid, err.Error())
		return "", nil, err
	}

	if authToken.Role != models.UserRoleAdmin && authToken.UserID != nodeUid && !node.IsParentOf(nodeUid) {
		log.Printf("cannot access to node %s policies", nodeUid)
		return "", nil, errors.New("cannot access to node policies")
	}

	result, err := getPortfolioPolicies([]string{nodeUid}, req.Queries, req.Limit)

	producerName := reqNode.GetName()
	for _, p := range result {
		policies = append(policies, policyToPolicyInfo(p, producerName))
	}

	rawPolices, err := json.Marshal(policies)

	log.Println("Handler End -------------------------------------------------")

	return string(rawPolices), policies, err
}
