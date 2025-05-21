package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/handlers/channelFlow"
)

func getNodeFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/node_flows.json")
	if e != nil {
		return nil, e
	}
	//hard coded, need to be on json
	callbackConf := flow.CallbackConfig{
		Proposal:        true,
		RequestApproval: true,
		Emit:            true,
		Pay:             true,
		Sign:            true,
		Approved:        true,
		Rejected:        true,
	}
	if e := store.AddLocal("config", &callbackConf); e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("baseCallback", channelFlow.BaseRequest),
		builder.AddHandler("winEmit", channelFlow.CallBackEmit),
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
