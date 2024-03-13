package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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

	producers := make([]string, 0)
	if authToken.Role != models.UserRoleAdmin {
		producers = append(producers, authToken.UserID)
		children, err := getNodeChildren(authToken.UserID)
		if err != nil {
			log.Printf("error fetching node %s children: %s", authToken.UserID, err.Error())
			return "", nil, err
		}
		producers = append(producers, children...)
	}

	resp.Policies, err = getPortfolioPoliciesV2(producers, req.Queries, req.Limit)
	if err != nil {
		log.Printf("error query: %s", err.Error())
		return "", nil, err
	}
	log.Printf("found %02d policies", len(resp.Policies))

	rawResp, err := json.Marshal(resp)

	log.Printf("response: %s", string(rawResp))
	log.Println("Handler end -------------------------------------------------")

	return string(rawResp), resp, err
}

type Children struct {
	Uid string `bigquery:"childUid"`
}

func getNodeChildren(nodeUid string) ([]string, error) {
	baseQuery := fmt.Sprintf("SELECT nodeUid FROM `%s.%s` WHERE ", models.WoptaDataset, models.NetworkTreeStructureTable)
	whereClause := fmt.Sprintf("rootUid = '%s'", nodeUid)
	query := fmt.Sprintf("%s %s", baseQuery, whereClause)
	result, err := lib.QueryRowsBigQuery[Children](query)
	if err != nil {
		log.Printf("error fetching children from BigQuery for node %s: %s", nodeUid, err.Error())
		return nil, err
	}

	ch := make([]string, 0)
	for _, r := range result {
		ch = append(ch, r.Uid)
	}

	return ch, nil
}

func getPortfolioPoliciesV2(producerUid []string, requestQueries []models.Query, limit int) ([]PolicyInfo, error) {
	var (
		err error
	)
	if len(requestQueries) == 0 {
		return nil, errors.New("no query specified")
	}

	var (
		fieldName  = "producerUid"
		limitValue = 25
		queries    []models.Query
	)
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
	for _, p := range producerUid {
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

	result := make([]PolicyInfo, 0)
	for _, policy := range policies {
		result = append(result, PolicyInfo{
			Uid:            policy.Uid,
			ProductName:    policy.Name,
			CodeCompany:    policy.CodeCompany,
			ProposalNumber: policy.ProposalNumber,
			NameDesc:       policy.NameDesc,
			Status:         policy.Status,
			Contractor:     policy.Contractor.Name + " " + policy.Contractor.Surname,
			Price:          policy.PriceGross,
			PriceMonthly:   policy.PriceGrossMonthly,
			StartDate:      policy.StartDate,
			EndDate:        policy.EndDate,
			PaymentSplit:   policy.PaymentSplit,
		})
	}

	return result, err
}
