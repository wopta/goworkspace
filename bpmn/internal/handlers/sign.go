package handlers

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func AddSignHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("setSign", setSign),
		builder.AddHandler("addContract", addContract),
		builder.AddHandler("sendMailContract", sendMailContract),
		builder.AddHandler("fillAttachments", fillAttachments),
		builder.AddHandler("setToPay", setToPay),
		builder.AddHandler("sendMailPay", sendMailPay),
	)
}

func setSign(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
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
func addContract(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}

	plc.AddSignedDocumentsInPolicy(policy.Policy, origin.String)

	return nil
}
func sendMailContract(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var flowName *flow.String
	var origin *flow.String
	var sendEmail *flow.BoolBpmn
	var addresses *flow.Addresses
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("sendEmail", &sendEmail, state),
	)
	if err != nil {
		return err
	}

	addresses.ToAddress = mail.GetContractorEmail(policy.Policy)
	switch flowName.String {
	case models.ProviderMgaFlow, models.RemittanceMgaFlow:
		addresses.CcAddress = mail.GetNetworkNodeEmail(networkNode.NetworkNode)
	}

	log.Printf(
		"from '%s', to '%s', cc '%s'",
		addresses.FromAddress.String(),
		addresses.ToAddress.String(),
		addresses.CcAddress.String(),
	)
	return mail.SendMailContract(*policy.Policy, policy.Attachments, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)
}

func fillAttachments(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
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

func setToPay(state *bpmnEngine.StorageBpnm) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
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

func sendMailPay(state *bpmnEngine.StorageBpnm) error {
	log.AddPrefix("sendMailPay")
	defer log.PopPrefix()

	var policy *flow.Policy
	var flowName *flow.String
	var networkNode *flow.Network
	var sendEmail *flow.BoolBpmn
	var addresses *flow.Addresses
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("sendEmail", &sendEmail, state),
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
