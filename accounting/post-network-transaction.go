package accounting

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func CreateNetworkTransactionFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.AddPrefix("CreateNetworkTransactionFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	transactionUid := chi.URLParam(r, "transactionUid")
	log.Printf("transactionUid %s", transactionUid)

	transaction := tr.GetTransactionByUid(transactionUid, origin)
	if transaction == nil {
		log.ErrorF("could not retrieve transaction")
		return "", "", fmt.Errorf("error transaction %s not found", transactionUid)
	}

	err := CreateNetworkTransaction(transaction, origin)
	if err != nil {
		log.ErrorF("error creating network transactions: %s", err.Error())
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", "", err
}

func CreateNetworkTransaction(transaction *models.Transaction, origin string) error {
	policy := plc.GetPolicyByUid(transaction.PolicyUid, origin)
	producerNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	return tr.CreateNetworkTransactions(&policy, transaction, producerNode, mgaProduct)
}
