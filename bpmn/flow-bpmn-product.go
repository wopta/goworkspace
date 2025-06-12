package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers/productFlow"
)

func getProductFlow() (*bpmnEngine.BpnmBuilder, error) {
	store := bpmnEngine.NewStorageBpnm()
	builder, e := bpmnEngine.NewBpnmBuilder("flows/draft/product_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmnEngine.IsError(
		builder.AddHandler("catnatIntegration", productFlow.CatnatIntegration),
		builder.AddHandler("catnatdownloadPolicy", productFlow.CatnatDownloadCertification),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
