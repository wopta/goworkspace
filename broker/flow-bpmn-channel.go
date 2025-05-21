package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/handlers/callbackFlow"
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
		builder.AddHandler("baseCallback", callbackFlow.BaseRequest),
		builder.AddHandler("winEmit", callbackFlow.CallBackEmit),
		builder.AddHandler("winSign", callbackFlow.CallBackSigned),
		builder.AddHandler("saveAudit", callbackFlow.SaveAudit),
		builder.AddHandler("winPay", callbackFlow.CallBackPaid),
		builder.AddHandler("winProposal", callbackFlow.CallBackProposal),
		builder.AddHandler("winRequestApproval", callbackFlow.CallBackRequestApproval),
		builder.AddHandler("winApproved", callbackFlow.CallBackApproved),
		builder.AddHandler("winRejected", callbackFlow.CallBackRejected),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
