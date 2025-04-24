package flow

import (
	"errors"
	"net/http"
	"strings"
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
		builder.AddHandler("winLead", funcTestWithInfo("winLead", log)),
		builder.AddHandler("winPay", funcTestWithInfo("winPay", log)),
		builder.AddHandler("winProposal", funcTestWithInfo("winProposal", log)),
		builder.AddHandler("winRequestApproval", funcTestWithInfo("winRequestApproval", log)),
		builder.AddHandler("winSign", funcTestWithInfo("winSign", log)),
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

var processesToTest = [...]string{"emit", "lead", "pay", "proposal", "requestApproval", "sign"}

func TestEmitForWinNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &winNode)
	for _, process := range processesToTest {
		namelog := strings.Replace(process, string(process[0]), string(process[0]-32), 1) //upper case first letter
		exps := []string{
			"win" + namelog,
			"saveAudit ciao",
		}
		testFlow(t, process, exps, store, getBuilderFlowNode)
	}
}

func TestBaseNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &baseNode)
	exps := []string{
		"baseCallback",
		"saveAudit ciao",
	}
	for _, process := range processesToTest {
		testFlow(t, process, exps, store, getBuilderFlowNode)
	}
}

func TestBrokenNode(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("node", &brokenNode)

	log := mockLog{}
	build := getBuilderFlowNode(&log, store)
	flow, err := build.Build()
	if err != nil {
		t.Fatal(err)
	}
	for _, process := range processesToTest {
		err = flow.Run(process)
		if err == nil {
			t.Fatal("Should have an error")
		}
	}
}
