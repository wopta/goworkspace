package _script

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mohae/deepcopy"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func UpdateCompanyNetworkTransactions() {
	var (
		netTransactions []models.NetworkTransaction
		transaction     *models.Transaction
		err             error
		originalAmount  float64
		modifiedCounter = make([]string, 0)
	)
	log.AddPrefix("UpdateNetworkTransactions")
	defer log.PopPrefix()
	// get all network transactions of RemittanceCompany
	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE paymentType = '%s'",
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		models.PaymentTypeRemittanceCompany,
	)
	netTransactions, err = lib.QueryRowsBigQuery[models.NetworkTransaction](query)
	if err != nil {
		fmt.Printf("error getting network transactions: %s", err.Error())
		return
	}
	fmt.Printf("found %d netTransactions\n", len(netTransactions))
	// loop nt
	for _, nt := range netTransactions {
		// for each nt get its parent transaction (t)
		transaction = tr.GetTransactionByUid(nt.TransactionUid)
		// update the nt.Amount and nt.AmountNet with t.Amount - nt.Amount
		if transaction == nil {
			fmt.Printf("error getting transaction '%s': %s", nt.TransactionUid, err.Error())
			return
		}

		policy, err := plc.GetPolicy(transaction.PolicyUid)
		if err != nil {
			log.Error(err)
			return
		}
		mgaProduct := product.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
		commissionMga := product.GetCommissionByProduct(&policy, mgaProduct, false)

		if commissionMga != nt.Amount {
			fmt.Printf("netTransaction '%s' with amount '%f' already modified\n", nt.Uid, nt.Amount)
			continue
		}

		originalAmount = nt.Amount
		nt.Amount = lib.RoundFloat(transaction.Amount-nt.Amount, 2)
		nt.AmountNet = nt.Amount

		// save to bigquery
		// TODO: remember to manually allow for the modification of amount and amountNet fields
		err = nt.SaveBigQuery()
		if err != nil {
			fmt.Printf("error updating network transaction '%s': %s\n", nt.Uid, err.Error())
			break
		}

		modifiedCounter = append(modifiedCounter, nt.Uid)
		fmt.Printf("netTransaction '%s' original amount '%f' modified amount '%f'\n", nt.Uid, originalAmount, nt.Amount)
	}
	fmt.Printf("modified %d network transactions %s\n", len(modifiedCounter), modifiedCounter)
	fmt.Println("script done")
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
		fmt.Printf("error getting network transactions: %s", err.Error())
		return
	}
	fmt.Printf("found %d netTransactions\n", len(netTransactions))

	for _, nt := range netTransactions {
		nn := network.GetNetworkNodeByUid(nt.NetworkNodeUid)

		originalName = nt.Name
		nodeName := nn.GetName()

		if strings.HasSuffix(strings.ToLower(originalName), strings.ToLower(nodeName)) {
			fmt.Printf("netTransaction '%s' with name '%s' already contains node name '%s'\n", nt.Uid, originalName, nodeName)
			continue
		}

		nt.Name = strings.ToUpper(nt.Name + nn.GetName())

		// TODO: remember to manually allow for the modification of name field
		err = nt.SaveBigQuery()
		if err != nil {
			fmt.Printf("error updating network transaction '%s': %s\n", nt.Uid, err.Error())
			break
		}

		modifiedCounter = append(modifiedCounter, nt.Uid)
		fmt.Printf("netTransaction '%s' original name '%s' modified name '%s'\n", nt.Uid, originalName, nt.Name)
	}

	fmt.Printf("modified network %d transactions %s\n", len(modifiedCounter), modifiedCounter)
	fmt.Println("script done")
}

type OutputNT struct {
	Input  models.NetworkTransaction `json:"input"`
	Output models.NetworkTransaction `json:"output"`
}

type OutputComplete struct {
	Modified    []map[string]OutputNT `json:"modified"`
	NotModified []string              `json:"notModified"`
}

/*
Script used to update network transactions that were created with wrong data,
caused by the manual payment. They we put in remittanceMga when should be
commissions. The wrong fields were: paymentType, accountType, amount and
amountNet.
*/
func UpdateManualPaymentNetworkTransactions(policyUids ...string) {
	output := OutputComplete{
		Modified:    make([]map[string]OutputNT, 0),
		NotModified: make([]string, 0),
	}
	log.AddPrefix("UpdateManualPaymentNetworkTransactions")
	defer log.PopPrefix()
	for _, policyUid := range policyUids {
		fmt.Printf("quering %s", policyUid)
		// get nettransaction by id
		query := fmt.Sprintf(
			"SELECT * FROM `%s.%s` WHERE policyUid = '%s' AND paymentType = '%s'",
			models.WoptaDataset,
			models.NetworkTransactionCollection,
			policyUid,
			models.PaymentTypeRemittanceMga,
		)
		netTransactions, err := lib.QueryRowsBigQuery[models.NetworkTransaction](query)
		if err != nil {
			fmt.Printf("error getting network transactions: %s", err.Error())
			return
		}
		if len(netTransactions) != 1 {
			fmt.Printf("expected 1 networkTransaction, got %d\n", len(netTransactions))
			output.NotModified = append(output.NotModified, policyUid)
			continue
		}
		fmt.Printf("found %d netTransactions\n", len(netTransactions))

		originalNetTransaction := netTransactions[0]
		modifiedNetTransaction := deepcopy.Copy(originalNetTransaction).(models.NetworkTransaction)

		policy, err := plc.GetPolicy(originalNetTransaction.PolicyUid)
		if err != nil {
			log.Error(err)
			return
		}
		networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode == nil {
			fmt.Println("error getting network node")
			return
		}
		warrant := networkNode.GetWarrant()
		if warrant == nil {
			fmt.Println("error getting warrant")
			return
		}
		prod := warrant.GetProduct(policy.Name)
		isActive := policy.ProducerUid == originalNetTransaction.NetworkNodeUid

		// update data
		commission := product.GetCommissionByProduct(&policy, prod, isActive)

		modifiedNetTransaction.PaymentType = models.PaymentTypeCommission
		modifiedNetTransaction.AccountType = models.AccountTypePassive
		modifiedNetTransaction.Amount = lib.RoundFloat(commission, 2)
		modifiedNetTransaction.AmountNet = lib.RoundFloat(commission, 2)

		output.Modified = append(output.Modified, map[string]OutputNT{
			originalNetTransaction.Uid: {
				Input:  originalNetTransaction,
				Output: modifiedNetTransaction,
			},
		})

		// save to bigquery
		err = saveBigQuery(modifiedNetTransaction)
		if err != nil {
			fmt.Printf("error saving to db: %s", err.Error())
			return
		}
		fmt.Println("NetworkTransaction saved!")
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		fmt.Printf("error marshaling output: %s", err.Error())
	}

	now := time.Now().UTC().Format(time.RFC3339)

	err = os.WriteFile(fmt.Sprintf("./%s_update_nt_manual_payment.json", now), outputJson, 0777)
	if err != nil {
		fmt.Printf("error writing output: %s", err.Error())
	}
}

func saveBigQuery(nt models.NetworkTransaction) error {
	updatedFields := make(map[string]interface{})

	updatedFields["paymentType"] = nt.PaymentType
	updatedFields["accountType"] = nt.AccountType
	updatedFields["amount"] = nt.Amount
	updatedFields["amountNet"] = nt.AmountNet

	return lib.UpdateRowBigQueryV2(
		models.WoptaDataset,
		models.NetworkTransactionCollection,
		updatedFields,
		fmt.Sprintf("WHERE uid = '%s'", nt.Uid),
	)
}
