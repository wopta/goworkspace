package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/internal/handlers/productFlow"
)

func getProductFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/product_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("catnatIntegration", productFlow.CatnatIntegration),
		builder.AddHandler("catnatdownloadPolicy", productFlow.CatnatDownloadCertification),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}
