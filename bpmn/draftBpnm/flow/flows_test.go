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

func (m *mockLog) Println(mes string) {
	m.log = append(m.log, mes)
}
func funcTest(message string, log *mockLog) func(bpnm.StorageData) error {
	return func(sd bpnm.StorageData) error {
		log.Println(message)
		return nil
	}
}

func getBuilder(log *mockLog, store bpnm.StorageData) (*bpnm.FlowBpnm, error) {

	builder, e := bpnm.NewBpnmBuilder("channel_flows.json")
	if e != nil {
		return nil, e
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
		return nil, e
	}

	return builder.Build()
}

func testFlow(t *testing.T, process string, expectedACtivities []string, store bpnm.StorageData) {
	log := mockLog{}
	flow, e := getBuilder(&log, store)
	if e != nil {
		t.Fatal(e)
	}
	if err := flow.Run(process); err != nil {
		t.Fatal(err)
	}
	if len(expectedACtivities) != len(log.log) {
		for _, mes := range log.log {
			t.Log(mes)
		}
		t.Fatalf("exp n message: %v,got: %v", len(expectedACtivities), len(log.log))
	}
	for i, exp := range expectedACtivities {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

var policyEcommerce = models.Policy{Channel: lib.ECommerceChannel}
var policyMga = models.Policy{Channel: lib.MgaChannel}
var policyNetwork = models.Policy{Channel: lib.NetworkChannel}

// network
var productEcommerce = models.Product{Flow: models.ECommerceFlow}
var productMga = models.Product{Flow: models.MgaFlow}
var productProviderMga = models.Product{Flow: models.ProviderMgaFlow}
var productRemittanceMga = models.Product{Flow: models.RemittanceMgaFlow}

func TestEmitForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store)
}

func TestLeadForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)

	exps := []string{
		"setLeadData",
		"sendLeadMail",
	}
	testFlow(t, "lead", exps, store)
}

func TestProposalForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)

	exps := []string{}
	testFlow(t, "proposal", exps, store)
}

func TestPayForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store)
}

func TestSignForEcommerce(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store)
}

func TestEmitForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{}
	testFlow(t, "emit", exps, store)
}

func TestLeadForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store)
}

func TestProposalForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{
		"setProposalData",
	}
	testFlow(t, "proposal", exps, store)
}

func TestApprovalForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{}
	testFlow(t, "requestApproval", exps, store)
}

func TestPayForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store)
}

func TestSignForMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
	}
	testFlow(t, "sign", exps, store)
}

func TestEmitForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store)
}

func TestLeadForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store)
}

func TestProposalForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"setProposalData",
		"sendProposalMail",
	}
	testFlow(t, "proposal", exps, store)
}

func TestApprovalForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store)
}

func TestPayForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	testFlow(t, "pay", exps, store)
}

func TestSignForProviderMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store)
}

func TestEmitForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"sendMailSign",
		"setAdvice",
		"putUser",
	}
	testFlow(t, "emit", exps, store)
}

func TestLeadForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store)
}

func TestProposalForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{
		"setProposalData",
		"sendProposalMail",
	}
	testFlow(t, "proposal", exps, store)
}

func TestApprovalForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store)
}

func TestPayForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{}
	testFlow(t, "pay", exps, store)
}

func TestSignForRemittanceMga(t *testing.T) {
	store := bpnm.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)

	exps := []string{
		"setSign",
		"addContract",
		"sendMailContract",
	}
	testFlow(t, "sign", exps, store)
}
