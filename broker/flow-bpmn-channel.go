package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/handlers/callback"
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
		builder.AddHandler("baseCallback", callback.BaseRequest),
		builder.AddHandler("winEmit", callback.CallBackEmit),
		builder.AddHandler("winSign", callback.CallBackSigned),
		builder.AddHandler("saveAudit", callback.SaveAudit),
		builder.AddHandler("winPay", callback.CallBackPaid),
		builder.AddHandler("winProposal", callback.CallBackProposal),
		builder.AddHandler("winRequestApproval", callback.CallBackRequestApproval),
		builder.AddHandler("winApproved", callback.CallBackApproved),
		builder.AddHandler("winRejected", callback.CallBackRejected),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
