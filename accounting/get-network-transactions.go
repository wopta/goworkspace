package accounting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type GetNetworkTransactionsResponse struct {
	NetworkTransactions []models.NetworkTransaction `json:"networkTransactions"`
}

func GetNetworkTransactionsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.Println("[GetNetworkTransactionsFx] Handler start ---------------------")

	var response GetNetworkTransactionsResponse

	transactionUid := r.Header.Get("transactionUid")

	log.Printf("[GetNetworkTransactionsFx] transactionUid %s", transactionUid)

	netTranscations := GetNetworkTransactionsByTransactionUid(transactionUid)
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

func GetNetworkTransactionsByTransactionUid(transactionUid string) []models.NetworkTransaction {
	log.Printf("[GetNetworkTransactionsByTransactionUid] transactionUid %s", transactionUid)

	var (
		netTransactions []models.NetworkTransaction
		err             error
	)

	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE transactionUid='%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		transactionUid,
	)

	netTransactions, err = lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil {
		log.Printf("[GetNetworkTransactionsByTransactionUid] error getting network transactions: %s", err.Error())
	}

	log.Printf("[GetNetworkTransactionsByTransactionUid] found %d network transactions", len(netTransactions))
	return netTransactions
}
