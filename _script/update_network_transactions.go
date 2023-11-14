package _script

import (
	"fmt"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

func UpdateNetworkTransactions() {
	var (
		netTransactions []models.NetworkTransaction
		transaction     *models.Transaction
		err             error
		originalAmount  float64
		origin          = ""
		modifiedCounter = make([]string, 0)
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
		fmt.Printf("[UpdateNetworkTransactions] error getting network transactions: %s", err.Error())
		return
	}
	fmt.Printf("[UpdateNetworkTransactions] found %d netTransactions\n", len(netTransactions))
	// loop nt
	for _, nt := range netTransactions {
		// for each nt get its parent transaction (t)
		transaction = tr.GetTransactionByUid(nt.TransactionUid, origin)
		// update the nt.Amount and nt.AmountNet with t.Amount - nt.Amount
		if transaction == nil {
			fmt.Printf("[UpdateNetworkTransactions] error getting transaction '%s': %s", nt.TransactionUid, err.Error())
			return
		}

		policy := plc.GetPolicyByUid(transaction.PolicyUid, origin)
		mgaProduct := product.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
		commissionMga := product.GetCommissionByProduct(&policy, mgaProduct, false)

		if commissionMga != nt.Amount {
			fmt.Printf("[UpdateNetworkTransactions] netTransaction '%s' with amount '%f' already modified\n", nt.Uid, nt.Amount)
			continue
		}

		originalAmount = nt.Amount
		nt.Amount = lib.RoundFloat(transaction.Amount-nt.Amount, 2)
		nt.AmountNet = nt.Amount

		// save to bigquery
		// TODO: remember to manually allow for the modification of amount and amountNet fields
		err = nt.SaveBigQuery()
		if err != nil {
			fmt.Printf("[UpdateNetworkTransactions] error updating network transaction '%s': %s\n", nt.Uid, err.Error())
			break
		}

		modifiedCounter = append(modifiedCounter, nt.Uid)
		fmt.Printf("[UpdateNetworkTransactions] netTransaction '%s' original amount '%f' modified amount '%f'\n", nt.Uid, originalAmount, nt.Amount)
	}
	fmt.Printf("[UpdateNetworkTransactions] modified network transactions %s\n", modifiedCounter)
	fmt.Println("[UpdateNetworkTransactions] script done")
}
