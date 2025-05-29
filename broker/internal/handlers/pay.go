package handlers

import (
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func AddPayHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("updatePolicy", updatePolicy),
		builder.AddHandler("payTransaction", payTransaction),
	)
}

func updatePolicy(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var origin *flow.StringBpmn
	var networkNode *flow.NetworkDraft
	var flowName *flow.StringBpmn
	var addresses *flow.Addresses
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("addresses", &addresses, state),
	)
	if err != nil {
		return err
	}

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.ErrorF("policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("ERROR SetUserIntoPolicyContractor %s", err.Error())
		return err
	}

	// Add Policy contract
	err = plc.AddSignedDocumentsInPolicy(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("ERROR AddContract %s", err.Error())
		return err
	}

	// Update Policy as paid
	err = plc.Pay(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("ERROR Policy Pay %s", err.Error())
		return err
	}

	err = network.UpdateNetworkNodePortfolio(origin.String, policy.Policy, networkNode.NetworkNode)
	if err != nil {
		log.ErrorF("error updating %s portfolio %s", networkNode.Type, err.Error())
		return err
	}

	policy.BigquerySave(origin.String)

	switch flowName.String {
	case models.ProviderMgaFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	case models.RemittanceMgaFlow:
		addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	// Send mail with the contract to the user
	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)

	return nil
}

func payTransaction(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var paymentInfo *flow.PaymentInfoBpmn
	var origin *flow.StringBpmn
	var networkNode *flow.NetworkDraft
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("paymentInfo", &paymentInfo, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}

	providerId := paymentInfo.FabrickCallback.PaymentID
	transaction, _ := tr.GetTransactionToBePaid(policy.Uid, *providerId, paymentInfo.Schedule, lib.TransactionsCollection)
	err = tr.Pay(&transaction, origin.String, paymentInfo.PaymentMethod)
	if err != nil {
		log.Error(err)
		return err
	}

	transaction.BigQuerySave(origin.String)

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	return tr.CreateNetworkTransactions(policy.Policy, &transaction, networkNode.NetworkNode, mgaProduct)
}
