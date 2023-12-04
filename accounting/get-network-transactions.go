package accounting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

type GetNetworkTransactionsResponse struct {
	NetworkTransactions []models.NetworkTransaction `json:"networkTransactions"`
}

func GetNetworkTransactionsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.Println("[GetNetworkTransactionsFx] Handler start ---------------------")

	var response GetNetworkTransactionsResponse

	transactionUid := r.Header.Get("transactionUid")

	log.Printf("[GetNetworkTransactionsFx] transactionUid %s", transactionUid)

	netTranscations := transaction.GetNetworkTransactionsByTransactionUid(transactionUid)
	if len(netTranscations) == 0 {
		return "", "", fmt.Errorf("no network transactions found for transaction %s", transactionUid)
	}

	response.NetworkTransactions = netTranscations
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.Printf("[GetNetworkTransactionsFx] error marshaling response: %s", err.Error())
		return "", "", err
	}

	return string(responseByte), response, nil
}
