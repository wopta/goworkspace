package broker

import (
	"encoding/base64"
	"os"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
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
		builder.AddHandler("catnatdownloadPolicy", catnatDownloadCertification),
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
	res, e := client.Emit(requestDto)
	if e != nil {
		return e
	}
	store.AddLocal("numeroPoliza", &flow.StringBpmn{String: res.PolicyNumber})
	return nil
}

func catnatDownloadCertification(store bpmn.StorageData) error {
	var policy *flow.PolicyDraft
	var numeroPoliza *flow.StringBpmn
	err := bpmn.IsError(
		bpmn.GetDataRef("policy", &policy, store),
		bpmn.GetDataRef("numeroPoliza", &numeroPoliza, store),
	)
	if err != nil {
		return err
	}

	client := catnat.NewNetClient()
	resp, err := client.Download(numeroPoliza.String)
	if err != nil {
		return err
	}
	bytes, err := base64.StdEncoding.DecodeString(resp.Documento[0].DatiDocumento)
	if err != nil {
		return err
	}
	os.WriteFile("prova.pdf", bytes, 0644)

	return nil
}
