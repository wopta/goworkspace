package renew

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
	"github.com/wopta/goworkspace/transaction"
)

type DraftReq struct {
	PolicyUid        string `json:"policyUid"`
	Date             string `json:"date"`
	DryRun           *bool  `json:"dryRun"`
	CollectionPrefix string `json:"collectionPrefix"`
}

func DraftFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		dryRun = true
		wg     = new(sync.WaitGroup)
		req    DraftReq
		resp   = RenewResp{
			Success: make([]RenewReport, 0),
			Failure: make([]RenewReport, 0),
		}
		today       = time.Now().UTC()
		productsMap = make(map[string]models.Product)
	)

	log.SetPrefix("[DraftFx] ")
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
		today = tmpDate
	}
	if req.DryRun != nil {
		dryRun = *req.DryRun
	}
	collectionPrefix = req.CollectionPrefix

	saveFn := func(p models.Policy, trs []models.Transaction) error {
		data := createDraftSaveBatch(p, trs)

		if !dryRun {
			return saveToDatabases(data)
		}

		return nil
	}

	policyType, quoteType, err := getQueryParameters(r)
	if err != nil {
		return "", nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	productsMap = getProductsMapByPolicyType(policyType, quoteType)
	log.Printf("products: %s", strings.Join(lib.GetMapKeys(productsMap), ", "))

	policies, err := getPolicies(req.PolicyUid, policyType, quoteType, productsMap, today)
	if err != nil {
		return "", nil, err
	}
	log.Printf("found %02d policies", len(policies))

	ch := make(chan RenewReport, len(policies))

	for _, policy := range policies {
		wg.Add(1)
		key := fmt.Sprintf("%s-%s", policy.Name, policy.ProductVersion)
		go draft(policy, productsMap[key], ch, wg, saveFn)
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
	if err != nil {
		return "", nil, err
	}

	if !dryRun {
		filename := fmt.Sprintf("renew/promote/report-%s-%d.json", today.Format(time.DateOnly), time.Now().Unix())
		if _, err = lib.PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filename, rawResp); err != nil {
			return "", nil, err
		}
	}

	sendReportMail(today, resp, true)

	return string(rawResp), resp, err
}

func getQueryParameters(r *http.Request) (policyType, quoteType string, err error) {
	policyType = r.URL.Query().Get("policyType")
	if policyType == "" {
		log.Printf("no policyType specified")
		return "", "", errors.New("no policyType specified")
	}

	quoteType = r.URL.Query().Get("quoteType")
	if quoteType == "" {
		log.Printf("no quoteType specified")
		return "", "", errors.New("no quoteType specified")
	}
	return policyType, quoteType, nil
}

func getPolicies(policyUid, policyType, quoteType string, products map[string]models.Product, today time.Time) ([]models.Policy, error) {
	var (
		err      error
		query    bytes.Buffer
		params   = make(map[string]interface{})
		policies []models.Policy
	)

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` p WHERE ", lib.WoptaDataset, lib.PoliciesViewCollection))

	if policyUid != "" {
		query.WriteString(" uid = @policyUid ")
		params["policyUid"] = policyUid

	} else if len(products) > 0 {
		tmpProducts := lib.GetMapValues(products)
		params["isRenewable"] = true
		params["policyType"] = policyType
		params["quoteType"] = quoteType
		params["isPay"] = true
		params["isDeleted"] = false

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
			targetDate := today.AddDate(0, 0, product.RenewOffset)

			productNameKey := fmt.Sprintf("%s%sProductName", product.Name, product.Version)
			productVersionKey := fmt.Sprintf("%s%sProductVersion", product.Name, product.Version)
			targetMonthKey := fmt.Sprintf("%s%sMonth", product.Name, product.Version)
			targetDayKey := fmt.Sprintf("%s%sDay", product.Name, product.Version)
			targetDateKey := fmt.Sprintf("%s%sDate", product.Name, product.Version)
			params[productNameKey] = product.Name
			params[productVersionKey] = product.Version
			params[targetMonthKey] = int64(targetDate.Month())
			params[targetDayKey] = int64(targetDate.Day())
			params[targetDateKey] = lib.GetBigQueryNullDateTime(targetDate)
			query.WriteString("(name = @" + productNameKey)
			query.WriteString(" AND productVersion = @" + productVersionKey)
			query.WriteString(" AND endDate > @" + targetDateKey)
			query.WriteString(" AND EXTRACT(MONTH FROM startDate) = @" + targetMonthKey)
			query.WriteString(" AND EXTRACT(DAY FROM startDate) = @" + targetDayKey + ")")
		}
		query.WriteString(") AND ")
		query.WriteString(fmt.Sprintf("(EXISTS(SELECT uid FROM `%s.%s` "+
			"WHERE uid = p.uid AND annuity = p.annuity + 1 AND isDeleted = false)) = false",
			lib.WoptaDataset, lib.RenewPolicyViewCollection))
	}

	log.Printf("query: %s", query.String())
	log.Printf("params: %v", params)

	policies, err = lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		log.Printf("error getting policies: %v", err)
		return nil, err
	}

	policies = lib.SliceMap(policies, func(p models.Policy) models.Policy {
		var tmpPolicy models.Policy
		err = json.Unmarshal([]byte(p.Data), &tmpPolicy)
		return tmpPolicy
	})

	return policies, nil
}

func draft(policy models.Policy, product models.Product, ch chan<- RenewReport, wg *sync.WaitGroup, save func(models.Policy, []models.Transaction) error) {
	var (
		err          error
		r            RenewReport
		transactions []models.Transaction
	)

	defer func() {
		r.Policy = policy
		r.Transactions = transactions
		if err != nil {
			r.Error = err.Error()
		}
		ch <- r
		wg.Done()
	}()

	policy.Annuity = policy.Annuity + 1

	err = calculatePricesByGuarantees(&policy)
	if err != nil {
		return
	}

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

	err = save(policy, transactions)
	if err != nil {
		return
	}
}

func calculatePricesByGuarantees(policy *models.Policy) error {
	var priceGross, priceNett, taxAmount, priceGrossMonthly, priceNettMonthly, taxAmountMonthly float64

	if policy.Name != models.LifeProduct {
		return errors.New("product not supported")
	}

	for index, guarantee := range policy.Assets[0].Guarantees {
		if policy.Annuity > guarantee.Value.Duration.Year || guarantee.IsDeleted {
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

	return nil
}
