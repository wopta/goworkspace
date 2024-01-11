package _script

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
	"log"
	"os"
	"strings"
	"time"
)

func RestoreTransactions(inputStatus string) {
	outputTransactions := make(map[string][]models.Transaction)

	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE status = '%s'",
		models.WoptaDataset,
		models.TransactionsCollection,
		inputStatus,
	)
	transactions, err := lib.QueryRowsBigQuery[models.Transaction](query)
	if err != nil {
		log.Printf("error fetching transaction from BigQuery: %s", err.Error())
		return
	}

	if len(transactions) == 0 {
		fmt.Println("no transactions found")
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

	for _, tr := range transactions {
		log.Printf("Transaction Uid: %s Policy Uid: %s", tr.Uid, tr.PolicyUid)

		if policyMap[tr.PolicyUid] == nil {
			policyMap[tr.PolicyUid] = new(models.Policy)
			*policyMap[tr.PolicyUid] = policy.GetPolicyByUid(tr.PolicyUid, "")
		}

		plc := *policyMap[tr.PolicyUid]

		tr.ProviderName = plc.Payment

		_updateTransactionCommission(&tr, &plc, productMap[tr.PolicyName])

		_reverseBigQueryFields(&tr)

		if outputTransactions[tr.PolicyUid] == nil {
			outputTransactions[tr.PolicyUid] = make([]models.Transaction, 0)
		}

		outputTransactions[tr.PolicyUid] = append(outputTransactions[tr.PolicyUid], tr)

		err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			fmt.Printf("[_updateTransactionCommission] error saving transaction to firestore: %s", err.Error())
			return
		}
		tr.BigQuerySave("")
	}

	rawOutput, err := json.Marshal(outputTransactions)
	if err != nil {
		log.Printf(err.Error())
	}
	err = os.WriteFile("./transaction_restore_output_"+inputStatus+".json", rawOutput, 777)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

}

func _updateTransactionCommission(tr *models.Transaction, policy *models.Policy, mgaProduct *models.Product) {
	commissionMga := lib.RoundFloat(product.GetCommissionByProduct(policy, mgaProduct, false), 2)
	fmt.Printf("[_updateTransactionCommission] new commission %.2f\n", commissionMga)

	tr.Commissions = commissionMga
	tr.UpdateDate = time.Now().UTC()
}

func _reverseBigQueryFields(tr *models.Transaction) {
	tr.CreationDate = tr.BigCreationDate.In(time.UTC)
	tr.StatusHistory = strings.Split(tr.BigStatusHistory, ",")
	if tr.BigPayDate.Valid {
		tr.PayDate = tr.BigPayDate.DateTime.In(time.UTC)
	}
	if tr.BigTransactionDate.Valid {
		tr.TransactionDate = tr.BigTransactionDate.DateTime.In(time.UTC)
	}
}
