package accounting

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type GetNetworkTransactionsResponse struct {
	NetworkTransactions []models.NetworkTransaction `json:"networkTransactions"`
}

func getNetworkTransactionsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var response GetNetworkTransactionsResponse

	log.AddPrefix("GetNetworkTransactionsFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	transactionUid := chi.URLParam(r, "transactionUid")

	log.Printf("transactionUid %s", transactionUid)

	netTranscations := transaction.GetNetworkTransactionsByTransactionUid(transactionUid)
	if len(netTranscations) == 0 {
		return "", "", fmt.Errorf("no network transactions found for transaction %s", transactionUid)
	}

	response.NetworkTransactions = netTranscations
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseByte), response, nil
}
