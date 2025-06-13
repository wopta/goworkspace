package bpmn

import (
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/bpmn/internal/handlers"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func GetFlow(policy *models.Policy, originStr string, storage *bpmnEngine.StorageBpnm) (*bpmnEngine.FlowBpnm, error) {
	builder, err := bpmnEngine.NewBpnmBuilder("flows/draft/base-flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlersDraft(builder)
	if err != nil {
		return nil, err
	}

	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	var warrant *models.Warrant
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowNameStr, _ := policy.GetFlow(networkNode, warrant)

	mgaProduct := flow.Product{Product: prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}

	storage.AddGlobal("policy", &flow.Policy{Policy: policy})
	storage.AddGlobal("product", &flow.Product{Product: product})
	storage.AddGlobal("networkNode", &flow.Network{NetworkNode: networkNode})
	storage.AddGlobal("mgaProduct", &mgaProduct)
	storage.AddGlobal("flowName", &flow.String{String: flowNameStr})
	storage.AddGlobal("origin", &flow.String{String: originStr})
	builder.SetStorage(storage)

	if networkNode != nil && networkNode.CallbackConfig != nil {
		injected, err := getNodeFlow(networkNode)
		if err != nil {
			return nil, err
		}
		err = builder.Inject(injected)
		if err != nil {
			return nil, err
		}
	} else {
		log.InfoF("no node or callback config available, no callback")
	}
	injected, err := getProductFlow()
	if err != nil {
		return nil, err
	}
	err = builder.Inject(injected)
	if err != nil {
		return nil, err
	}
	return builder.Build()
}

func addHandlersDraft(builder *bpmnEngine.BpnmBuilder) error {
	return bpmnEngine.IsError(
		handlers.AddAcceptanceHandlers(builder),
		handlers.AddEmitHandlers(builder),
		handlers.AddLeadHandlers(builder),
		handlers.AddPayHandlers(builder),
		handlers.AddProposalHandlers(builder),
		handlers.AddRequestApprovaHandlers(builder),
		handlers.AddSignHandlers(builder),
	)
}
