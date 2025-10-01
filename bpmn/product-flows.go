package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers/productFlow"
)

func injectProductFlow(mainBuilder *bpmnEngine.BpnmBuilder) error {
	store := bpmnEngine.NewStorageBpnm()
	builder, e := bpmnEngine.NewBpnmBuilder("flows/product-flows.json")
	if e != nil {
		return e
	}
	builder.SetStorage(store)
	err := bpmnEngine.IsError(
		builder.AddHandler("catnatIntegration", productFlow.CatnatIntegration),
		builder.AddHandler("catnatdownloadPolicy", productFlow.CatnatDownloadCertification),
	)
	if err != nil {
		return err
	}
	return mainBuilder.Inject(builder)
}
