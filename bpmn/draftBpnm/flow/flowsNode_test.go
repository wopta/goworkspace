package flow

import (
	"errors"
	"net/http"
	"testing"

	bpnm "github.com/wopta/goworkspace/bpmn/draftBpnm"
	"github.com/wopta/goworkspace/models"
)

type callbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
}

func (c *callbackInfo) GetType() string {
	return "callbackInfo"
}

func funcTestWithInfo(message string, log *mockLog) func(bpnm.StorageData) error {
	return func(st bpnm.StorageData) error {
		log.Println(message)
		st.AddLocal("callbackInfo", &callbackInfo{})
		return nil
	}
}

func builderFlowNode(log *mockLog, store bpnm.StorageData) (*bpnm.FlowBpnm, error) {

	builder, e := bpnm.NewBpnmBuilder("node_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpnm.IsError(
		builder.AddHandler("WinEmit", funcTestWithInfo("WinEmit", log)),
		builder.AddHandler("BaseCallback", funcTestWithInfo("BaseCallback", log)),
		builder.AddHandler("ErrorCallbackConfig", func(sd bpnm.StorageData) error { return errors.New("callback client not set") }),
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

	log := mockLog{}
	flow, err := builderFlowNode(&log, store)
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err == nil {
		t.Fatal("Should have an error")
	}
}
