package flow

import (
	"testing"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type mockLog struct {
	log []string
}

func (m *mockLog) println(mes string) {
	m.log = append(m.log, mes)
}

func (m *mockLog) printlnForTesting(t *testing.T) {
	t.Log("Actual log: ")
	for _, mes := range m.log {
		t.Log(" ", mes)
	}
}
func funcTest(message string, log *mockLog) func(*bpmnEngine.StorageBpnm) error {
	return func(sd *bpmnEngine.StorageBpnm) error {
		log.println(message)
		return nil
	}
}

func getBuilderFlowChannel(log *mockLog, store *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
	builder, e := bpmnEngine.NewBpnmBuilderRawPath("../../../../function-data/dev//flows/draft/base-flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpmnEngine.IsError(
		builder.AddHandler("setProposalData", funcTest("setProposalData", log)),
		builder.AddHandler("emitWithSequence", funcTest("emitWithSequence", log)),
		builder.AddHandler("emitNoSequence", funcTest("emitNoSequence", log)),
		builder.AddHandler("sendMailSign", funcTest("sendMailSign", log)),
		builder.AddHandler("pay", funcTest("pay", log)),
		builder.AddHandler("setAdvice", funcTest("setAdvice", log)),
		builder.AddHandler("putUser", funcTest("putUser", log)),
		builder.AddHandler("sendEmitProposalMail", funcTest("sendEmitProposalMail", log)),
		builder.AddHandler("setLeadData", funcTest("setLeadData", log)),
		builder.AddHandler("sendLeadMail", funcTest("sendLeadMail", log)),
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
		builder.AddHandler("rejected", funcTest("rejected", log)),
		builder.AddHandler("approved", funcTest("approved", log)),
		builder.AddHandler("sendAcceptanceMail", funcTest("sendAcceptanceMail", log)),
		builder.AddHandler("generateInvoice", funcTest("generateInvoice", log)),
		builder.AddHandler("updatePolicyAsPaid", funcTest("updatePolicyAsPaid", log)),
		builder.AddHandler("promotePolicy", funcTest("promotePolicy", log)),
		builder.AddHandler("saveTransactionAndPolicy", funcTest("saveTransactionAndPolicy", log)),
		builder.AddHandler("createNetworkTransaction", funcTest("createNetworkTransaction", log)),
	)
	if e != nil {
		return nil, e
	}

	return builder, nil
}

type builderFlow func(*mockLog, *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error)

func testFlow(t *testing.T, process string, expectedACtivities []string, store *bpmnEngine.StorageBpnm, builbuilderFlow builderFlow) {
	log := mockLog{}
	build, e := builbuilderFlow(&log, store)
	if e != nil {
		t.Fatal(e)
	}
	flow, e := build.Build()
	if e != nil {
		t.Fatal(e)
	}
	if err := flow.Run(bpmnEngine.BpmnFlow(process)); err != nil {
		t.Fatal(err)
	}
	if len(expectedACtivities) != len(log.log) {
		log.printlnForTesting(t)
		t.Fatalf("exp n message: %v,got: %v", len(expectedACtivities), len(log.log))
	}
	for i, exp := range expectedACtivities {
		if log.log[i] != exp {
			log.printlnForTesting(t)
			t.Fatalf("exp: %v,got: %v", exp, log.log[i])
		}
	}
}

var (
	policyEcommerce = Policy{&models.Policy{Channel: lib.ECommerceChannel, Name: "test policy", PaymentMode: models.PaymentModeRecurrent}}
	policyMga       = Policy{&models.Policy{Channel: lib.MgaChannel, Name: "test policy"}}
	policyNetwork   = Policy{&models.Policy{Channel: lib.NetworkChannel, Name: "test policy"}}
	policyCatnat    = Policy{&models.Policy{Channel: lib.NetworkChannel, Name: models.CatNatProduct}}
)

var (
	paymentInfo = PaymentInfoBpmn{}
	transaction = Transaction{}
	addresses   = Addresses{}
)

// product
var (
	productEcommerce     = Product{&models.Product{Flow: models.ECommerceFlow}}
	productMga           = Product{&models.Product{Flow: models.MgaFlow}}
	productProviderMga   = Product{&models.Product{Flow: models.ProviderMgaFlow}}
	productRemittanceMga = Product{&models.Product{Flow: models.RemittanceMgaFlow}}
)

func initBaseStorage(storage *bpmnEngine.StorageBpnm) {
	storage.AddGlobal("mgaProduct", &Product{})
	storage.AddGlobal("flowName", &String{})
	storage.AddGlobal("addresses", &addresses)
	storage.AddGlobal("origin", &String{})
	storage.AddGlobal("paymentSplit", &String{})
	storage.AddGlobal("paymentMode", &String{})
	storage.AddGlobal("sendEmail", &BoolBpmn{})
	storage.AddGlobal("clientCallback", &callbackClient)
}

func TestEmitForEcommerceForCatnat(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setProposalData",
		"emitWithSequence",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}
func TestLeadForEcommerce(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setLeadData",
		"sendLeadMail",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForEcommerce(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("is_PROPOSAL_V2", &BoolBpmn{true})
	initBaseStorage(store)

	exps := []string{}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestPayAnnuity0FirstRate(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	policyEcommerce.Annuity = 0
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("paymentInfo", &paymentInfo)
	store.AddGlobal("transaction", &transaction)
	store.AddGlobal("skipInvoice", &BoolBpmn{false})
	store.AddGlobal("sendEmail", &BoolBpmn{true})
	initBaseStorage(store)

	exps := []string{
		"payTransaction",
		"promotePolicy",
		"updatePolicyAsPaid",
		"generateInvoice",
		"createNetworkTransaction",
		"saveTransactionAndPolicy",
		"sendMailContract",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}
func TestPayAnnuityNot0FirstRate(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	initBaseStorage(store)
	policyEcommerce.Annuity = 1
	store.AddGlobal("sendEmail", &BoolBpmn{false})
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("paymentInfo", &paymentInfo)
	store.AddGlobal("transaction", &transaction)
	store.AddGlobal("skipInvoice", &BoolBpmn{true})

	exps := []string{
		"payTransaction",
		"updatePolicyAsPaid",
		"createNetworkTransaction",
		"saveTransactionAndPolicy",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestPaySingleRate(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	policyEcommerce.PaymentMode = models.PaymentModeSingle
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("paymentInfo", &paymentInfo)
	store.AddGlobal("transaction", &transaction)
	store.AddGlobal("skipInvoice", &BoolBpmn{false})
	initBaseStorage(store)

	exps := []string{
		"payTransaction",
		"updatePolicyAsPaid",
		"generateInvoice",
		"createNetworkTransaction",
		"saveTransactionAndPolicy",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestPayAnnuityNoFirstRateNoFirstPayment(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	initBaseStorage(store)
	policyEcommerce.Annuity = 1
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("paymentInfo", &paymentInfo)
	store.AddGlobal("transaction", &transaction)
	store.AddGlobal("sendEmail", &BoolBpmn{false})
	store.AddGlobal("skipInvoice", &BoolBpmn{false})

	exps := []string{
		"payTransaction",
		"updatePolicyAsPaid",
		"generateInvoice",
		"createNetworkTransaction",
		"saveTransactionAndPolicy",
	}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}
func TestSignForEcommerce(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productEcommerce)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("is_PROPOSAL_V2", &BoolBpmn{true})
	initBaseStorage(store)

	exps := []string{
		"setProposalData",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestSignForMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyMga)
	store.AddGlobal("product", &productMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForProviderMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setProposalData",
		"emitWithSequence",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForProviderMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForProviderMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("is_PROPOSAL_V2", &BoolBpmn{true})
	initBaseStorage(store)

	exps := []string{
		"setProposalData",
		"sendProposalMail",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForProviderMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestSignForProviderMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyEcommerce)
	store.AddGlobal("product", &productProviderMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"fillAttachments",
		"setSign",
		"setToPay",
		"sendMailPay",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

func TestEmitForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setProposalData",
		"emitWithSequence",
		"sign",
		"sendMailSign",
		"setAdvice",
		"putUser",
	}
	testFlow(t, "emit", exps, store, getBuilderFlowChannel)
}

func TestLeadForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "lead", exps, store, getBuilderFlowChannel)
}

func TestProposalForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("is_PROPOSAL_V2", &BoolBpmn{false})
	initBaseStorage(store)

	exps := []string{
		"setLeadData",
	}
	testFlow(t, "proposal", exps, store, getBuilderFlowChannel)
}

func TestApprovalForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setRequestApprovalData",
		"sendRequestApprovalMail",
	}
	testFlow(t, "requestApproval", exps, store, getBuilderFlowChannel)
}

func TestPayForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	store.AddGlobal("paymentInfo", &paymentInfo)
	store.AddGlobal("transaction", &transaction)
	store.AddGlobal("skipInvoice", &BoolBpmn{true})
	initBaseStorage(store)

	exps := []string{}
	testFlow(t, "pay", exps, store, getBuilderFlowChannel)
}

func TestSignForRemittanceMga(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyNetwork)
	store.AddGlobal("product", &productRemittanceMga)
	store.AddGlobal("networkNode", &winNode)
	initBaseStorage(store)

	exps := []string{
		"setSign",
		"addContract",
		"sendMailContract",
	}
	testFlow(t, "sign", exps, store, getBuilderFlowChannel)
}

// With node flow
func TestSignForRemittanceMgaWithNodeFlow(t *testing.T) {
	storeFlowChannel := bpmnEngine.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyNetwork)
	storeFlowChannel.AddGlobal("product", &productRemittanceMga)
	storeFlowChannel.AddGlobal("networkNode", &winNode)
	initBaseStorage(storeFlowChannel)

	storeNode := bpmnEngine.NewStorageBpnm()
	storeNode.AddLocal("config", &CallbackConfig{Sign: false})

	exps := []string{
		"setSign",
		"addContract",
		"sendMailContract",
	}
	testFlow(t, "sign", exps, storeFlowChannel, func(log *mockLog, sd *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
		build, e := getBuilderFlowChannel(log, storeFlowChannel)
		if e != nil {
			return nil, e
		}
		nodeBuild, e := getBuilderFlowNode(log, storeNode)
		if e != nil {
			return nil, e
		}
		if e := build.Inject(nodeBuild); e != nil {
			t.Fatal(e)
		}

		return build, nil
	})
}

