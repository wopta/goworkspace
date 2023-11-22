package _script

import (
	"fmt"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
)

func UpdateCompanyNetworkTransactions() {
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

func UpdateAreaManagerName() {
	var (
		netTransactions []models.NetworkTransaction
		err             error
		originalName    string
		modifiedCounter = make([]string, 0)
	)

	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE networkNodeType = '%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		models.AreaManagerNetworkNodeType,
	)
	netTransactions, err = lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil {
		fmt.Printf("[UpdateAreaManagerName] error getting network transactions: %s", err.Error())
		return
	}
	fmt.Printf("[UpdateAreaManagerName] found %d netTransactions\n", len(netTransactions))

	for _, nt := range netTransactions {
		nn := network.GetNetworkNodeByUid(nt.NetworkNodeUid)

		originalName = nt.Name
		nodeName := nn.GetName()

		if strings.HasSuffix(originalName, nodeName) {
			fmt.Printf("[UpdateAreaManagerName] netTransaction '%s' with name '%s' already contains node name '%s'\n", nt.Uid, originalName, nodeName)
			continue
		}

		nt.Name = nt.Name + nn.GetName()

		err = nt.SaveBigQuery()
		if err != nil {
			fmt.Printf("[UpdateAreaManagerName] error updating network transaction '%s': %s\n", nt.Uid, err.Error())
			break
		}

		modifiedCounter = append(modifiedCounter, nt.Uid)
		fmt.Printf("[UpdateAreaManagerName] netTransaction '%s' original name '%s' modified name '%s'\n", nt.Uid, originalName, nt.Name)
	}

	fmt.Printf("[UpdateAreaManagerName] modified network transactions %s\n", modifiedCounter)
	fmt.Println("[UpdateAreaManagerName] script done")
}
