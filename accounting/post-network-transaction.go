package accounting

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

func CreateNetworkTransactionFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.SetPrefix("[CreateNetworkTransactionFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	transactionUid := chi.URLParam(r, "transactionUid")
	log.Printf("transactionUid %s", transactionUid)

	transaction := tr.GetTransactionByUid(transactionUid, origin)
	if transaction == nil {
		log.Println("could not retrieve transaction")
		return "", "", fmt.Errorf("error transaction %s not found", transactionUid)
	}

	err := CreateNetworkTransaction(transaction, origin)
	if err != nil {
		log.Printf("error creating network transactions: %s", err.Error())
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
