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
		st.AddLocal("callbackInfo", &callbackInfo{RequestBody: []byte("ciao")})
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
		builder.AddHandler("baseCallback", funcTestWithInfo("baseCallback", log)),
		builder.AddHandler("errorCallbackConfig", func(sd bpnm.StorageData) error { return errors.New("callback client not set") }),
		builder.AddHandler("saveAudit", func(sd bpnm.StorageData) error {
			d, e := bpnm.GetData[*callbackInfo]("callbackInfo", sd)
			if e != nil {
				return e
			}
			if string(d.RequestBody) != "ciao" {
				return errors.New("no correct body request")
			}
			log.Println("saveAudit " + string(d.RequestBody))
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

func TestEmitForWinNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddLocal("policy", &policyEcommerce)
	store.AddLocal("node", &winNode)
	exps := []string{
		"winEmit",
		"saveAudit ciao",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowNode)
}

func TestEmitForBaseNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddLocal("policy", &policyEcommerce)
	store.AddLocal("node", &baseNode)
	exps := []string{
		"baseCallback",
		"saveAudit ciao",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowNode)
}

func TestEmitForBrokenNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddLocal("policy", &policyEcommerce)
	store.AddLocal("node", &brokenNode)

	log := mockLog{}
	build := getBuilderFlowNode(&log, store)
	flow, err := build.Build()
	if err != nil {
		t.Fatal(err)
	}
	err = flow.Run("emit")
	if err == nil {
		t.Fatal("Should have an error")
	}
}
