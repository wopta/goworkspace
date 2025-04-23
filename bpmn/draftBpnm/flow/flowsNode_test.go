package flow

import (
	"testing"

	bpnm "github.com/wopta/goworkspace/bpmn/draftBpnm"
	"github.com/wopta/goworkspace/models"
)

func builderFlowNode(log *mockLog, store bpnm.StorageData) (*bpnm.FlowBpnm, error) {

	builder, e := bpnm.NewBpnmBuilder("node_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpnm.IsError(
		builder.AddHandler("WinEmit", funcTest("WinEmit", log)),
		builder.AddHandler("BaseCallback", funcTest("BaseCallback", log)),
		builder.AddHandler("ErrorCallbackConfig", funcTest("ErrorCallbackConfig", log)),
		builder.AddHandler("SaveAudit", funcTest("SaveAudit", log)),
	)
	if e != nil {
		return nil, e
	}
	return builder.Build()
}

var (
	winNode    = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "winClient"}}
	baseNode   = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "facileBrokerClient"}}
	brokenNode = models.NetworkNode{CallbackConfig: &models.CallbackConfig{Name: "booo"}}
)

func TestEmitForWinNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &winNode)
	exps := []string{
		"WinEmit",
		"SaveAudit",
	}
	testFlow(t, "emit", exps, store, builderFlowNode)
}

func TestEmitForBaseNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &baseNode)
	exps := []string{
		"BaseCallback",
		"SaveAudit",
	}
	testFlow(t, "emit", exps, store, builderFlowNode)
}

func TestEmitForBrokenNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &brokenNode)
	exps := []string{
		"ErrorCallbackConfig",
	}
	testFlow(t, "emit", exps, store, builderFlowNode)
}
