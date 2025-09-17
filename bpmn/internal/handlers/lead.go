package handlers

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func AddLeadHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("setLeadData", setLead),
		builder.AddHandler("sendLeadMail", sendLeadMail),
		builder.AddHandler("setProposalNumber", setProposalNumber),
		builder.AddHandler("end_lead", savePolicy),
	)
}

func setProposalNumber(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
	)
	if err != nil {
		return err
	}
	utility.SetProposalNumber(policy.Policy)
	return nil
}

func setLead(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	return utility.SetLeadData(policy.Policy, *mgaProduct, networkNode.NetworkNode)
}

func sendLeadMail(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var addresses *flow.Addresses
	var flowName *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
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
