package handlers

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddAcceptanceHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		builder.AddHandler("rejected", draftRejectPolicy),
		builder.AddHandler("approved", draftApprovePolicy),
		builder.AddHandler("sendAcceptanceMail", sendAcceptanceMail),
		builder.AddHandler("end_accepance", savePolicy),
	)
}

func draftRejectPolicy(storage bpmnEngine.StorageData) error {
	var policy *flow.Policy
	var action *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, storage),
		bpmnEngine.GetDataRef("action", &action, storage),
	)
	if err != nil {
		return err
	}
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = action.String
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s REJECTED", policy.Uid)
	policy.Updated = time.Now().UTC()
	return nil
}

func draftApprovePolicy(storage bpmnEngine.StorageData) error {
	var policy *flow.Policy
	var action *flow.String
	err := bpmnEngine.IsError(
		bpmnEngine.GetDataRef("policy", &policy, storage),
		bpmnEngine.GetDataRef("action", &action, storage),
	)
	if err != nil {
		return err
	}
	policy.Status = models.PolicyStatusApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = action.String
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s APPROVED", policy.Uid)
	policy.Updated = time.Now().UTC()
	return nil
}

func sendAcceptanceMail(state bpmnEngine.StorageData) error {
	policy, err := bpmnEngine.GetData[*flow.Policy]("policy", state)
	if err != nil {
		return err
	}
	addresses, err := bpmnEngine.GetData[*flow.Addresses]("addresses", state)
	if err != nil {
		return err
	}

	node, err := bpmnEngine.GetData[*flow.Network]("networkNode", state)
	if err != nil {
		return err
	}
	log.Printf("toAddress '%s'", addresses.ToAddress.String())
	var warrant *models.Warrant
	if node.NetworkNode != nil {
		warrant = node.GetWarrant()
	}
	flowName, _ := policy.GetFlow(node.NetworkNode, warrant)
	mail.SendMailReservedResult(
		*policy.Policy,
		mail.AddressAssunzione,
		addresses.ToAddress,
		mail.Address{},
		flowName,
	)
	return nil
}
