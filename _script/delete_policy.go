package _script

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	tr "github.com/wopta/goworkspace/transaction"
	"log"
	"time"
)

func DeletePolicy(policyUid string) {
	var (
		err          error
		policy       models.Policy
		transactions []models.Transaction
	)

	policy, err = plc.GetPolicy(policyUid, "")
	if err != nil {
		log.Printf("error retrieving policy from Firestore: %s", err.Error())
		return
	}

	policy.IsDeleted = true
	policy.Status = models.PolicyStatusDeleted
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()

	// save policy in Firestore
	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		log.Printf("error saving policy in Firestore: %s", err.Error())
		return
	}

	// save policy in BigQuery
	policy.BigquerySave("")

	transactions = tr.GetPolicyTransactions("", policyUid)
	for index, _ := range transactions {
		transactions[index].IsDelete = true
		transactions[index].Status = models.TransactionStatusDeleted
		transactions[index].StatusHistory = append(transactions[index].StatusHistory, models.TransactionStatusDeleted)
		transactions[index].UpdateDate = time.Now().UTC()

		// save transaction in Firestore
		err = lib.SetFirestoreErr(models.TransactionsCollection, transactions[index].Uid, transactions[index])
		if err != nil {
			log.Printf("error saving transaction in Firestore: %s", err.Error())
			return
		}

		// save transaction in BigQuery
		transactions[index].BigQuerySave("")
	}

}

func SaveTransactionBigQuery(policyUid string) {
	transactions := tr.GetPolicyTransactions("", policyUid)
	for index, _ := range transactions {
		// save transaction in BigQuery
		transactions[index].BigQuerySave("")
	}
}
