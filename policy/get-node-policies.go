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

type GetNodePoliciesReq struct {
	NodeUid string `json:"nodeUid"`
}

func GetNodePoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		req      GetNodePoliciesReq
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

	if authToken.Role != models.UserRoleAdmin && authToken.UserID != req.NodeUid && !node.IsParentOf(req.NodeUid) {
		log.Printf("cannot access to node %s policies", req.NodeUid)
		return "", nil, errors.New("cannot access to node policies")
	}

	// TODO: implement query on BigQuery
	docsnap := lib.WhereFirestore(models.PolicyCollection, "producerUid", "==", req.NodeUid)
	tmpPolicies := models.PolicyToListData(docsnap)

	for _, p := range tmpPolicies {
		policies = append(policies, PolicyInfo{
			Uid:            p.Uid,
			ProductName:    p.Name,
			CodeCompany:    p.CodeCompany,
			ProposalNumber: p.ProposalNumber,
			NameDesc:       p.NameDesc,
			Status:         p.Status,
			Contractor:     p.Contractor.Name + " " + p.Contractor.Surname,
			Price:          p.PriceGross,
			PriceMonthly:   p.PriceGrossMonthly,
			Producer:       node.GetName(),
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			PaymentSplit:   p.PaymentSplit,
		})
	}

	rawPolices, err := json.Marshal(policies)

	log.Println("Handler End -------------------------------------------------")

	return string(rawPolices), policies, err
}
