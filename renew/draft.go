package renew

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/client"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type DraftReq struct {
	PolicyUid        string `json:"policyUid"`
	Date             string `json:"date"`
	DryRun           *bool  `json:"dryRun"`
	CollectionPrefix string `json:"collectionPrefix"`
	SendMail         *bool  `json:"sendMail"`
}

type NodeFlowRelation struct {
	Node models.NetworkNode
	Flow string
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
		productsMap = make(map[string]map[string]models.Product)
		sendMail    = false
	)

	log.AddPrefix("[DraftFx] ")
	defer func() {
		collectionPrefix = ""
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
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
			log.ErrorF("error parsing request date: %s", err.Error())
			return "", nil, err
		}
		today = tmpDate
	}
	if req.DryRun != nil {
		dryRun = *req.DryRun
	}
	if req.SendMail != nil {
		sendMail = *req.SendMail
	}
	collectionPrefix = req.CollectionPrefix

	log.Printf("running pipeline with set config. Today as: %v, DryRun: %v, SendMail: %v", today, dryRun, sendMail)

	policyType, quoteType, err := getQueryParameters(r)
	if err != nil {
		return "", nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	productsMap = getProducts(policyType, quoteType)
	log.Printf("products: %s", strings.Join(lib.GetMapKeys(productsMap), ", "))

	policies, err := getPolicies(req.PolicyUid, policyType, quoteType, productsMap[models.MgaChannel], today)
	if err != nil {
		return "", nil, err
	}
	log.Printf("found %02d policies", len(policies))

	saveFn := func(p models.Policy, trs []models.Transaction, hasMandate bool) error {
		data := createDraftSaveBatch(p, trs)

		if !dryRun {
			if err := saveToDatabases(data); err != nil {
				return err
			}
			return nil
		}

		log.Println("dryRun active - not saving to DB and not sending mail")

		return nil
	}

	ch := make(chan RenewReport, len(policies))

	for _, policy := range policies {
		wg.Add(1)
		key := fmt.Sprintf("%s-%s", policy.Name, policy.ProductVersion)
		go draft(policy, productsMap[policy.Channel][key], ch, wg, saveFn)
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
		filename := fmt.Sprintf("renew/draft/report-draft-%s-%d.json", today.Format(time.DateOnly), time.Now().Unix())
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
		log.ErrorF("no quoteType specified")
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
		params["policyUid"] = policyUid

		query.WriteString(" uid = @policyUid ")
		query.WriteString(" AND ")
	} else if len(products) > 0 {
		tmpProducts := lib.GetMapValues(products)
		params["isRenewable"] = true
		params["policyType"] = policyType
		params["quoteType"] = quoteType
		params["isPay"] = true
		params["isDeleted"] = false
		params["year"] = today.Year()

		query.WriteString("isRenewable = @isRenewable")
		query.WriteString(" AND policyType = @policyType")
		query.WriteString(" AND quoteType = @quoteType")
		query.WriteString(" AND isPay = @isPay")
		query.WriteString(" AND isDeleted = @isDeleted")
		query.WriteString(" AND EXTRACT(YEAR FROM startDate) <= @year")
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
	}

	query.WriteString(fmt.Sprintf("(EXISTS(SELECT uid FROM `%s.%s` "+
		"WHERE uid = p.uid AND annuity = p.annuity + 1 AND isDeleted = false)) = false",
		lib.WoptaDataset, lib.RenewPolicyViewCollection))

	log.Printf("query: %s", query.String())
	log.Printf("params: %v", params)

	policies, err = lib.QueryParametrizedRowsBigQuery[models.Policy](query.String(), params)
	if err != nil {
		log.ErrorF("error getting policies: %v", err)
		return nil, err
	}

	policies = lib.SliceMap(policies, func(p models.Policy) models.Policy {
		var tmpPolicy models.Policy
		err = json.Unmarshal([]byte(p.Data), &tmpPolicy)
		return tmpPolicy
	})

	return policies, nil
}

