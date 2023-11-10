package _script

import (
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	tr "github.com/wopta/goworkspace/transaction"
)

func UpdateNetworkTransactions() {
	var (
		netTransactions []models.NetworkTransaction
		transaction     *models.Transaction
		err             error
		originalAmount  float64
	)

	// get all network transactions of RemittanceCompany
	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE paymentType = '%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		models.PaymentTypeRemittanceCompany,
	)
	netTransactions, err = lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil {
		log.Printf("[UpdateNetworkTransactions] error getting network transactions: %s", err.Error())
		return
	}
	log.Printf("[UpdateNetworkTransactions] found %d netTransactions", len(netTransactions))
	// loop nt
	for _, nt := range netTransactions {

		if nt.Uid == "ef8a69c7-820e-450b-8691-708e04b17b9d" {
			originalAmount = nt.Amount
			nt.Amount = lib.RoundFloat(nt.Amount, 2)
			nt.AmountNet = nt.Amount
		} else {
			// for each nt get its parent transaction (t)
			transaction = tr.GetTransactionByUid(nt.TransactionUid, "")
			// update the nt.Amount and nt.AmountNet with t.Amount - nt.Amount
			if transaction == nil {
				log.Printf("[UpdateNetworkTransactions] error getting transaction '%s': %s", nt.TransactionUid, err.Error())
				return
			}
			// log.Printf("[UpdateNetworkTransactions] transaction amount: '%f'", transaction.Amount)
			// log.Printf("[UpdateNetworkTransactions] original amount: '%f' amountNet: '%f'", nt.Amount, nt.AmountNet)
			originalAmount = nt.Amount
			nt.Amount = lib.RoundFloat(transaction.Amount-nt.Amount, 2)
			nt.AmountNet = nt.Amount
			// log.Printf("[UpdateNetworkTransactions] modified amount: '%f' amountNet: '%f'", nt.Amount, nt.AmountNet)
		}

		// save to bigquery
		// TODO: remember to manually allow for the modification of amount and amountNet fields
		nt.SaveBigQuery()
		log.Printf("[UpdateNetworkTransactions] netTransaction '%s' original amount '%f' modified amount '%f'", nt.Uid, originalAmount, nt.Amount)
	}
	log.Println("[UpdateNetworkTransactions] script done")
}
