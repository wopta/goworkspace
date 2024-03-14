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
	"time"
)

type PolicyInfo struct {
	Uid            string    `json:"uid" bigquery:"uid"`
	ProductName    string    `json:"productName" bigquery:"productName"`
	CodeCompany    string    `json:"codeCompany" bigquery:"codeCompany"`
	ProposalNumber int       `json:"proposalNumber" bigquery:"proposalNumber"`
	NameDesc       string    `json:"nameDesc" bigquery:"nameDesc"`
	Status         string    `json:"status" bigquery:"status"`
	Contractor     string    `json:"contractor" bigquery:"contractor"`
	Price          float64   `json:"price" bigquery:"price"`
	PriceMonthly   float64   `json:"priceMonthly" bigquery:"priceMonthly"`
	Producer       string    `json:"producer" bigquery:"producer"`
	ProducerCode   string    `json:"producerCode" bigquery:"-"`
	StartDate      time.Time `json:"startDate" bigquery:"startDate"`
	EndDate        time.Time `json:"endDate" bigquery:"endDate"`
	PaymentSplit   string    `json:"paymentSplit" bigquery:"paymentSplit"`
}

type GetSubtreePortfolioResp struct {
	Policies []PolicyInfo `json:"policies"`
}

func GetSubtreePortfolioFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req  GetPoliciesReq
		resp GetSubtreePortfolioResp
	)

	log.SetPrefix("[GetSubtreePortfolioFx] ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error getting authToken: %s", err.Error())
		return "", nil, err
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("request: %s", string(body))

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshiling request: %s", err.Error())
		return "", nil, err
	}

	producersMap, err := getProducersMap(authToken.Role, authToken.UserID)
	if err != nil {
		return "", nil, err
	}

	result, err := getPortfolioPoliciesV2(lib.GetMapKeys(producersMap), req.Queries, req.Limit)
	if err != nil {
		log.Printf("error query: %s", err.Error())
		return "", nil, err
	}
	log.Printf("found %02d policies", len(resp.Policies))

	for _, policy := range result {
		resp.Policies = append(resp.Policies, policyToPolicyInfo(policy, producersMap[policy.ProducerUid].Name))
	}

	rawResp, err := json.Marshal(resp)

	log.Printf("response: %s", string(rawResp))
	log.Println("Handler end -------------------------------------------------")

	return string(rawResp), resp, err
}

func getProducersMap(role string, nodeUid string) (map[string]models.NetworkTreeElement, error) {
	producersMap := make(map[string]models.NetworkTreeElement)
	if role != models.UserRoleAdmin {
		node, err := network.GetNodeByUid(nodeUid)
		if err != nil {
			log.Printf("error fetching node %s from Firestore: %s", nodeUid, err.Error())
			return nil, err
		}

		children, err := node.GetChildren()
		if err != nil {
			log.Printf("error fetching node %s children: %s", node.Uid, err.Error())
			return nil, err
		}

		producersMap[nodeUid] = models.NetworkTreeElement{
			ParentUid: node.ParentUid,
			NodeUid:   node.Uid,
			Name:      node.GetName(),
		}

		for _, child := range children {
			producersMap[child.NodeUid] = child
		}
	}
	return producersMap, nil
}

func getPortfolioPoliciesV2(producers []string, requestQueries []models.Query, limit int) ([]models.Policy, error) {
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

	values := make([]interface{}, 0)
	for _, p := range producers {
		values = append(values, p)
	}

	queries = append(queries, models.Query{
		Field:  fieldName,
		Op:     "in",
		Values: values,
	})

	policies, err := GetPoliciesByQueriesBigQuery(models.WoptaDataset, models.PoliciesViewCollection, queries, limitValue)
	if err != nil {
		return nil, err
	}

	return policies, err
}