func draft(policy models.Policy, product models.Product, ch chan<- RenewReport, wg *sync.WaitGroup, save func(models.Policy, []models.Transaction, bool) error) {
	var (
		err          error
		r            RenewReport
		transactions []models.Transaction
		customerId   string
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
	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusDraftRenew, models.PolicyStatusToPay)

	transactions = transaction.CreateTransactions(policy, product, func() string {
		return lib.NewDoc(models.TransactionsCollection)
	})

	if policy.Payment == models.FabrickPaymentProvider {
		var isTransactionPaid bool = true
		trs := transaction.GetPolicyValidTransactions(policy.Uid, &isTransactionPaid)
		if len(trs) > 0 {
			customerId = trs[len(trs)-1].UserToken
		}
	}

	// TODO: code smell
	if policy.Channel == lib.ECommerceChannel {
		policy.PaymentMode = models.PaymentModeRecurrent
	}

	client := client.NewClient(policy.Payment, policy, product, transactions, customerId != "", customerId)
	payUrl, hasMandate, transactions, err := client.Renew()
	if err != nil {
		return
	}

	if payUrl != "" {
		policy.PayUrl = payUrl
	}
	policy.Updated = time.Now().UTC()
	policy.HasMandate = hasMandate
	policy.IsRenew = true

	err = save(policy, transactions, hasMandate)
	if err != nil {
		return
	}
}

func calculatePricesByGuarantees(policy *models.Policy) error {
	var priceGross, priceNett, taxAmount, priceGrossMonthly, priceNettMonthly, taxAmountMonthly float64

	if policy.Name != models.LifeProduct {
		return errors.New("product not supported")
	}

	if policy.OffersPrices[policy.OfferlName] == nil {
		return errors.New("invalid offer name")
	}

	if policy.OffersPrices[policy.OfferlName][policy.PaymentSplit] == nil {
		return errors.New("no offer found for payment split")
	}

	for index, guarantee := range policy.Assets[0].Guarantees {
		if policy.Annuity >= guarantee.Value.Duration.Year || guarantee.IsDeleted {
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
	policy.PriceGross = lib.RoundFloat(priceGross, 2)
	policy.PriceNett = lib.RoundFloat(priceNett, 2)
	policy.TaxAmount = lib.RoundFloat(taxAmount, 2)
	policy.PriceGrossMonthly = lib.RoundFloat(priceGrossMonthly, 2)
	policy.PriceNettMonthly = lib.RoundFloat(priceNettMonthly, 2)
	policy.TaxAmountMonthly = lib.RoundFloat(taxAmountMonthly, 2)
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Tax = lib.RoundFloat(policy.TaxAmount, 2)
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Net = lib.RoundFloat(policy.PriceNett, 2)
	policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Gross = lib.RoundFloat(policy.PriceGross, 2)

	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Tax = policy.TaxAmountMonthly
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Net = policy.PriceNettMonthly
		policy.OffersPrices[policy.OfferlName][policy.PaymentSplit].Gross = policy.PriceGrossMonthly
	}

	return nil
}

func getPolicyMailDataMap(ps []models.Policy) map[string]NodeFlowRelation {
	var (
		policyNodeFlowMap = make(map[string]NodeFlowRelation)
		nodeMap           = make(map[string]models.NetworkNode)
		warrants          []models.Warrant
		warrantMap        = make(map[string]models.Warrant)
		err               error
	)

	if warrants, err = network.GetWarrants(); err != nil {
		log.ErrorF("error loading warrants: %s", err)
		return nil
	}

	for _, w := range warrants {
		warrantMap[w.Name] = w
	}

	for _, p := range ps {
		if p.ProducerUid == "" {
			continue
		}
		currentNode := nodeMap[p.ProducerUid]
		if _, ok := nodeMap[p.ProducerUid]; !ok {
			nn := network.GetNetworkNodeByUid(p.ProducerUid)
			if nn == nil {
				log.ErrorF("error loading networkNode: %s", p.ProducerUid)
				return nil
			}
			currentNode = *nn
		}
		w := warrantMap[currentNode.Warrant]
		flowName := w.GetFlowName(p.Name)

		policyNodeFlowMap[p.Uid] = NodeFlowRelation{
			Node: currentNode,
			Flow: flowName,
		}
	}

	return policyNodeFlowMap
}
