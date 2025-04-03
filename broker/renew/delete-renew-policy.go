package renew

import (
	"errors"
	"github.com/wopta/goworkspace/lib/log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/fabrick"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"
)

func DeleteRenewPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.AddPrefix("DeleteRenewPolicyFx")
	log.AddPrefix("[DeleteRenewPolicyFx] ")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	uid := chi.URLParam(r, "uid")

	if policy, err = plcRenew.GetRenewPolicyByUid(uid); err != nil {
		log.ErrorF("error getting renew policy %v", err)
		return "", nil, err
	}

	if policy.IsPay {
		log.Printf("cannot delete renew policy %s already paid", policy.Uid)
		return "", nil, errors.New("cannot delete renew policy")
	}

	transactions, err := trxRenew.GetRenewActiveTransactionsByPolicyUid(policy.Uid, policy.Annuity)
	if err != nil {
		log.ErrorF("error getting renew transactions %v", err)
		return "", nil, err
	}

	err = providerDeleteTransactions(policy.Payment, transactions)
	if err != nil {
		log.ErrorF("error deleting transaction on fabrick system %v", err)
		return "", nil, err
	}

	deletedPolicy := deleteRenewPolicy(policy)
	deletedTransactions := deleteRenewTransactions(transactions)

	batchData := createBatch(deletedPolicy, deletedTransactions)

	err = saveToDatabases(batchData)
	if err != nil {
		log.ErrorF("error saving batch to DB %v", err)
		return "", nil, err
	}

	return "{}", nil, nil
}

func providerDeleteTransactions(providerName string, transactions []models.Transaction) error {
	if providerName != models.FabrickPaymentProvider {
		return nil
	}

	for _, trx := range transactions {
		err := fabrick.FabrickExpireBill(trx.ProviderId)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteRenewPolicy(p models.Policy) models.Policy {
	p.IsDeleted = true
	p.Status = models.PolicyStatusDeleted
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.DeleteDesc = "Cancellata per annullo quietanza"
	p.Updated = time.Now().UTC()
	return p
}

func deleteRenewTransactions(transactions []models.Transaction) []models.Transaction {
	deletedTransactions := make([]models.Transaction, 0)
	for _, tr := range transactions {
		tr.IsDelete = true
		tr.Status = models.TransactionStatusDeleted
		tr.StatusHistory = append(tr.StatusHistory, tr.Status)
		tr.UpdateDate = time.Now().UTC()
		tr.PaymentNote = "Cancellata per annullo quietanza"
		deletedTransactions = append(deletedTransactions, tr)
	}
	return deletedTransactions
}

func createBatch(policy models.Policy, transactions []models.Transaction) map[string]map[string]interface{} {
	var (
		polCollection = lib.RenewPolicyCollection
		trsCollection = lib.RenewTransactionCollection
	)

	policy.Updated = time.Now().UTC()
	policy.BigQueryParse()
	batch := map[string]map[string]interface{}{
		polCollection: {
			policy.Uid: policy,
		},
		trsCollection: make(map[string]interface{}),
	}

	for idx, tr := range transactions {
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		batch[trsCollection][tr.Uid] = tr
		transactions[idx] = tr
	}

	return batch
}

func saveToDatabases(data map[string]map[string]interface{}) error {
	err := lib.SetBatchFirestoreErr(data)
	if err != nil {
		return err
	}

	for collection, values := range data {
		dataToSave := make([]interface{}, 0)
		for _, value := range values {
			dataToSave = append(dataToSave, value)
		}
		err = lib.InsertRowsBigQuery(models.WoptaDataset, collection, dataToSave)
		if err != nil {
			return err
		}
	}

	return nil
}
