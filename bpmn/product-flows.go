package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers"
)

func injectProductFlow(mainBuilder *bpmnEngine.BpnmBuilder) error {
	store := bpmnEngine.NewStorageBpnm()
	builder, e := bpmnEngine.NewBpnmBuilder("flows/product-flows.json")
	if e != nil {
		return e
	}
	builder.SetStorage(store)
	err := bpmnEngine.IsError(
		handlers.AddProductsHandlers(builder),
	)
	if err != nil {
		return err
	}
	return mainBuilder.Inject(builder)
}
