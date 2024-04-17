package accounting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

type GetNetworkTransactionsResponse struct {
	NetworkTransactions []models.NetworkTransaction `json:"networkTransactions"`
}

func GetNetworkTransactionsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var response GetNetworkTransactionsResponse

	log.SetPrefix("[GetNetworkTransactionsFx] ")
	defer log.SetPrefix("")

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
		log.Printf("error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseByte), response, nil
}
