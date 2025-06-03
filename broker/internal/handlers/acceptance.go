package handlers

import (
	"time"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	draftbpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func AddAcceptanceHandlers(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		builder.AddHandler("rejected", draftRejectPolicy),
		builder.AddHandler("approved", draftApprovePolicy),
		builder.AddHandler("sendAcceptanceMail", sendAcceptanceMail),
		builder.AddHandler("end_accepance", savePolicy),
	)
}

func draftRejectPolicy(storage draftbpmn.StorageData) error {
	var policy *flow.Policy
	var action *flow.String
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, storage),
		bpmn.GetDataRef("action", &action, storage),
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

func draftApprovePolicy(storage draftbpmn.StorageData) error {
	var policy *flow.Policy
	var action *flow.String
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, storage),
		bpmn.GetDataRef("action", &action, storage),
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

func sendAcceptanceMail(state bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.Policy]("policy", state)
	if err != nil {
		return err
	}
	addresses, err := bpmn.GetData[*flow.Addresses]("addresses", state)
	if err != nil {
		return err
	}

	node, err := bpmn.GetData[*flow.Network]("networkNode", state)
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
