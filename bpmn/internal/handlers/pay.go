package handlers

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

func AddPayHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("updatePolicy", updatePolicy),
		builder.AddHandler("payTransaction", payTransaction),
	)
}

func updatePolicy(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var origin *flow.String
	var networkNode *flow.Network
	var flowName *flow.String
	var addresses *flow.Addresses
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
	)
	if err != nil {
		return err
	}

	if policy.IsPay || policy.Status != models.PolicyStatusToPay {
		log.ErrorF("policy already updated with isPay %t and status %s", policy.IsPay, policy.Status)
		return nil
	}

	// Add Policy contract
	err = plc.AddSignedDocumentsInPolicy(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("ERROR AddContract %s", err.Error())
		return err
	}

	// promote documents from temp bucket to user and connect it to policy
	err = plc.SetUserIntoPolicyContractor(policy.Policy, origin.String)
	if err != nil {
		log.ErrorF("ERROR SetUserIntoPolicyContractor %s", err.Error())
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

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	}

	// Send mail with the contract to the user
	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	return mail.SendMailContract(*policy.Policy, policy.Attachments, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)
}

func payTransaction(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var paymentInfo *flow.PaymentInfoBpmn
	var origin *flow.String
	var networkNode *flow.Network
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("paymentInfo", &paymentInfo, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
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
