package transaction

import (
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetNetworkTransactionByUid(uid string) *models.NetworkTransaction {
	log.Printf("[GetNetworkTransactionByUid] uid %s", uid)

	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE uid='%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		uid,
	)

	netTransactions, err := lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil || len(netTransactions) == 0 {
		log.Printf("[GetNetworkTransactionByUid] error getting network transactions: %s", err.Error())
		return nil
	}

	return &netTransactions[0]
}

func GetNetworkTransactionsByTransactionUid(transactionUid string) []models.NetworkTransaction {
	log.Printf("[GetNetworkTransactionsByTransactionUid] transactionUid %s", transactionUid)

	var (
		netTransactions = make([]models.NetworkTransaction, 0)
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
