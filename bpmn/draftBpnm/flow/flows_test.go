package flow

import (
	"testing"

	bpnm "github.com/wopta/goworkspace/bpmn/draftBpnm"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type mockLog struct {
	log []string
}

func (m *mockLog) println(mes string) {
	m.log = append(m.log, mes)
}

func (m *mockLog) printlnToTesting(t *testing.T) {
	t.Log("Actual log: ")
	for _, mes := range m.log {
		t.Log(" ", mes)
	}
}
func funcTest(message string, log *mockLog) func(bpnm.StorageData) error {
	return func(sd bpnm.StorageData) error {
		log.println(message)
		return nil
	}
}

func getBuilderFlowChannel(log *mockLog, store bpnm.StorageData) *bpnm.BpnmBuilder {

	builder, e := bpnm.NewBpnmBuilder("channel_flows.json")
	if e != nil {
		return nil
	}
	builder.SetStorage(store)
	e = bpnm.IsError(
		builder.AddHandler("setProposalData", funcTest("setProposalData", log)),
		builder.AddHandler("emitData", funcTest("emitData", log)),
		builder.AddHandler("sendMailSign", funcTest("sendMailSign", log)),
		builder.AddHandler("pay", funcTest("pay", log)),
		builder.AddHandler("setAdvice", funcTest("setAdvice", log)),
		builder.AddHandler("putUser", funcTest("putUser", log)),
		builder.AddHandler("sendEmitProposalMail", funcTest("sendEmitProposalMail", log)),
		builder.AddHandler("setLeadData", funcTest("setLeadData", log)),
		builder.AddHandler("sendLeadMail", funcTest("sendLeadMail", log)),
		builder.AddHandler("updatePolicy", funcTest("updatePolicy", log)),
		builder.AddHandler("sign", funcTest("sign", log)),
		builder.AddHandler("payTransaction", funcTest("payTransaction", log)),
		builder.AddHandler("sendProposalMail", funcTest("sendProposalMail", log)),
		builder.AddHandler("fillAttachments", funcTest("fillAttachments", log)),
		builder.AddHandler("setToPay", funcTest("setToPay", log)),
		builder.AddHandler("setSign", funcTest("setSign", log)),
		builder.AddHandler("sendMailContract", funcTest("sendMailContract", log)),
		builder.AddHandler("sendMailPay", funcTest("sendMailPay", log)),
		builder.AddHandler("setRequestApprovalData", funcTest("setRequestApprovalData", log)),
		builder.AddHandler("sendRequestApprovalMail", funcTest("sendRequestApprovalMail", log)),
		builder.AddHandler("addContract", funcTest("addContract", log)),
	)
	if e != nil {
		return nil
	}

	return builder
}

type builderFlow func(*mockLog, bpnm.StorageData) *bpnm.BpnmBuilder

func testFlow(t *testing.T, process string, expectedACtivities []string, store bpnm.StorageData, builbuilderFlow builderFlow) {
	log := mockLog{}
	build := builbuilderFlow(&log, store)
	flow, e := build.Build()
	if e != nil {
		t.Fatal(e)
	}
	if err := flow.Run(process); err != nil {
		t.Fatal(err)
	}
	if len(expectedACtivities) != len(log.log) {
		log.printlnToTesting(t)
		for _, mes := range log.log {
			t.Log(mes)
		}
		t.Fatalf("exp n message: %v,got: %v", len(expectedACtivities), len(log.log))
	}
	for i, exp := range expectedACtivities {
		log.printlnToTesting(t)
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

var (
	policyEcommerce = models.Policy{Channel: lib.ECommerceChannel}
	policyMga       = models.Policy{Channel: lib.MgaChannel}
	policyNetwork   = models.Policy{Channel: lib.NetworkChannel}
)

// product
var (
	productEcommerce     = models.Product{Flow: models.ECommerceFlow}
	productMga           = models.Product{Flow: models.MgaFlow}
	productProviderMga   = models.Product{Flow: models.ProviderMgaFlow}
	productRemittanceMga = models.Product{Flow: models.RemittanceMgaFlow}
)

func TestEmitForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("node", &winNode)
	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("node", &winNode)
	exps := []string{
		"setLeadData",
		"sendLeadMail",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("node", &winNode)

	exps := []string{}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestPayForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestSignForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setProposalData",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestPayForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestSignForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setProposalData",
		"sendProposalMail",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestPayForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestSignForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"sendMailSign",
		"setAdvice",
		"putUser",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setProposalData",
		"sendProposalMail",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestPayForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestSignForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("node", &winNode)

	exps := []string{
		"setSign",
		"addContract",
		"sendMailContract",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

// With node flow
func TestEmitForEcommerceWithNodeFlow(t *testing.T) {
	storeFlowChannel := bpnm.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyEcommerce)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("node", &winNode)

	storeNode := bpnm.NewStorageBpnm()
	storeNode.AddLocal("config", &callbackConfig{Events: map[string]bool{"emit": true}})

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
		"winEmit",
		"saveAudit ciao",
	}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd bpnm.StorageData) *bpnm.BpnmBuilder {
		build := getBuilderFlowChannel(log, storeFlowChannel)
		nodeBuild := getBuilderFlowNode(log, storeNode)
		if e := build.Inject(nodeBuild); e != nil {
			t.Fatal(e)
		}

		return build
	})
}
func TestEmitForWgaWithNodeFlow(t *testing.T) {
	storeFlowChannel := bpnm.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyMga)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("node", &winNode)

	storeNode := bpnm.NewStorageBpnm()
	storeNode.AddLocal("config", &callbackConfig{Events: map[string]bool{"emit": true}})

	exps := []string{}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd bpnm.StorageData) *bpnm.BpnmBuilder {
		build := getBuilderFlowChannel(log, storeFlowChannel)
		nodeBuild := getBuilderFlowNode(log, storeNode)
		if e := build.Inject(nodeBuild); e != nil {
			t.Fatal(e)
		}

		return build
	})
}
func TestEmitForEcommerceWithNodeFlowConfFalse(t *testing.T) {
	storeFlowChannel := bpnm.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyEcommerce)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("node", &winNode)

	storeNode := bpnm.NewStorageBpnm()
	storeNode.AddLocal("config", &callbackConfig{Events: map[string]bool{"emit": false}})

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd bpnm.StorageData) *bpnm.BpnmBuilder {
		build := getBuilderFlowChannel(log, storeFlowChannel)
		nodeBuild := getBuilderFlowNode(log, storeNode)
		if e := build.Inject(nodeBuild); e != nil {
			t.Fatal(e)
		}

		return build
	})
}
