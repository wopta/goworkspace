package flow

import (
	"os"
	"testing"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func funcTestWithInfo(message string, log *mockLog) func(*bpmnEngine.StorageBpnm) error {
	return func(st *bpmnEngine.StorageBpnm) error {
		log.println(message)
		return st.AddLocal("callbackInfo", &CallbackInfo{base.CallbackInfo{ReqBody: []byte("prova request " + message)}})
	}
}

func getBuilderFlowNode(log *mockLog, store *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
	os.Setenv("env", env.LocalTest)
	builder, e := bpmnEngine.NewBpnmBuilderRawPath("../../../../function-data/dev//flows/draft/callback-flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpmnEngine.IsError(
		builder.AddHandler("Emit", funcTestWithInfo("Emit", log)),
		builder.AddHandler("Pay", funcTestWithInfo("Pay", log)),
		builder.AddHandler("Proposal", funcTestWithInfo("Proposal", log)),
		builder.AddHandler("RequestApproval", funcTestWithInfo("RequestApproval", log)),
		builder.AddHandler("Sign", funcTestWithInfo("Sign", log)),
		builder.AddHandler("Approved", funcTestWithInfo("Approved", log)),
		builder.AddHandler("Rejected", funcTestWithInfo("Approved", log)),
		builder.AddHandler("saveAudit", func(sd *bpmnEngine.StorageBpnm) error {
			d, e := bpmnEngine.GetData[*CallbackInfo]("callbackInfo", sd)
			if e != nil {
				return e
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
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	store.AddLocal("config", &CallbackConfig{Emit: true})
	exps := []string{
		"Emit",
		"saveAudit prova request Emit",
	}
	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinNodeWithConfigFalse(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	exps := []string{}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}

func TestEmitForWinWithProductFlowWinEmitRemittance(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("clientCallback", &callbackClient)

	exps := []string{
		"Emit",
		"saveAudit prova request Emit",
		"Pay",
		"saveAudit prova request Pay",
	}
	store.AddLocal("config", &CallbackConfig{Emit: true})

	testFlow(t, "emitCallBack", exps, store, getBuilderFlowNode)
}
func TestCallbackProposalWithIsReservedTrue(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
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
	store := bpmnEngine.NewStorageBpnm()
	policyRes := Policy{&models.Policy{Channel: lib.ECommerceChannel, IsReserved: false}}
	store.AddGlobal("policy", &policyRes)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("clientCallback", &callbackClient)

	store.AddLocal("config", &CallbackConfig{Proposal: true})
	exps := []string{
		"Proposal",
		"saveAudit prova request Proposal",
	}
	testFlow(t, "proposalCallback", exps, store, getBuilderFlowNode)
}
