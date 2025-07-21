package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

func getAllPolicies(numPolicies int) ([]models.Policy, error) {
	var policies = make([]models.Policy, 0)
	docIterator := lib.OrderFirestore(lib.PolicyCollection, "uid", firestore.Asc)

	snapshots, err := docIterator.GetAll()
	if err != nil {
		log.ErrorF("error getting polcies from Firestore: %s", err.Error())
		return policies, err
	}

	for _, snapshot := range snapshots {
		var policy models.Policy
		err = snapshot.DataTo(&policy)
		if err != nil {
			log.ErrorF("error parsing policy %s: %s", snapshot.Ref.ID, err.Error())
		} else {
			policies = append(policies, policy)
		}
	}

	return policies[:numPolicies], nil
}

func policyTransactionsUpdate(request int) {
	var (
		wg          sync.WaitGroup
		productsMap = make(map[string]models.Product)
	)

	productsInfo := product.GetAllProductsByChannel(models.MgaChannel)
	for _, pr := range productsInfo {
		prd := product.GetProductV2(pr.Name, pr.Version, models.MgaChannel, nil, nil)
		productsMap[fmt.Sprintf("%s-%s", prd.Name, prd.Version)] = *prd
	}

	policies, err := getAllPolicies(request)
	if err != nil {
		return
	}

	startDate := time.Now().UTC()
	log.Printf("Started at %s", startDate.String())

	wg.Add(len(policies))

	for _, p := range policies {
		go func(p models.Policy) {
			defer wg.Done()

			m := map[string]map[string]interface{}{
				models.PolicyCollection:       make(map[string]interface{}),
				models.TransactionsCollection: make(map[string]interface{}),
			}
			transactionsList := make([]models.Transaction, 0)

			productIdentifier := fmt.Sprintf("%s-%s", p.Name, p.ProductVersion)

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
				t.BigPayDate = lib.GetBigQueryNullDateTime(t.PayDate)
				t.BigTransactionDate = lib.GetBigQueryNullDateTime(t.TransactionDate)
				t.BigCreationDate = civil.DateTimeOf(t.CreationDate)
				t.BigStatusHistory = strings.Join(t.StatusHistory, ",")
				t.BigUpdateDate = lib.GetBigQueryNullDateTime(t.UpdateDate)
				t.BigEffectiveDate = lib.GetBigQueryNullDateTime(t.EffectiveDate)
				m[models.TransactionsCollection][t.Uid] = t
				transactionsList = append(transactionsList, t)
			}

			err = setBatchFirestoreErr(m)
			if err != nil {
				log.ErrorF("error saving policy and transactiosn into firestore: %s", err.Error())
				return
			}

			p.BigquerySave()

			err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, transactionsList)
			if err != nil {
				log.ErrorF("error saving transactions into BigQuery %s", err)
			}
			log.Printf("Updated data for policy %s", p.Uid)
		}(p)
	}

	wg.Wait()

	endDate := time.Now().UTC()
	log.Printf("End at %s duration %s", endDate.String(), endDate.Sub(startDate).String())

}

func setBatchFirestoreErr[T any](operations map[string]map[string]T) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return err
	}

	bulk := client.BulkWriter(ctx)

	for collection, values := range operations {
		c := client.Collection(collection)

		for k, v := range values {
			col := c.Doc(k)
			_, err = bulk.Set(col, v)
			if err != nil {
				log.ErrorF("error batch firestore: %s", err.Error())
				return err
			}
		}
	}

	bulk.End()

	return nil
}
