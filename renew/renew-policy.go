package renew

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"time"
)

type RenewReport struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
	Error        string               `json:"error,omitempty"`
}

type RenewPolicyReq struct {
	PolicyUid string `json:"policyUid"`
}

type RenewPolicyResp struct {
	Success []RenewReport `json:"success"`
	Failure []RenewReport `json:"failure"`
}

func RenewPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err  error
		req  RenewPolicyReq
		resp = RenewPolicyResp{
			Success: make([]RenewReport, 0),
			Failure: make([]RenewReport, 0),
		}
	)

	log.SetPrefix("[RenewPolicyFx] ")
	defer func() {
		log.SetPrefix("")
		log.Println("Handler end -------------------------------------------------")
	}()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	policyType := chi.URLParam(r, "policyType")
	if policyType == "" {
		log.Printf("no policyType specified")
		return "", "", errors.New("no policyType specified")
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return "", nil, err
	}

	// TODO: solve issue that non active products are not fetched
	products := getProductsByPolicyType(policyType)

	policies, err := getPolicies(req.PolicyUid, products)

	log.Printf("found %02d policies", len(policies))

	rawResp, err := json.Marshal(resp)

	return string(rawResp), resp, err
}

func getPolicies(policyUid string, products []models.Product) ([]models.Policy, error) {
	var (
		err      error
		query    bytes.Buffer
		params   = make(map[string]interface{})
		policies []models.Policy
	)

	query.WriteString("SELECT * FROM `wopta.policiesView` WHERE ")

	if policyUid != "" {
		query.WriteString(" uid = @policyUid ")
		params["policyUid"] = policyUid
	} else if len(products) > 1 {
		//today := time.Now().UTC()
		for index, product := range products {
			if index != 0 {
				query.WriteString(" OR ")
			}
			//targetDate := today.AddDate(0, 0, product.RenewOffset)
			// TODO: restore commented lines
			targetDate := time.Date(2024, 03, 21, 0, 0, 0, 0, time.UTC)
			productNameKey := fmt.Sprintf("%sProductName", product.Name)
			productVersionKey := fmt.Sprintf("%sProductVersion", product.Version)
			targetMonthKey := fmt.Sprintf("%s%sMonth", product.Name, product.Version)
			targetDayKey := fmt.Sprintf("%s%sDay", product.Name, product.Version)
			params[productNameKey] = product.Name
			params[productVersionKey] = product.Version
			params[targetMonthKey] = int64(targetDate.Month())
			params[targetDayKey] = int64(targetDate.Day())
			query.WriteString("(name = @" + productNameKey)
			query.WriteString(" AND productVersion = @" + productVersionKey)
			query.WriteString(" AND EXTRACT(MONTH FROM startDate) = @" + targetMonthKey)
			query.WriteString(" AND EXTRACT(DAY FROM startDate) = @" + targetDayKey + ")")
		}
	}

	policies, err = lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		log.Printf("error getting policies: %v", err)
		return nil, err
	}

	return policies, nil
}
