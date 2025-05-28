package common

import (
	"fmt"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
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
	case string(models.PaySplitSemestral):
		//TODO: what is allowed?
	}
	if !lib.SliceContains(allowedModes, policy.PaymentMode) {
		return fmt.Errorf("mode '%s' is incompatible with split '%s'", policy.PaymentMode, policy.PaymentSplit)
	}

	return nil
}

func SaveTransactionsToDB(transactions []models.Transaction, collection string) error {
	return tr.SaveTransactionsToDB(transactions, collection)
}

func checkProviderCompatibility(provider, flow string) error {
	if provider == models.ManualPaymentProvider {
		if flow != models.RemittanceMgaFlow {
			return fmt.Errorf("can't update payment because flow %s doesn't support provider %s", flow, provider)
		}
	} else if provider == models.FabrickPaymentProvider {
		if flow == models.RemittanceMgaFlow {
			return fmt.Errorf("can't update payment because flow %s doesn't support provider %s", flow, provider)
		}
	} else {
		return fmt.Errorf("can't update payment because provider %s is not supported", provider)
	}
	return nil
}

func UpdatePaymentProvider(policy *models.Policy, provider string) error {
	if policy.Channel == models.NetworkChannel {
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
		err := checkProviderCompatibility(provider, flow)
		if err != nil {
			return err
		}
	}
	policy.Payment = provider
	return nil
}
