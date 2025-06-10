package flow

import (
	"errors"
	"os"
	"testing"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func funcTestWithInfo(message string, log *mockLog) func(bpmn.StorageData) error {
	return func(st bpmn.StorageData) error {
		log.println(message)
		st.AddLocal("callbackInfo", &CallbackInfo{base.CallbackInfo{ReqBody: []byte("prova request")}})
		return nil
	}
}

func getBuilderFlowNode(log *mockLog, store bpmn.StorageData) (*bpmn.BpnmBuilder, error) {
	os.Setenv("env", env.LocalTest)
	builder, e := bpmn.NewBpnmBuilderRawPath("../../../../function-data/dev//flows/draft/node_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpmn.IsError(
		builder.AddHandler("winEmit", funcTestWithInfo("winEmit", log)),
		builder.AddHandler("winEmitRemittance", funcTestWithInfo("winEmitRemittance", log)),
		builder.AddHandler("winPay", funcTestWithInfo("winPay", log)),
		builder.AddHandler("winProposal", funcTestWithInfo("winProposal", log)),
		builder.AddHandler("winRequestApproval", funcTestWithInfo("winRequestApproval", log)),
		builder.AddHandler("winSign", funcTestWithInfo("winSign", log)),
		builder.AddHandler("winApproved", funcTestWithInfo("winApproved", log)),
		builder.AddHandler("winRejected", funcTestWithInfo("winApproved", log)),
		builder.AddHandler("baseCallback", funcTestWithInfo("baseCallback", log)),
		builder.AddHandler("saveAudit", func(sd bpmn.StorageData) error {
			d, e := bpmn.GetData[*CallbackInfo]("callbackInfo", sd)
			if e != nil {
				return e
			}
			if string(d.ReqBody) != "prova request" {
				return errors.New("no correct body request")
			}
			log.println("saveAudit " + string(d.ReqBody))
			return nil
		}),
	)
	if e != nil {
		return nil, e
	}
	return builder, nil
}

var (
	winNode    = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "winClient"}}}
	baseNode   = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "facileBrokerClient"}}}
	brokenNode = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "booo"}}}
)

func TestEmitForWinNodeWithConfigTrue(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)

	store.AddLocal("config", &CallbackConfig{Emit: true})
	exps := []string{
		"winEmit",
		"saveAudit prova request",
	}
	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinNodeWithConfigFalse(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)

	exps := []string{}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinWithProductFlowWinEmitRemittance(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{
		"winEmitRemittance",
	}
	store.AddLocal("config", &CallbackConfig{Emit: true})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}
