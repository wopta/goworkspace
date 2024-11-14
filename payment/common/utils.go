package common

import (
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func CheckPaymentModes(policy models.Policy) error {
	var allowedModes []string

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		allowedModes = models.GetAllowedMonthlyModes()
	case string(models.PaySplitYearly):
		allowedModes = models.GetAllowedYearlyModes()
	case string(models.PaySplitSingleInstallment):
		allowedModes = models.GetAllowedSingleInstallmentModes()
	}

	if !lib.SliceContains(allowedModes, policy.PaymentMode) {
		return fmt.Errorf("mode '%s' is incompatible with split '%s'", policy.PaymentMode, policy.PaymentSplit)
	}

	return nil
}

func SaveTransactionsToDB(transactions []models.Transaction, collection string) error {
	batch := make(map[string]map[string]models.Transaction)
	batch[collection] = make(map[string]models.Transaction)

	for idx := range transactions {
		transactions[idx].BigQueryParse()
		batch[collection][transactions[idx].Uid] = transactions[idx]
	}

	if err := lib.SetBatchFirestoreErr(batch); err != nil {
		log.Printf("error saving transactions to firestore: %s", err.Error())
		return err
	}

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, collection, transactions); err != nil {
		log.Printf("error saving transactions to bigquery: %s", err.Error())
		return err
	}

	return nil
}

func UpdatePaymentProvider(policy *models.Policy, provider string) error {
	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode == nil {
		return fmt.Errorf("can't find network node for policy %s ", policy.Uid)
	}
	warrant := networkNode.GetWarrant()
	if warrant == nil {
		return fmt.Errorf("can't find warrant for network node %s ", networkNode.Uid)
	}
	flow := warrant.GetFlowName(policy.Name)
	if flow == "" {
		return fmt.Errorf("can't find a flow for policy %s", policy.Uid)
	}

	switch provider {
	case models.ManualPaymentProvider:
		if flow == models.RemittanceMgaFlow {
			policy.Payment = provider
			return nil
		}
		return fmt.Errorf("can't update payment because flow %s doesn't support provider %s", flow, provider)
	case models.FabrickPaymentProvider:
		if flow != models.RemittanceMgaFlow {
			policy.Payment = provider
			return nil
		}
		return fmt.Errorf("can't update payment because flow %s doesn't support provider %s", flow, provider)
	default:
		return fmt.Errorf("can't update payment because provider %s is not supported", provider)
	}
}
