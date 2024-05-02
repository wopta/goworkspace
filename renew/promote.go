package renew

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type PromoteReq struct {
	Date string `json:"date"`
}

type PromoteResp struct {
	Success []RenewReport `json:"success"`
	Failure []RenewReport `json:"failure"`
}

type RenewReport struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
	Error        error                `json:"error,omitempty"`
}

func PromoteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err        error
		targetDate time.Time = time.Now().UTC()
		request    PromoteReq
		response   PromoteResp
	)

	log.SetPrefix("[PromoteFx] ")
	log.Println("Handler start -----------------------------------------------")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(reqBytes, &request)
	if err != nil {
		log.Printf("error unmarshalling request: %s", err.Error())
		return "", nil, err
	}

	if request.Date != "" {
		if targetDate, err = time.Parse(time.DateOnly, request.Date); err != nil {
			log.Printf("error parsing request.Date: %s", err.Error())
			return "", nil, err
		}
	}

	policies, err := getRenewingPolicies(targetDate)
	if err != nil {
		log.Printf("error querying bigquery: %s", err.Error())
		return "", nil, err
	}

	response, err = Promote(policies, saveToDatabases)
	if err != nil {
		return "", nil, err
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return "", nil, err
	}
	_, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), fmt.Sprintf("renew/promote/promote-%d.json", time.Now().Unix()), responseJson)
	if err != nil {
		return "", nil, err
	}

	return string(responseJson), response, err
}

func Promote(policies []models.Policy, saveFn func(map[string]map[string]interface{}) error) (PromoteResp, error) {
	var (
		err            error
		promoteChannel chan RenewReport = make(chan RenewReport, len(policies))
		wg             sync.WaitGroup
		response       PromoteResp
	)

	wg.Add(len(policies))
	for _, p := range policies {
		if p.IsPay {
			go promotePolicyData(p, promoteChannel, &wg, saveFn)
		} else {
			go setPolicyNotPaid(p, promoteChannel, &wg, saveFn)
		}
	}

	go func() {
		wg.Wait()
		close(promoteChannel)
	}()

	for res := range promoteChannel {
		if res.Error != nil {
			response.Failure = append(response.Failure, res)
		} else {
			response.Success = append(response.Success, res)
		}
	}

	return response, err
}

func getRenewingPolicies(renewDate time.Time) ([]models.Policy, error) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	params["month"] = int64(renewDate.Month())
	params["day"] = int64(renewDate.Day())

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE "+
		"EXTRACT(MONTH FROM startDate) = @month AND "+
		"EXTRACT(DAY FROM startDate) = @day",
		models.WoptaDataset,
		renewPolicyCollection))

	policies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		return nil, err
	}

	for index, policy := range policies {
		var temp models.Policy
		err := json.Unmarshal([]byte(policy.Data), &temp)
		if err != nil {
			return nil, err
		}
		policies[index] = temp
	}

	return policies, nil
}

func getTransactionsByPolicyAnnuity(policyUid string, annuity int) ([]models.Transaction, error) {
	queries := []firestoreQuery{
		{field: "policyUid", operator: "==", queryValue: policyUid},
		{field: "annuity", operator: "==", queryValue: annuity},
	}

	return firestoreWhere[models.Transaction](renewTransactionCollection, queries)
}

func promotePolicyData(p models.Policy, promoteChannel chan<- RenewReport, wg *sync.WaitGroup, saveFn func(map[string]map[string]interface{}) error) {
	var (
		err          error
		transactions []models.Transaction
	)

	defer func() {
		promoteChannel <- RenewReport{
			Policy:       p,
			Transactions: transactions,
			Error:        err,
		}

		wg.Done()
	}()

	if transactions, err = getTransactionsByPolicyAnnuity(p.Uid, p.Annuity); err != nil {
		return
	}

	p.Status = policyStatusRenewed
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	batch := createSaveBatch(p, transactions)

	err = saveFn(batch)
}

func setPolicyNotPaid(p models.Policy, promoteChannel chan<- RenewReport, wg *sync.WaitGroup, saveFn func(map[string]map[string]interface{}) error) {
	var (
		err error
	)

	defer func() {
		promoteChannel <- RenewReport{
			Policy: p,
			Error:  err,
		}
		wg.Done()
	}()

	p.Status = policyStatusPaymentUnsolved
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	batch := createSaveBatch(p, nil)

	err = saveFn(batch)
}
