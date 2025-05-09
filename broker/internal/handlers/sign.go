package handlers

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

func AddSignHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setSign", setSignDraft),
		builder.AddHandler("addContract", addContractDraft),
		builder.AddHandler("sendMailContract", sendMailContractDraft),
		builder.AddHandler("fillAttachments", fillAttachmentsDraft),
		builder.AddHandler("setToPay", setToPayDraft),
		builder.AddHandler("sendMailPay", sendMailPayDraft),
	)
}

func setSignDraft(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var origin *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}

	err = plc.Sign(policy.Policy, origin.String)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
func addContractDraft(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var origin *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}

	plc.AddContract(policy.Policy, origin.String)

	return nil
}
func sendMailContractDraft(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var networkNode *flow.NetworkDraft
	var flowName *flow.StringBpmn
	var origin *flow.StringBpmn
	var sendEmail *flow.BoolBpmn
	var addresses *flow.Addresses
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("addresses", &addresses, state),
	)
	if err != nil {
		return err
	}

	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail.Bool {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailContract(*policy.Policy, nil, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)

	return nil
}

func fillAttachmentsDraft(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var origin *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}
	err = plc.FillAttachments(policy.Policy, origin.String)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func setToPayDraft(state bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var origin *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}
	err = plc.SetToPay(policy.Policy, origin.String)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func sendMailPayDraft(state bpmn.StorageData) error {
	log.AddPrefix("sendMailPay")
	defer log.PopPrefix()

	var policy *flow.PolicyDraft
	var flowName *flow.StringBpmn
	var networkNode *flow.NetworkDraft
	var sendEmail *flow.BoolBpmn
	var addresses *flow.Addresses
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("sendEmail", &sendEmail, state),
	)
	if err != nil {
		return err
	}

	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		if sendEmail.Bool {
			addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
			addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		} else {
			addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
		}
	case models.MgaFlow, models.ECommerceFlow:
		addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	mail.SendMailPay(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)

	return nil
}
