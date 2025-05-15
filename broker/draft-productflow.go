package broker

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/quote/catnat"
)

func getProductFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/product_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("catnatIntegration", catnatIntegration),
	)
	if err != nil {
		return nil, err
	}
	return builder, nil
}

func catnatIntegration(store bpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", store)
	if err != nil {
		return err
	}
	client := catnat.NewNetClient()
	requestDto := catnat.RequestDTO{}
	err = requestDto.FromPolicy(policy.Policy, true)
	if err != nil {
		return err
	}
	log.PrintStruct("-----------------------dto for quote", requestDto)
	res, e := client.Emit(requestDto)
	if e != nil {
		return e
	}
	log.PrintStruct("emit catnat response", res)
	return nil
}
