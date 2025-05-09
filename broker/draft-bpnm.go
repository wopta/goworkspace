package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/broker/internal/handlers"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
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
		injected, err := getNodeFlow()
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
