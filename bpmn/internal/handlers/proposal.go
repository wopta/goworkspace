package handlers

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddProposalHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("sendProposalMail", sendProposalMail),
		builder.AddHandler("setProposalData", setProposalData),
		builder.AddHandler("end_proposal", endProposal),
	)
}

func endProposal(state bpmnEngine.StorageData) error {
	var policy *flow.Policy
	var isProposal *flow.BoolBpmn
	var origin *flow.String
	var networkNode *flow.Network
	var mgaProduct *flow.Product
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("is_PROPOSAL_V2", &isProposal, state),
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
	)
	if err != nil {
		return err
	}
	if !isProposal.Bool {
		utility.SetProposalNumber(policy.Policy)
		policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	}
	return nil
}

func setProposalData(state bpmnEngine.StorageData) error {
	var origin *flow.String
	var policy *flow.Policy
	var networkNode *flow.Network
	var mgaProduct *flow.Product
	var err = bpmnEngine.IsError(
		bpmnEngine.GetDataRef("origin", &origin, state),
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
		bpmnEngine.GetDataRef("mgaProduct", &mgaProduct, state),
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

func sendProposalMail(state bpmnEngine.StorageData) error {
	var policy *flow.Policy
	var addresses *flow.Addresses
	var sendEmail *flow.BoolBpmn
	var flowName *flow.String
	var networkNode *flow.Network
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, state),
		bpmnEngine.GetDataRef("addresses", &addresses, state),
		bpmnEngine.GetDataRef("sendEmail", &sendEmail, state),
		bpmnEngine.GetDataRef("flowName", &flowName, state),
		bpmnEngine.GetDataRef("networkNode", &networkNode, state),
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
