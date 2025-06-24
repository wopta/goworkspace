package handlers

import (
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func AddLeadHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setLeadData", setLead),
		builder.AddHandler("sendLeadMail", sendLeadMail),
		builder.AddHandler("end_lead", savePolicy),
	)
}

func setLead(state bpmn.StorageData) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	return utility.SetLeadData(policy.Policy, *mgaProduct, networkNode.NetworkNode)
}

func sendLeadMail(state bpmn.StorageData) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var addresses *flow.Addresses
	var flowName *flow.String
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("flowName", &flowName, state),
	)
	if err != nil {
		return err
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	addresses.CcAddress = mail.Address{}
	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	}

	log.Printf(
		"[sendLeadMail] from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailLead(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String, []string{})
	return nil
}
