package bpmn

import (
	"encoding/json"
	"fmt"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers/channelFlow"
	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func getNodeFlow(networkNode *models.NetworkNode) (*bpmnEngine.BpnmBuilder, error) {
	store := bpmnEngine.NewStorageBpnm()
	builder, err := bpmnEngine.NewBpnmBuilder("flows/draft/callback-flows.json")
	if err != nil {
		return nil, err
	}
	var callbackConfigFile []byte
	var client callback_out.CallbackClient
	switch networkNode.CallbackConfig.Name {
	case "winClient":
		callbackConfigFile, err = lib.GetFilesByEnvV2("flows/draft/callback/win.json")
		client = win.NewClient(networkNode.ExternalNetworkCode)
	case "facileBrokerClient":
		callbackConfigFile, err = lib.GetFilesByEnvV2("flows/draft/callback/base.json")
		client = base.NewClient(networkNode, "facile_broker")
	default:
		return nil, fmt.Errorf("CallbackCConfigName not valid '%v'", networkNode.CallbackConfig.Name)
	}

	var callbackConf flow.CallbackConfig
	err = json.Unmarshal(callbackConfigFile, &callbackConf)
	if err != nil {
		return nil, err
	}
	if err = store.AddLocal("config", &callbackConf); err != nil {
		return nil, err
	}
	if err = store.AddGlobal("clientCallback", &flow.ClientCallback{CallbackClient: client}); err != nil {
		return nil, err
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
		builder.AddHandler("saveAudit", channelFlow.CallBackRejected),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
