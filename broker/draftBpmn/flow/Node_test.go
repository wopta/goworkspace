package flow

import (
	"errors"
	"os"
	"testing"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
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
		builder.AddHandler("Emit", funcTestWithInfo("Emit", log)),
		builder.AddHandler("EmitRemittance", funcTestWithInfo("EmitRemittance", log)),
		builder.AddHandler("Pay", funcTestWithInfo("Pay", log)),
		builder.AddHandler("Proposal", funcTestWithInfo("Proposal", log)),
		builder.AddHandler("RequestApproval", funcTestWithInfo("RequestApproval", log)),
		builder.AddHandler("Sign", funcTestWithInfo("Sign", log)),
		builder.AddHandler("Approved", funcTestWithInfo("Approved", log)),
		builder.AddHandler("Rejected", funcTestWithInfo("Approved", log)),
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
	winNode        = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "winClient"}}}
	baseNode       = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "facileBrokerClient"}}}
	brokenNode     = Network{&models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "booo"}}}
	callbackClient = ClientCallback{}
)

func TestEmitForWinNodeWithConfigTrue(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	store.AddLocal("config", &CallbackConfig{Emit: true})
	exps := []string{
		"Emit",
		"saveAudit prova request",
	}
	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinNodeWithConfigFalse(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	exps := []string{}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinWithProductFlowWinEmitRemittance(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("clientCallback", &callbackClient)

	exps := []string{
		"EmitRemittance",
	}
	store.AddLocal("config", &CallbackConfig{Emit: true})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}
func TestCallbackProposalWithIsReservedTrue(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	policyRes := Policy{&models.Policy{Channel: lib.ECommerceChannel, IsReserved: true}}
	store.AddGlobal("policy", &policyRes)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	store.AddLocal("config", &CallbackConfig{Proposal: true})
	exps := []string{}
	testFlow(t, "proposalCallback", exps, store, getBuilderFlowNode)
}

func TestCallbackProposalWithIsReservedFalse(t *testing.T) {
	store := bpmn.NewStorageBpnm()
	policyRes := Policy{&models.Policy{Channel: lib.ECommerceChannel, IsReserved: false}}
	store.AddGlobal("policy", &policyRes)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	store.AddLocal("config", &CallbackConfig{Proposal: true})
	exps := []string{
		"Proposal",
		"saveAudit prova request",
	}
	testFlow(t, "proposalCallback", exps, store, getBuilderFlowNode)
}
