package renew

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

type RenewReq struct {
	Date            string `json:"date"`
	DaysBeforeRenew string `json:"days_before_renew"`
}

func RenewMailFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err             error
		date            = time.Now().UTC()
		daysBeforeRenew = 10
		query           bytes.Buffer
		params          = make(map[string]interface{})
		req             RenewReq
	)

	log.SetPrefix("[RenewMailFx] ")
	defer func() {
		collectionPrefix = ""
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
	}()
	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	if req.Date != "" {
		tmpDate, err := time.Parse(time.DateOnly, req.Date)
		if err != nil {
			log.Printf("error parsing request date: %s", err.Error())
			return "", nil, err
		}
		date = tmpDate
	}

	if req.DaysBeforeRenew != "" {
		tmpDays, err := strconv.Atoi(req.DaysBeforeRenew)
		if err != nil {
			log.Printf("error parsing target date: %s", err.Error())
			return "", nil, err
		}
		daysBeforeRenew = tmpDays
	}

	targetDate := date.AddDate(0, 0, daysBeforeRenew)
	log.Printf("Date is: %v; targetDate is: %v", date, targetDate)

	params["isRenewable"] = true
	params["channel"] = models.ECommerceChannel
	params["targetYear"] = int64(targetDate.Year())
	params["targetMonth"] = int64(targetDate.Month())
	params["targetDay"] = int64(targetDate.Day())

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE", lib.WoptaDataset, lib.RenewPolicyViewCollection))
	query.WriteString(" isRenewable = @isRenewable")
	query.WriteString(" AND channel = @channel")
	query.WriteString(" AND EXTRACT(YEAR FROM RenewDate) = @targetYear")
	query.WriteString(" AND EXTRACT(MONTH FROM RenewDate) = @targetMonth")
	query.WriteString(" AND EXTRACT(DAY FROM RenewDate) = @targetDay")

	policies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		log.Printf("error fetching policies from BigQuery: %s", err)
		return "", nil, err
	}

	policies = lib.SliceMap(policies, func(p models.Policy) models.Policy {
		var tmpPolicy models.Policy
		err = json.Unmarshal([]byte(p.Data), &tmpPolicy)
		return tmpPolicy
	})

	for _, policy := range policies {
		from := mail.AddressAnna
		to := mail.GetContractorEmail(&policy)
		flowName := models.ECommerceFlow
		log.Printf("Sending email from %s to %s", from, to)
		mail.SendMailRenewDraft(policy, from, to, mail.Address{}, flowName, policy.HasMandate)
	}

	return "", nil, nil
}
