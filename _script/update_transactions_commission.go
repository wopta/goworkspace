package _script

import (
	"fmt"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
)

func UpdateTransactionsCommission() {
	var modifiedTransactions = make([]string, 0)

	// TODO: wait for the updated schema where we should have isDelete field and updateDate
	query := fmt.Sprintf(
		"SELECT uid FROM `%s.%s` WHERE commissions = 0",
		models.WoptaDataset,
		models.TransactionsViewCollection,
	)
	transactionUids, err := lib.QueryRowsBigQuery[string](query)
	if err != nil {
		fmt.Printf("[UpdateNetworkTransactions] error getting network transactions: %s", err.Error())
		return
	}

	fmt.Printf("[UpdateTransactionsCommission] found %d transactions\n", len(transactionUids))

	if len(transactionUids) == 0 {
		fmt.Println("[UpdateTransactionsCommission] done")
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

	for _, uid := range transactionUids {
		tr := transaction.GetTransactionByUid(uid, "")
		var plc = policyMap[tr.PolicyUid]

		if plc == nil {
			temp := policy.GetPolicyByUid(tr.PolicyUid, "")
			plc = &temp
			policyMap[tr.PolicyUid] = plc
		}
		fmt.Printf("[UpdateTransactionsCommission] transaction.PolicyUid - '%s' | policy.Uid - '%s'\n", tr.PolicyUid, plc.Uid)

		err := updateTransactionCommission(tr, plc, productMap[tr.PolicyName])
		if err != nil {
			return
		}
		modifiedTransactions = append(modifiedTransactions, tr.Uid)
	}

	fmt.Printf("[UpdateTransactionsCommission] modified %d transactions => %s\n", len(modifiedTransactions), modifiedTransactions)
	fmt.Println("[UpdateTransactionsCommission] done")
}

func updateTransactionCommission(tr *models.Transaction, policy *models.Policy, mgaProduct *models.Product) error {
	commissionMga := lib.RoundFloat(product.GetCommissionByProduct(policy, mgaProduct, false), 2)
	fmt.Printf("[updateTransactionCommission] new commission %.2f\n", commissionMga)

	tr.Commissions = commissionMga

	// err := lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
	// if err != nil {
	// 	log.Printf("[updateTransactionCommission] error saving transaction to firestore: %s", err.Error())
	// 	return err
	// }
	// tr.BigQuerySave("")

	return nil
}
