package _script

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
	"google.golang.org/api/iterator"
)

func UpdateTransactions() {
	var (
		trs                  = make([]models.Transaction, 0)
		tr                   models.Transaction
		modifiedTransactions = make([]string, 0)
	)
	log.AddPrefix("UpdateTransactions")
	defer log.PopPrefix()
	// GET all transactions from firestore
	ctx := context.Background()
	client, _ := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	iter := client.Collection(models.TransactionsCollection).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("ERROR iterator: %s", err.Error())
			return
		}
		if err := doc.DataTo(&tr); err != nil {
			fmt.Printf("ERROR datato: %s", err.Error())
			return
		}
		trs = append(trs, tr)
	}

	fmt.Printf("Found %d transactions/n", len(trs))

	for _, tr := range trs {
		// set updateDate to now
		tr.UpdateDate = time.Now().UTC()

		// save to firestore
		err := lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			fmt.Printf("ERROR saving to firestore: %s", err.Error())
			return
		}
		modifiedTransactions = append(modifiedTransactions, tr.Uid)
		// save to bigquery
		tr.BigQuerySave()
	}

	fmt.Printf("Modified %d transactions: %v/n", len(modifiedTransactions), modifiedTransactions)
	fmt.Println("done")
}

type QueryResult struct {
	Uid string `bigquery:"uid"`
}

func UpdateTransactionsCommission() {
	var modifiedTransactions = make([]string, 0)
	log.AddPrefix("UpdateNetworkTransactions")
	defer log.PopPrefix()
	query := fmt.Sprintf(
		"SELECT uid FROM `%s.%s` WHERE isDelete = false AND commissions = 0",
		models.WoptaDataset,
		models.TransactionsViewCollection,
	)
	transactionUids, err := lib.QueryRowsBigQuery[QueryResult](query)
	if err != nil {
		fmt.Printf("error getting network transactions: %s", err.Error())
		return
	}

	fmt.Printf("found %d transactions\n", len(transactionUids))

	if len(transactionUids) == 0 {
		fmt.Println("done")
		return
	}

	lifeMgaProduct := product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil, nil)
	gapMgaProduct := product.GetProductV2(models.GapProduct, models.ProductV1, models.MgaChannel, nil, nil)
	pmiMgaProduct := product.GetProductV2(models.PmiProduct, models.ProductV1, models.MgaChannel, nil, nil)
	personaMgaProduct := product.GetProductV2(models.PersonaProduct, models.ProductV1, models.MgaChannel, nil, nil)

	productMap := map[string]*models.Product{
		models.LifeProduct:    lifeMgaProduct,
		models.GapProduct:     gapMgaProduct,
		models.PmiProduct:     pmiMgaProduct,
		models.PersonaProduct: personaMgaProduct,
	}

	policyMap := make(map[string]*models.Policy)

	for _, t := range transactionUids {
		tr := transaction.GetTransactionByUid(t.Uid)
		var plc = policyMap[tr.PolicyUid]

		if plc == nil {
			temp, err := policy.GetPolicy(tr.PolicyUid)
			if err != nil {
				log.Error(err)
				return
			}
			policyMap[tr.PolicyUid] = plc
			plc = &temp
		}

		err := updateTransactionCommission(tr, plc, productMap[tr.PolicyName])
		if err != nil {
			return
		}
		modifiedTransactions = append(modifiedTransactions, tr.Uid)
	}

	fmt.Printf("modified %d transactions => %s\n", len(modifiedTransactions), modifiedTransactions)
	fmt.Println("done")
}

func updateTransactionCommission(tr *models.Transaction, policy *models.Policy, mgaProduct *models.Product) error {
	log.AddPrefix("updateTransactionCommission")
	defer log.PopPrefix()
	commissionMga := lib.RoundFloat(product.GetCommissionByProduct(policy, mgaProduct, false), 2)
	fmt.Printf("new commission %.2f\n", commissionMga)

	tr.Commissions = commissionMga
	tr.UpdateDate = time.Now().UTC()

	err := lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
	if err != nil {
		fmt.Printf("error saving transaction to firestore: %s", err.Error())
		return err
	}
	tr.BigQuerySave()

	return nil
}
