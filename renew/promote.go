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
	Error        string               `json:"error,omitempty"`
}

func PromoteFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err        error
		targetDate time.Time = time.Now().UTC()
		request    PromoteReq
		response   PromoteResp
	)

	log.SetPrefix("[PromoteFx] ")
	defer log.Println("Handler end -------------------------------------------")
	defer log.SetPrefix("")

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

	response, err = Promote(targetDate)
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

func Promote(targetDate time.Time) (PromoteResp, error) {
	var (
		err       error
		okChannel chan RenewReport = make(chan RenewReport)
		koChannel chan RenewReport = make(chan RenewReport)
		wg        sync.WaitGroup
		response  PromoteResp
	)

	query, params := buildQuery(targetDate)
	policies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query, params)
	if err != nil {
		log.Printf("error querying bigquery: %s", err.Error())
		return PromoteResp{}, err
	}

	wg.Add(len(policies))
	for _, p := range policies {
		if p.IsPay {
			go promotePolicyData(p, okChannel, koChannel, &wg)
		} else {
			go setPolicyNotPaid(p, okChannel, koChannel, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(okChannel)
		close(koChannel)
	}()

	for res := range okChannel {
		response.Success = append(response.Success, res)
	}
	for res := range koChannel {
		response.Failure = append(response.Failure, res)
	}

	return response, err
}

func buildQuery(date time.Time) (string, map[string]interface{}) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	// SELECT * FROM `wopta.renewPolicyDraft` WHERE EXTRACT(MONTH FROM startDate) = @date.Month() AND EXTRACT(DAY FROM startDate) = @date.Day()
	params["month"] = int64(date.Month())
	params["day"] = int64(date.Day())

	query.WriteString("SELECT * FROM `wopta.policiesView` WHERE " +
		"EXTRACT(MONTH FROM startDate) = @month AND " +
		"EXTRACT(DAY FROM startDate) = @day")

	return query.String(), params
}

func promotePolicyData(p models.Policy, okChannel, koChannel chan<- RenewReport, wg *sync.WaitGroup) {
	defer wg.Done()

	r := RenewReport{
		Policy: p,
	}

	trs, err := GetTransactionsByPolicyAnnuity(p.Uid, p.Annuity)
	if err != nil {
		log.Printf("error: %s", err.Error())
		r.Error = err.Error()
		koChannel <- r
		return
	}

	r.Transactions = trs
	p.Status = "RENEWED"
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	fireBatch := map[string]map[string]interface{}{
		"policyRenewedTest": {
			p.Uid: p,
		},
		"transactionsRenewdTest": {},
	}
	for idx, tr := range trs {
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		fireBatch["transactionsRenewdTest"][tr.Uid] = tr
		trs[idx] = tr
	}

	err = lib.SetBatchFirestoreErr(fireBatch)
	if err != nil {
		r.Error = err.Error()
		koChannel <- r
		return
	}

	// err = lib.InsertRowsBigQuery(models.WoptaDataset, models.PolicyCollection, p)
	// if err != nil {
	// 	koChannel <- r
	// 	return
	// }

	// err = lib.InsertRowsBigQuery(models.WoptaDataset, models.TransactionsCollection, trs)
	// if err != nil {
	// 	koChannel <- r
	// 	return
	// }

	okChannel <- r
}

func setPolicyNotPaid(p models.Policy, okChannel, koChannel chan<- RenewReport, wg *sync.WaitGroup) {
	defer wg.Done()
	p.Status = "INSOLUTO"
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()

	r := RenewReport{
		Policy: p,
	}
	fireBatch := map[string]map[string]interface{}{
		"policyRenewedTest": {
			p.Uid: p,
		},
	}

	err := lib.SetBatchFirestoreErr(fireBatch)
	if err != nil {
		r.Error = err.Error()
		koChannel <- r
		return
	}

	// err = lib.InsertRowsBigQuery(models.WoptaDataset, models.PolicyCollection, p)
	// if err != nil {
	// 	koChannel <- r
	// 	return
	// }

	okChannel <- r
}

// transactions
func GetTransactionsByPolicyAnnuity(policyUid string, annuity int) ([]models.Transaction, error) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	// SELECT * FROM `wopta.renewTransactionsDraft` WHERE policyUid = '@policyUid' AND annuity = @annuity
	params["policyUid"] = policyUid
	params["annuity"] = annuity

	query.WriteString("SELECT * FROM `wopta.transactionsView` WHERE " +
		"policyUid = '@policyUid' AND " +
		"annuity = @annuity")

	return lib.QueryParametrizedRowsBigQuery[models.Transaction](query.String(), params)
}
