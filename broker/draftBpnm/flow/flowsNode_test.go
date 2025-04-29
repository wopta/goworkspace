package flow

import (
	"errors"
	"testing"

	bpnm "github.com/wopta/goworkspace/broker/draftBpnm"
	"github.com/wopta/goworkspace/models"
)

func funcTestWithInfo(message string, log *mockLog) func(bpnm.StorageData) error {
	return func(st bpnm.StorageData) error {
		log.println(message)
		st.AddLocal("callbackInfo", &CallbackInfo{RequestBody: []byte("prova request")})
		return nil
	}
}

func getBuilderFlowNode(log *mockLog, store bpnm.StorageData) *bpnm.BpnmBuilder {
	builder, e := bpnm.NewBpnmBuilder("node_flows.json")
	if e != nil {
		return nil
	}
	builder.SetStorage(store)
	e = bpnm.IsError(
		builder.AddHandler("winEmit", funcTestWithInfo("winEmit", log)),
		builder.AddHandler("winLead", funcTestWithInfo("winLead", log)),
		builder.AddHandler("winPay", funcTestWithInfo("winPay", log)),
		builder.AddHandler("winProposal", funcTestWithInfo("winProposal", log)),
		builder.AddHandler("winRequestApproval", funcTestWithInfo("winRequestApproval", log)),
		builder.AddHandler("winSign", funcTestWithInfo("winSign", log)),
		builder.AddHandler("baseCallback", funcTestWithInfo("baseCallback", log)),
		builder.AddHandler("saveAudit", func(sd bpnm.StorageData) error {
			d, e := bpnm.GetData[*CallbackInfo]("callbackInfo", sd)
			if e != nil {
				return e
			}
			if string(d.RequestBody) != "prova request" {
				return errors.New("no correct body request")
			}
			log.println("saveAudit " + string(d.RequestBody))
			return nil
		}),
	)
	if e != nil {
		return nil
	}
	return builder
}

var (
	winNode    = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "winClient"}}
	baseNode   = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "facileBrokerClient"}}
	brokenNode = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "booo"}}
)

func TestEmitForWinNodeWithConfigTrue(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &winNode)

	store.AddLocal("config", &CallbackConfig{Emit: true})
	exps := []string{
		"winEmit",
		"saveAudit prova request",
	}
	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinNodeWithConfigFalse(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &winNode)

	exps := []string{}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}
