package broker

import (
	"encoding/json"
	"errors"
	"fmt"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/models/client"
	"github.com/wopta/goworkspace/models/dto/net"
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
	client := client.NewNetClient()
	client.Authenticate()
	requestDto := net.RequestDTO{}
	requestDto.FromPolicy(policy.Policy, true)
	requestDto.Emission = "si"
	res, _, e := client.Quote(requestDto)
	if e != nil {
		return e
	}
	if len(res.Errors) != 0 {
		return errors.New(fmt.Sprintln(res.Errors))
	}
	var resString []byte
	resString, _ = json.Marshal(res)
	println(string(resString))
	return nil
}
