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
	Date             string `json:"date"`
	DryRun           *bool  `json:"dryRun"`
	CollectionPrefix string `json:"collectionPrefix"`
}

func PromoteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		dryRun     bool = true
		err        error
		targetDate time.Time = time.Now().UTC()
		request    PromoteReq
		response   RenewResp
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
	if request.DryRun != nil {
		dryRun = *request.DryRun
	}
	collectionPrefix = request.CollectionPrefix

	policies, err := getRenewingPolicies(targetDate)
	if err != nil {
		log.Printf("error querying bigquery: %s", err.Error())
		return "", nil, err
	}

	saveFn := func(p models.Policy, trs []models.Transaction) error {
		data := createPromoteSaveBatch(p, trs)

		if !dryRun {
			return saveToDatabases(data)
		}

		return nil
	}

	response, err = Promote(policies, saveFn)
	if err != nil {
		return "", nil, err
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return "", nil, err
	}
	if !dryRun {
		filename := fmt.Sprintf("renew/promote/report-%s-%d.json", targetDate.Format(time.DateOnly), time.Now().Unix())
		if _, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, responseJson); err != nil {
			return "", nil, err
		}
	}

	sendReportMail(targetDate, response, false)

	return string(responseJson), response, err
}

func Promote(policies []models.Policy, saveFn func(models.Policy, []models.Transaction) error) (RenewResp, error) {
	var (
		err            error
		promoteChannel chan RenewReport = make(chan RenewReport, len(policies))
		wg             sync.WaitGroup
		response       RenewResp
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
		if res.Error != "" {
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
		collectionPrefix+lib.RenewPolicyCollection))

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

	return firestoreWhere[models.Transaction](collectionPrefix+lib.RenewTransactionCollection, queries)
}

func promotePolicyData(p models.Policy, promoteChannel chan<- RenewReport, wg *sync.WaitGroup, saveFn func(models.Policy, []models.Transaction) error) {
	var (
		err          error
		transactions []models.Transaction
	)

	defer func() {
		var r RenewReport

		r.Policy = p
		r.Transactions = transactions
		if err != nil {
			r.Error = err.Error()
		}
		promoteChannel <- r

		wg.Done()
	}()

	if transactions, err = getTransactionsByPolicyAnnuity(p.Uid, p.Annuity); err != nil {
		return
	}

	p.Status = policyStatusRenewed
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	err = saveFn(p, transactions)
}

func setPolicyNotPaid(p models.Policy, promoteChannel chan<- RenewReport, wg *sync.WaitGroup, saveFn func(models.Policy, []models.Transaction) error) {
	var (
		err error
	)

	defer func() {
		var r RenewReport

		r.Policy = p
		if err != nil {
			r.Error = err.Error()
		}
		promoteChannel <- r
		wg.Done()
	}()

	p.Status = policyStatusPaymentUnsolved
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	err = saveFn(p, nil)
}
