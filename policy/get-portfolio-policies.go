package policy

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	md "github.com/wopta/goworkspace/policy/models"
)

type GetPoliciesReq struct {
	Queries []models.Query `json:"queries,omitempty"`
	Limit   int            `json:"limit"`
	Page    int            `json:"page"`
}

type GetPortfolioPoliciesResp struct {
	Policies []md.PolicyInfo `json:"policies"`
}

func GetPortfolioPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req  GetPoliciesReq
		resp = GetPortfolioPoliciesResp{
			Policies: make([]md.PolicyInfo, 0),
		}
		nodeUid string
	)

	log.SetPrefix("[GetSubtreePortfolioFx] ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error getting authToken: %s", err.Error())
		return "", nil, err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshiling request: %s", err.Error())
		return "", nil, err
	}

	if authToken.Role != models.UserRoleAdmin {
		nodeUid = authToken.UserID
	} else {
		toBeRemoveIndex := -1
		for index, query := range req.Queries {
			if query.Field == "producerUid" {
				nodeUid = query.Value.(string)
				toBeRemoveIndex = index
				break
			}
		}
		if toBeRemoveIndex != -1 {
			req.Queries = append(req.Queries[:toBeRemoveIndex], req.Queries[toBeRemoveIndex+1:]...)
		}
	}

	producersList, err := getProducersList(nodeUid)
	if err != nil {
		return "", nil, err
	}

	resp.Policies, err = getPortfolioPolicies(producersList, req.Queries, req.Limit)
	if err != nil {
		log.Printf("error query: %s", err.Error())
		return "", nil, err
	}
	log.Printf("found %02d policies", len(resp.Policies))

	rawResp, err := json.Marshal(resp)

	log.Println("Handler end -------------------------------------------------")

	return string(rawResp), resp, err
}

func getPortfolioPolicies(producers []string, requestQueries []models.Query, limit int) ([]md.PolicyInfo, error) {
	var (
		err        error
		fieldName  = "producerUid"
		limitValue = 25
		queries    []models.Query
	)
	if len(requestQueries) == 0 {
		return nil, errors.New("no query specified")
	}

	if limit != 0 {
		limitValue = limit
	}

	for _, q := range requestQueries {
		if q.Field == fieldName {
			log.Println("[getPortfolioPolicies] WARNING query with following field will be ignored: ", fieldName)
			continue
		} else {
			queries = append(queries, q)
		}
	}

	if len(producers) != 0 {
		values := make([]interface{}, 0)
		for _, p := range producers {
			values = append(values, p)
		}

		queries = append(queries, models.Query{
			Field:  fieldName,
			Op:     "in",
			Values: values,
		})
	}

	policies, err := getPoliciesInfoQueriesBigQuery(queries, limitValue)

	return policies, err
}

func getProducersList(nodeUid string) ([]string, error) {
	if nodeUid == "" {
		return []string{}, nil
	}
	node, err := network.GetNodeByUid(nodeUid)
	if err != nil {
		log.Printf("error fetching node %s from Firestore: %s", nodeUid, err.Error())
		return nil, err
	}

	children, err := network.GetNodeChildren(nodeUid)
	if err != nil {
		log.Printf("error fetching node %s children: %s", node.Uid, err.Error())
		return nil, err
	}

	producers := []string{nodeUid}

	for _, child := range children {
		producers = append(producers, child.NodeUid)
	}
	return producers, nil
}
