package draftbpnm

import (
	"testing"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func funcTest(message string, log *mockLog) func(StorageData) error {
	return func(sd StorageData) error {
		log.Println(message)
		return nil
	}
}
func getBuilder(log *mockLog, store StorageData) (*FlowBpnm, error) {

	builder, e := NewBpnmBuilder("flows.json")
	if e != nil {
		return nil, e
	}

	builder.SetStorage(store)

	builder.AddHandler("setProposalData", funcTest("setProposalData", log))
	builder.AddHandler("emitData", funcTest("emitData", log))
	builder.AddHandler("sendMailSign", funcTest("sendMailSign", log))
	builder.AddHandler("pay", funcTest("pay", log))
	builder.AddHandler("setAdvice", funcTest("setAdvice", log))
	builder.AddHandler("putUser", funcTest("putUser", log))
	builder.AddHandler("sendEmitProposalMail", funcTest("sendEmitProposalMail", log))
	builder.AddHandler("setLeadData", funcTest("setLeadData", log))
	builder.AddHandler("sendLeadMail", funcTest("sendLeadMail", log))
	builder.AddHandler("updatePolicy", funcTest("updatePolicy", log))
	builder.AddHandler("sign", funcTest("sign", log))
	builder.AddHandler("payTransaction", funcTest("payTransaction", log))
	builder.AddHandler("sendProposalMail", funcTest("sendProposalMail", log))
	builder.AddHandler("fillAttachments", funcTest("fillAttachments", log))
	builder.AddHandler("setToPay", funcTest("setToPay", log))
	builder.AddHandler("setSign", funcTest("setSign", log))
	builder.AddHandler("sendMailContract", funcTest("sendMailContract", log))
	builder.AddHandler("sendMailPay", funcTest("sendMailPay", log))
	return builder.Build()
}

func testFlow(t *testing.T, process string, expectedACtivities []string, store StorageData) {
	log := mockLog{}
	flow, e := getBuilder(&log, store)
	if e != nil {
		t.Fatal(e)
	}
	if err := flow.Run(process); err != nil {
		t.Fatal(err)
	}
	if len(expectedACtivities) != len(log.log) {

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
var productProviderMga = models.Product{Flow: models.ProviderMgaFlow}

func TestEmitForEcommerce(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productProviderMga)
	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	t.Helper()
	testFlow(t, "emit", exps, store)
}
func TestLeadForEcommerce(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)

	exps := []string{
		"setLeadData",
		"sendLeadMail",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	t.Helper()
	testFlow(t, "lead", exps, store)
}
func TestProposalForEcommerce(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)

	exps := []string{}
	t.Helper()
	testFlow(t, "proposal", exps, store)
}
func TestSignForEcommerce(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	t.Helper()
	testFlow(t, "sign", exps, store)
}

func TestPayForEcommerce(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)

	exps := []string{
		"updatePolicy",
		"payTransaction",
		"sign",
	}
	t.Helper()
	testFlow(t, "pay", exps, store)
}

func TestEmitForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productProviderMga)

	exps := []string{}
	t.Helper()
	testFlow(t, "emit", exps, store)
}
func TestLeadForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)

	exps := []string{
		"setLeadData",
	}
	t.Helper()
	testFlow(t, "lead", exps, store)
}
func TestProposalForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)

	exps := []string{
		"setProposalData",
	}
	t.Helper()
	testFlow(t, "proposal", exps, store)
}
func TestApprovalForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)

	exps := []string{}
	t.Helper()
	testFlow(t, "requestApproval", exps, store)
}
func TestPayForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)

	exps := []string{
		"updatePolicy",
		"payTransaction",
	}
	t.Helper()
	testFlow(t, "pay", exps, store)
}
func TestSignForMga(t *testing.T) {
	store := NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
	}
	t.Helper()
	testFlow(t, "sign", exps, store)
}
