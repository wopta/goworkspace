package handlers

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/utility"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
)

func AddLeadHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setLeadData", setLeadBpmn),
		builder.AddHandler("sendLeadMail", sendLeadMail),
		builder.AddHandler("end_lead", savePolicy),
	)
}

func setLeadBpmn(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}
	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	utility.SetLeadData(policy.Policy, *mgaProduct, networkNode.NetworkNode)
	return nil
}

func sendLeadMail(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var addresses *flow.Addresses
	var flowName *flow.StringBpmn
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