func TestEmitForEcommerceWithNodeFlow(t *testing.T) {
	storeFlowChannel := bpmnEngine.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyEcommerce)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("networkNode", &winNode)
	initBaseStorage(storeFlowChannel)

	storeNode := bpmnEngine.NewStorageBpnm()
	storeNode.AddLocal("config", &CallbackConfig{Emit: true})

	exps := []string{
		"setProposalData",
		"emitWithSequence",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
		"Emit",
		"saveAudit prova request Emit",
	}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
		build, e := getBuilderFlowChannel(log, storeFlowChannel)
		if e != nil {
			return nil, e
		}
		nodeBuild, e := getBuilderFlowNode(log, storeNode)
		if e != nil {
			return nil, e
		}
		if e := build.Inject(nodeBuild); e != nil {
			return nil, e
		}

		return build, nil
	})
}
func TestEmitForWgaWithNodeFlow(t *testing.T) {
	storeFlowChannel := bpmnEngine.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyMga)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("networkNode", &winNode)
	initBaseStorage(storeFlowChannel)

	storeNode := bpmnEngine.NewStorageBpnm()
	storeNode.AddLocal("config", &CallbackConfig{Emit: true})

	exps := []string{}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
		build, e := getBuilderFlowChannel(log, storeFlowChannel)
		if e != nil {
			return nil, e
		}
		nodeBuild, e := getBuilderFlowNode(log, storeNode)

		if e != nil {
			return nil, e
		}
		if e := build.Inject(nodeBuild); e != nil {
			t.Fatal(e)
		}

		return build, nil
	})
}
func TestEmitForEcommerceWithNodeFlowConfFalse(t *testing.T) {
	storeFlowChannel := bpmnEngine.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyEcommerce)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("networkNode", &winNode)
	initBaseStorage(storeFlowChannel)

	storeNode := bpmnEngine.NewStorageBpnm()
	storeNode.AddLocal("config", &CallbackConfig{Emit: false})

	storeProduct := bpmnEngine.NewStorageBpnm()
	exps := []string{
		"setProposalData",
		"emitWithSequence",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
		build, e := getBuilderFlowChannel(log, storeFlowChannel)
		if e != nil {
			return nil, e
		}
		nodeBuild, e := getBuilderFlowNode(log, storeNode)
		if e != nil {
			return nil, e
		}
		productBuild, e := getBuilderFlowProduct(log, storeProduct)
		if e != nil {
			return nil, e
		}
		if e := build.Inject(nodeBuild); e != nil {
			return nil, e
		}
		if e := build.Inject(productBuild); e != nil {
			return nil, e
		}

		return build, nil
	})
}
func TestEmitForEcommerceCatnat(t *testing.T) {
	storeFlowChannel := bpmnEngine.NewStorageBpnm()
	storeFlowChannel.AddGlobal("policy", &policyCatnat)
	storeFlowChannel.AddGlobal("product", &productEcommerce)
	storeFlowChannel.AddGlobal("networkNode", &winNode)
	initBaseStorage(storeFlowChannel)

	storeNode := bpmnEngine.NewStorageBpnm()
	storeNode.AddLocal("config", &CallbackConfig{Emit: false})

	storeProduct := bpmnEngine.NewStorageBpnm()
	exps := []string{
		"setProposalData",
		"emitNoSequence",
		"catnatIntegration",
		"catnatdownloadPolicy",
		"sign",
		"pay",
		"sendEmitProposalMail",
		"sendMailSign",
	}
	testFlow(t, "emit", exps, storeFlowChannel, func(log *mockLog, sd *bpmnEngine.StorageBpnm) (*bpmnEngine.BpnmBuilder, error) {
		build, e := getBuilderFlowChannel(log, storeFlowChannel)
		if e != nil {
			return nil, e
		}
		nodeBuild, e := getBuilderFlowNode(log, storeNode)
		if e != nil {
			return nil, e
		}
		productBuild, e := getBuilderFlowProduct(log, storeProduct)
		if e != nil {
			return nil, e
		}
		if e := build.Inject(nodeBuild); e != nil {
			return nil, e
		}
		if e := build.Inject(productBuild); e != nil {
			return nil, e
		}

		return build, nil
	})
}
