package handlers

import (
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func AddSignHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("setSign", setSign),
		builder.AddHandler("addContract", addContract),
		builder.AddHandler("sendMailContract", sendMailContract),
		builder.AddHandler("fillAttachments", fillAttachments),
		builder.AddHandler("setToPay", setToPay),
		builder.AddHandler("sendMailPay", sendMailPay),
	)
}

func setSign(state bpmn.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
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
func addContract(state bpmn.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
	)
	if err != nil {
		return err
	}

	plc.AddSignedDocumentsInPolicy(policy.Policy, origin.String)

	return nil
}
func sendMailContract(state bpmn.StorageData) error {
	var policy *flow.Policy
	var networkNode *flow.Network
	var flowName *flow.String
	var origin *flow.String
	var sendEmail *flow.BoolBpmn
	var addresses *flow.Addresses
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("origin", &origin, state),
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
	return mail.SendMailContract(*policy.Policy, policy.Attachments, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String)
}

func fillAttachments(state bpmn.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
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

func setToPay(state bpmn.StorageData) error {
	var policy *flow.Policy
	var origin *flow.String
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

func sendMailPay(state bpmn.StorageData) error {
	log.AddPrefix("sendMailPay")
	defer log.PopPrefix()

	var policy *flow.Policy
	var flowName *flow.String
	var networkNode *flow.Network
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
