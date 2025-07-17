package _script

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

func getAllPolicies() ([]models.Policy, error) {
	var policies = make([]models.Policy, 0)
	docIterator := lib.OrderFirestore(lib.PolicyCollection, "uid", firestore.Asc)

	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.Printf("error getting polcies from Firestore: %s", err.Error())
		return policies, err
	}

	for _, snapshot := range snapshots {
		var policy models.Policy
		err = snapshot.DataTo(&policy)
		if err != nil {
			log.Printf("error parsing policy %s: %s", snapshot.Ref.ID, err.Error())
		} else {
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

func PolicyTransactionsUpdate() {
	var (
		wg          *sync.WaitGroup = new(sync.WaitGroup)
		productsMap                 = make(map[string]models.Product)
	)

	// ATTENTION: the function needs an manual override at product/get-product.go:110
	// otherwise if will filter out old versions of a given product and won't updated
	// the needed policies for life/v1 renewal. Remove the isActive check
	productsInfo := product.GetAllProductsByChannel(models.MgaChannel)
	for _, pr := range productsInfo {
		prd := product.GetProductV2(pr.Name, pr.Version, models.MgaChannel, nil, nil)
		productsMap[fmt.Sprintf("%s-%s", prd.Name, prd.Version)] = *prd
	}

	policies, err := getAllPolicies()
	if err != nil {
		return
	}

	startDate := time.Now().UTC()
	log.Printf("Started at %s", startDate.String())

	// TODO: split policies in batch es.: 1000
	for _, p := range policies {
		wg.Add(1)
		go func(p models.Policy) {
			defer wg.Done()

			m := map[string]map[string]interface{}{
				models.PolicyCollection:       make(map[string]interface{}),
				models.TransactionsCollection: make(map[string]interface{}),
			}
			transactionsList := make([]models.Transaction, 0)

			productIdentifier := fmt.Sprintf("%s-%s", p.Name, p.ProductVersion)

			if strings.EqualFold(p.Name, models.LifeProduct) {
				p.OfferlName = "default"
			}

			if strings.EqualFold(p.PaymentSplit, string(models.PaySplitYear)) {
				p.PaymentSplit = string(models.PaySplitYearly)
			}

			if strings.EqualFold(p.PaymentSplit, string(models.PaySplitYearly)) && p.PaymentMode == "" {
				p.PaymentMode = models.PaymentModeSingle
			}

			if strings.EqualFold(p.PaymentSplit, string(models.PaySplitMonthly)) && p.PaymentMode == "" {
				p.PaymentMode = models.PaymentModeRecurrent
			}

			if p.Payment == "" {
				p.Payment = models.FabrickPaymentProvider
			}

			if p.Channel == "" {
				p.Channel = models.ECommerceChannel
			}

			p.Annuity = 0
			p.IsAutoRenew = productsMap[productIdentifier].IsAutoRenew
			p.IsRenewable = productsMap[productIdentifier].IsRenewable
			p.QuoteType = productsMap[productIdentifier].QuoteType
			p.PolicyType = productsMap[productIdentifier].PolicyType
			p.Updated = time.Now().UTC()
			m[models.PolicyCollection][p.Uid] = p

			transactions := transaction.GetPolicyTransactions(p.Uid)
			for _, t := range transactions {
				t.Annuity = 0
				t.UpdateDate = time.Now().UTC()
				t.BigQueryParse()
				m[models.TransactionsCollection][t.Uid] = t
				transactionsList = append(transactionsList, t)
			}

			err = lib.SetBatchFirestoreErr(m)
			if err != nil {
				log.Printf("error saving policy and transactiosn into firestore: %s", err.Error())
				return
			}

			p.BigquerySave()

			err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, transactionsList)
			if err != nil {
				log.Println("error saving transactions into BigQuery", err)
			}
			log.Printf("Updated data for policy %s", p.Uid)
		}(p)
	}

	wg.Wait()

	endDate := time.Now().UTC()
	log.Printf("End at %s duration %s", endDate.String(), endDate.Sub(startDate).String())
	log.Printf("N policies %d", len(policies))
}
