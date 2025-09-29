package renew

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type PromoteReq struct {
	Date             string `json:"date"`
	DryRun           *bool  `json:"dryRun"`
	CollectionPrefix string `json:"collectionPrefix"`
}

func promoteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		dryRun     bool = true
		err        error
		targetDate time.Time = time.Now().UTC()
		request    PromoteReq
		response   = RenewResp{
			Success: make([]RenewReport, 0),
			Failure: make([]RenewReport, 0),
		}
	)

	log.AddPrefix("PromoteFx")

	defer func() {
		collectionPrefix = ""
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
	}()

	log.Println("Handler start -----------------------------------------------")

	reqBytes := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(reqBytes, &request)
	if err != nil {
		log.ErrorF("error unmarshalling request: %s", err.Error())
		return "", nil, err
	}

	if request.Date != "" {
		if targetDate, err = time.Parse(time.DateOnly, request.Date); err != nil {
			log.ErrorF("error parsing request.Date: %s", err.Error())
			return "", nil, err
		}
	}
	if request.DryRun != nil {
		dryRun = *request.DryRun
	}
	collectionPrefix = request.CollectionPrefix

	log.Printf("running pipeline with set config. TargetDate: %v, DryRun: %v", targetDate, dryRun)

	policies, err := getRenewingPolicies(targetDate)
	if err != nil {
		log.ErrorF("error querying bigquery: %s", err.Error())
		return "", nil, err
	}
	log.Printf("found %02d policies", len(policies))

	saveFn := func(p models.Policy, trs []models.Transaction) error {
		data := createPromoteSaveBatch(p, trs)

		if !dryRun {
			err = saveToDatabases(data)
			if err != nil {
				return err
			}

			dataDelete := createPromoteProcessedBatch(p, trs)
			p.AddSystemNote(models.GetQuietanzamentoPolicyNote)
			return saveToDatabases(dataDelete)
		}

		log.Println("dryRun active - not saving to DB")

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
		filename := fmt.Sprintf("renew/promote/report-promote-%s-%d.json", targetDate.Format(time.DateOnly), time.Now().Unix())
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
			res.Policy.AddSystemNote(models.GetErrorNote("Quietanzamento"))
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
	params["isDeleted"] = false
	params["isRenewable"] = true

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE "+
		"EXTRACT(MONTH FROM startDate) = @month AND "+
		"EXTRACT(DAY FROM startDate) = @day AND "+
		"isDeleted = @isDeleted AND "+
		"isRenewable = @isRenewable",
		models.WoptaDataset,
		collectionPrefix+lib.RenewPolicyViewCollection))

	log.Printf("query: %s", query.String())
	log.Printf("params: %v", params)

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
		{field: "isDelete", operator: "==", queryValue: false},
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

	transactions, err = getTransactionsByPolicyAnnuity(p.Uid, p.Annuity)
	if err != nil {
		return
	}

	p.Status = models.PolicyStatusRenewed
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	err = saveFn(p, transactions)
}

func setPolicyNotPaid(p models.Policy, promoteChannel chan<- RenewReport, wg *sync.WaitGroup, saveFn func(models.Policy, []models.Transaction) error) {
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

	transactions, err = getTransactionsByPolicyAnnuity(p.Uid, p.Annuity)
	if err != nil {
		return
	}
	transactions = lib.SliceMap(transactions, func(t models.Transaction) models.Transaction {
		t.UpdateDate = time.Now().UTC()
		return t
	})

	p.Status = models.PolicyStatusUnsolved
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	err = saveFn(p, transactions)
}
