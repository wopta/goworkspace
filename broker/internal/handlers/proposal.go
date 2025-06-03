package handlers

import (
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddProposalHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("sendProposalMail", sendProposalMail),
		builder.AddHandler("setProposalData", setProposalData),
	)
}

func setProposalData(state bpmn.StorageData) error {
	var origin *flow.String
	var policy *flow.Policy
	var networkNode *flow.Network
	var mgaProduct *flow.Product
	var err = bpmn.IsError(
		bpmn.GetDataRef("origin", &origin, state),
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
		bpmn.GetDataRef("mgaProduct", &mgaProduct, state),
	)
	if err != nil {
		return err
	}
	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if policy.ProposalNumber != 0 {
		log.Printf("policy '%s' already has proposal with number '%d'", policy.Uid, policy.ProposalNumber)
		return nil
	}

	utility.SetProposalData(policy.Policy, origin.String, networkNode.NetworkNode, mgaProduct.Product)

	log.Printf("saving proposal n. %d to firestore...", policy.ProposalNumber)

	firePolicy := lib.GetDatasetByEnv(origin.String, lib.PolicyCollection)
	return lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
}

func sendProposalMail(state bpmn.StorageData) error {
	var policy *flow.Policy
	var addresses *flow.Addresses
	var sendEmail *flow.BoolBpmn
	var flowName *flow.String
	var networkNode *flow.Network
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, state),
		bpmn.GetDataRef("addresses", &addresses, state),
		bpmn.GetDataRef("sendEmail", &sendEmail, state),
		bpmn.GetDataRef("flowName", &flowName, state),
		bpmn.GetDataRef("networkNode", &networkNode, state),
	)
	if err != nil {
		return err
	}

	// TODO: remove when a proper solution to handle PMI is found
	if policy.Name == models.PmiProduct {
		return nil
	}

	if !sendEmail.Bool || policy.IsReserved {
		return nil
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

	mail.SendMailProposal(*policy.Policy, addresses.FromAddress, addresses.ToAddress, addresses.CcAddress, flowName.String, []string{models.ProposalAttachmentName})
	return nil
}
