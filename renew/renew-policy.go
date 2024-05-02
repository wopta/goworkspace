package renew

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
	"github.com/wopta/goworkspace/transaction"
	"io"
	"log"
	"net/http"
	"sync"
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
		wg   = new(sync.WaitGroup)
		req  RenewPolicyReq
		resp = RenewPolicyResp{
			Success: make([]RenewReport, 0),
			Failure: make([]RenewReport, 0),
		}
		productsMap = make(map[string]models.Product)
	)

	log.SetPrefix("[RenewPolicyFx] ")
	defer func() {
		log.SetPrefix("")
		log.Println("Handler end -------------------------------------------------")
	}()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	policyType, quoteType, err := getQueryParameters(r)
	if err != nil {
		return "", nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return "", nil, err
	}

	// TODO: solve issue that non active products are not fetched
	productsMap = getProductsMapByPolicyType(policyType, quoteType)

	policies, err := getPolicies(req.PolicyUid, policyType, quoteType, productsMap)
	if err != nil {
		log.Printf("error getting policies: %v", err)
		return "", nil, err
	}
	log.Printf("found %02d policies", len(policies))

	ch := make(chan RenewReport, len(policies))

	for _, policy := range policies {
		wg.Add(1)
		key := fmt.Sprintf("%s-%s", policy.Name, policy.ProductVersion)
		go draft(policy, productsMap[key], ch, wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.Error != "" {
			resp.Failure = append(resp.Failure, res)
			continue
		}
		resp.Success = append(resp.Success, res)
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), resp, err
}

func getQueryParameters(r *http.Request) (string, string, error) {
	policyType := r.URL.Query().Get("policyType")
	if policyType == "" {
		log.Printf("no policyType specified")
		return "", "", errors.New("no policyType specified")
	}

	quoteType := r.URL.Query().Get("quoteType")
	if quoteType == "" {
		log.Printf("no quoteType specified")
		return "", "", errors.New("no quoteType specified")
	}
	return policyType, quoteType, nil
}

func getPolicies(policyUid, policyType, quoteType string, products map[string]models.Product) ([]models.Policy, error) {
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

	} else if len(products) > 0 {
		today := time.Now().UTC()

		tmpProducts := lib.GetMapValues(products)
		params["isRenewable"] = true
		params["policyType"] = policyType
		params["quoteType"] = quoteType
		params["isPay"] = true
		params["isDeleted"] = false
		params["year"] = int64(today.Year())

		query.WriteString("isRenewable = @isRenewable")
		query.WriteString(" AND policyType = @policyType")
		query.WriteString(" AND quoteType = @quoteType")
		query.WriteString(" AND isPay = @isPay")
		query.WriteString(" AND isDeleted = @isDeleted")
		query.WriteString(" AND (")
		for index, product := range tmpProducts {
			if index != 0 {
				query.WriteString(" OR ")
			}
			//targetDate := today.AddDate(0, 0, product.RenewOffset)
			// TODO: restore commented lines
			targetDate := time.Date(2024, 03, 21, 0, 0, 0, 0, time.UTC)

			productNameKey := fmt.Sprintf("%s%sProductName", product.Name, product.Version)
			productVersionKey := fmt.Sprintf("%s%sProductVersion", product.Name, product.Version)
			targetYearKey := fmt.Sprintf("%s%sYear", product.Name, product.Version)
			targetMonthKey := fmt.Sprintf("%s%sMonth", product.Name, product.Version)
			targetDayKey := fmt.Sprintf("%s%sDay", product.Name, product.Version)
			targetDateKey := fmt.Sprintf("%s%sDate", product.Name, product.Version)
			params[productNameKey] = product.Name
			params[productVersionKey] = product.Version
			params[targetYearKey] = int64(targetDate.Year())
			params[targetMonthKey] = int64(targetDate.Month())
			params[targetDayKey] = int64(targetDate.Day())
			params[targetDateKey] = lib.GetBigQueryNullDateTime(targetDate)
			query.WriteString("(name = @" + productNameKey)
			query.WriteString(" AND productVersion = @" + productVersionKey)
			query.WriteString(" AND endDate != @" + targetDateKey)
			query.WriteString(" AND EXTRACT(MONTH FROM startDate) = @" + targetMonthKey)
			query.WriteString(" AND EXTRACT(DAY FROM startDate) = @" + targetDayKey + ")")
		}
		query.WriteString(")")
	}

	policies, err = lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		log.Printf("error getting policies: %v", err)
		return nil, err
	}

	policies = lib.SliceMap(policies, func(p models.Policy) models.Policy {
		// TODO: check if is it better to do so or is better to convert all bigquery fields to json fields
		var tmpPolicy models.Policy
		err = json.Unmarshal([]byte(p.Data), &tmpPolicy)
		return tmpPolicy
	})

	return policies, nil
}

func draft(policy models.Policy, product models.Product, ch chan<- RenewReport, wg *sync.WaitGroup) {
	var (
		err          error
		r            RenewReport
		transactions []models.Transaction
	)

	defer func() {
		r.Policy = policy
		r.Transactions = transactions
		r.Error = err.Error()
		ch <- r
		wg.Done()
	}()

	policy.Annuity = policy.Annuity + 1

	calculatePrices(&policy)

	policy.IsPay = false
	policy.Status = models.TransactionStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyDraftRenew, models.TransactionStatusToPay)

	transactions = transaction.CreateTransactions(policy, product, func() string {
		return lib.NewDoc(models.TransactionsCollection)
	})

	// TODO: value of scheduleFirstRate depends on if customer has an active "mandato"
	payUrl, transactions, err := payment.Controller(policy, product, transactions, false)
	if err != nil {
		return
	}

	policy.PayUrl = payUrl
	policy.Updated = time.Now().UTC()
	policy.IsRenew = true

	// TODO save policy and transactions to Firestore

	// TODO: save policy and transaction to BigQuery

}

func calculatePrices(policy *models.Policy) {
	var priceGross, priceNett, taxAmount, priceGrossMonthly, priceNettMonthly, taxAmountMonthly float64

	for index, guarantee := range policy.Assets[0].Guarantees {
		if policy.Annuity > guarantee.Value.Duration.Year {
			policy.Assets[0].Guarantees[index].IsDeleted = true
			continue
		}

		priceGross += guarantee.Value.PremiumGrossYearly
		priceNett += guarantee.Value.PremiumNetYearly
		taxAmount += guarantee.Value.PremiumTaxAmountYearly
		priceGrossMonthly += guarantee.Value.PremiumGrossMonthly
		priceNettMonthly += guarantee.Value.PremiumNetMonthly
		taxAmountMonthly += guarantee.Value.PremiumTaxAmountMonthly
	}
	policy.PriceGross = priceGross
	policy.PriceNett = priceNett
	policy.TaxAmount = taxAmount
	policy.PriceGrossMonthly = priceGrossMonthly
	policy.PriceNettMonthly = priceNettMonthly
	policy.TaxAmountMonthly = taxAmountMonthly
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Tax = policy.TaxAmount
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Net = policy.PriceNett
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Gross = policy.PriceGross

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Tax = policy.TaxAmountMonthly
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Net = policy.PriceNettMonthly
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Gross = policy.PriceGrossMonthly
	}
}
