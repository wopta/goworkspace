package draftbpnm

import (
	"testing"

	"github.com/wopta/goworkspace/models"
)

func funcTest(message string, log *mockLog) func(StorageData) error {
	return func(sd StorageData) error {
		log.Println(message)
		return nil
	}
}
func getBuilder(log *mockLog) (*FlowBpnm, error) {
	storage := NewStorageBpnm()
	p := new(models.Policy)
	pr := new(models.Product)
	n := new(models.NetworkNode)

	storage.AddGlobal("product", pr)
	storage.AddGlobal("mgaProduct", pr)
	storage.AddGlobal("policy", p)
	storage.AddGlobal("network", n)
	builder, e := NewBpnmBuilder("e-commerce.json")
	if e != nil {
		return nil, e
	}
	builder.SetPoolDate(storage)

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
	return builder.Build()
}
func TestEmit(t *testing.T) {
	log := mockLog{}
	flow, e := getBuilder(&log)
	if e != nil {
		t.Fatal(e)
	}
	flow.Run("emit")
	for _, m := range log.log {
		t.Log(m)
	}

	exps := []string{
		"setProposalData",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}
func TestLead(t *testing.T) {
	log := mockLog{}
	flow, e := getBuilder(&log)
	if e != nil {
		t.Fatal(e)
	}
	flow.Run("lead")
	for _, m := range log.log {
		t.Log(m)
	}

	exps := []string{
		"setLeadData",
		"sendLeadMail",
		"emitData",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}
func TestPay(t *testing.T) {
	log := mockLog{}
	flow, e := getBuilder(&log)
	if e != nil {
		t.Fatal(e)
	}
	flow.Run("pay")
	for _, m := range log.log {
		t.Log(m)
	}

	exps := []string{
		"updatePolicy",
		"payTransaction",
		"sign",
	}
	if len(exps) != len(log.log) {
		t.Fatalf("exp n message: %v,got: %v", len(exps), len(log.log))
	}
	for i, exp := range exps {
		if log.log[i] != exp {
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}
