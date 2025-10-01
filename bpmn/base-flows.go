package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func GetFlow(policy *models.Policy, storage *bpmnEngine.StorageBpnm) (*bpmnEngine.FlowBpnm, error) {
	builder, err := bpmnEngine.NewBpnmBuilder("flows/base-flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlers(builder)
	if err != nil {
		return nil, err
	}

	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	var warrant *models.Warrant
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowNameStr := policy.GetFlow(networkNode, warrant)

	mgaProduct := flow.Product{Product: prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}

	storage.AddGlobal("policy", &flow.Policy{Policy: policy})
	storage.AddGlobal("product", &flow.Product{Product: product})
	storage.AddGlobal("networkNode", &flow.Network{NetworkNode: networkNode})
	storage.AddGlobal("mgaProduct", &mgaProduct)
	storage.AddGlobal("flowName", &flow.String{String: flowNameStr})
	builder.SetStorage(storage)

	err = injectCallbackFlow(networkNode, builder)
	if err != nil {
		return nil, err
	}
	err = injectProductFlow(builder)
	if err != nil {
		return nil, err
	}
	return builder.Build()
}

func addHandlers(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		handlers.AddAcceptanceHandlers(builder),
		handlers.AddEmitHandlers(builder),
		handlers.AddLeadHandlers(builder),
		handlers.AddPayHandlers(builder),
		handlers.AddProposalHandlers(builder),
		handlers.AddRequestApprovaHandlers(builder),
		handlers.AddSignHandlers(builder),
		handlers.AddRecoverHandlers(builder),
	)
}
