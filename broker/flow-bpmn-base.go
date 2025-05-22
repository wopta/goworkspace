package broker

import (
	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/handlers"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

func getFlow(policy *models.Policy, originStr string, storage bpmn.StorageData) (*bpmn.FlowBpnm, error) {
	builder, err := bpmn.NewBpnmBuilder("flows/draft/channel_flows.json")
	if err != nil {
		return nil, err
	}
	err = addHandlersDraft(builder)
	if err != nil {
		return nil, err
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	var warrant *models.Warrant
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, networkNode, warrant)
	flowNameStr, _ := policy.GetFlow(networkNode, warrant)

	policyDraft := flow.PolicyDraft{Policy: policy}
	productDraft := flow.ProductDraft{Product: product}
	mgaProduct := flow.ProductDraft{Product: prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)}
	flowName := flow.StringBpmn{String: flowNameStr}
	networkDraft := flow.NetworkDraft{NetworkNode: networkNode}
	if product.Flow == "" {
		product.Flow = policy.Channel
	}

	storage.AddGlobal("policy", &policyDraft)
	storage.AddGlobal("product", &productDraft)
	storage.AddGlobal("networkNode", &networkDraft)
	storage.AddGlobal("mgaProduct", &mgaProduct)
	storage.AddGlobal("flowName", &flowName)
	storage.AddGlobal("origin", &flow.StringBpmn{String: originStr})
	builder.SetStorage(storage)

	if networkNode != nil && networkNode.CallbackConfig != nil {
		injected, err := getNodeFlow(networkNode.CallbackConfig.Name)
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

func addHandlersDraft(builder *bpmn.BpnmBuilder) error {
	return bpmn.IsError(
		handlers.AddAcceptanceHandlers(builder),
		handlers.AddEmitHandlers(builder),
		handlers.AddLeadHandlers(builder),
		handlers.AddPayHandlers(builder),
		handlers.AddProposalHandlers(builder),
		handlers.AddRequestApprovaHandlers(builder),
		handlers.AddSignHandlers(builder),
	)
}
