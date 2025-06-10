package broker

import (
	"encoding/json"
	"fmt"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/handlers/channelFlow"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

func getNodeFlow(callbackConfigName string) (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, err := bpmn.NewBpnmBuilder("flows/draft/node_flows.json")
	if err != nil {
		return nil, err
	}
	var callbackConfigFile []byte
	switch callbackConfigName {
	case "winClient":
		callbackConfigFile, err = lib.GetFilesByEnvV2("flows/draft/callback/win.json")
	case "facileBrokerClient":
		callbackConfigFile, err = lib.GetFilesByEnvV2("flows/draft/callback/base.json")
	default:
		return nil, fmt.Errorf("CallbackCConfigName not valid '%v'", callbackConfigName)
	}
	var callbackConf flow.CallbackConfig
	err = json.Unmarshal(callbackConfigFile, &callbackConf)
	if err != nil {
		return nil, err
	}
	if err = store.AddLocal("config", &callbackConf); err != nil {
		return nil, err
	}
	builder.SetStorage(store)
	err = bpmn.IsError(
		builder.AddHandler("baseCallback", channelFlow.BaseRequest),
		builder.AddHandler("winEmit", channelFlow.CallBackEmit),
		builder.AddHandler("winEmitRemittance", channelFlow.CallBackEmitRemittance),
		builder.AddHandler("winSign", channelFlow.CallBackSigned),
		builder.AddHandler("saveAudit", channelFlow.SaveAudit),
		builder.AddHandler("winPay", channelFlow.CallBackPaid),
		builder.AddHandler("winProposal", channelFlow.CallBackProposal),
		builder.AddHandler("winRequestApproval", channelFlow.CallBackRequestApproval),
		builder.AddHandler("winApproved", channelFlow.CallBackApproved),
		builder.AddHandler("winRejected", channelFlow.CallBackRejected),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
