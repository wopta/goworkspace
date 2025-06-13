package flow

import (
	"os"
	"testing"

	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func mock_catnatIntegration(log *mockLog) func(bpmnEngine.StorageData) error {
	return func(sd bpmnEngine.StorageData) error {
		log.println("catnatIntegration")
		sd.AddLocal("numeroPolizza", &String{String: "provissiamo"})
		return nil
	}
}

func getBuilderFlowProduct(log *mockLog, store bpmnEngine.StorageData) (*bpmnEngine.BpnmBuilder, error) {
	os.Setenv("env", env.LocalTest)
	builder, e := bpmnEngine.NewBpnmBuilderRawPath("../../../../function-data/dev/flows/draft/product-flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpmnEngine.IsError(
		builder.AddHandler("catnatIntegration", mock_catnatIntegration(log)),
		builder.AddHandler("catnatdownloadPolicy", funcTest("catnatdownloadPolicy", log)),
	)
	if e != nil {
		return nil, e
	}
	return builder, nil
}
func TestCatnatIntegrationWithNoCatnatPolicy(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &Policy{&models.Policy{Name: "noCatnat"}})
	store.AddGlobal("networkNode", &winNode)

	exps := []string{}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "catnatIntegration", exps, store, getBuilderFlowProduct)
}

func TestCatnatIntegrationWithCatnatPolicy(t *testing.T) {
	store := bpmnEngine.NewStorageBpnm()
	store.AddGlobal("policy", &policyCatnat)
	store.AddGlobal("networkNode", &winNode)

	exps := []string{
		"catnatIntegration",
		"catnatdownloadPolicy",
	}
	store.AddLocal("config", &CallbackConfig{Emit: false})

	testFlow(t, "catnatIntegration", exps, store, getBuilderFlowProduct)
}
