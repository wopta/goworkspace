package handlers

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers/productFlow"
)

func AddProductsHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		productFlow.AddCatnatHandlers(builder),
	)
}
