package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers/channelFlow"
	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func injectCallbackFlow(networkNode *models.NetworkNode, mainBuilder *bpmnEngine.BpnmBuilder) error {
	if networkNode == nil || networkNode.CallbackConfig == nil {
		log.InfoF("no node or callback config available, no callback")
		return nil
	}
	store := bpmnEngine.NewStorageBpnm()
	builder, err := bpmnEngine.NewBpnmBuilder("flows/draft/callback-flows.json")
	if err != nil {
		return err
	}
	client, conf, err := callback_out.NewClient(networkNode)
	if err != nil {
		return err
	}
	if err = store.AddGlobal("config", &flow.CallbackConfigBpmn{CallbackConfig: conf}); err != nil {
		return err
	}
	if err = store.AddGlobal("clientCallback", &flow.ClientCallback{CallbackClient: client}); err != nil {
		return err
	}
	builder.SetStorage(store)
	err = bpmnEngine.IsError(
		builder.AddHandler("Emit", channelFlow.CallBackEmit),
		builder.AddHandler("Sign", channelFlow.CallBackSigned),
		builder.AddHandler("Pay", channelFlow.CallBackPaid),
		builder.AddHandler("Proposal", channelFlow.CallBackProposal),
		builder.AddHandler("RequestApproval", channelFlow.CallBackRequestApproval),
		builder.AddHandler("Approved", channelFlow.CallBackApproved),
		builder.AddHandler("Rejected", channelFlow.CallBackRejected),
		builder.AddHandler("saveAudit", channelFlow.SaveAudit),
	)
	if err != nil {
		return err
	}
	return mainBuilder.Inject(builder)
}
