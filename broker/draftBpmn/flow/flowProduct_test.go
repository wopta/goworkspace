package flow

import (
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	env "github.com/wopta/goworkspace/lib/environment"
	"os"
)

func mock_catnatIntegration(log *mockLog) func(bpmn.StorageData) error {
	return func(sd bpmn.StorageData) error {
		log.println("catnatIntegration")
		sd.AddLocal("numeroPolizza", &StringBpmn{String: "provissiamo"})
		return nil
	}
}

func getBuilderFlowProduct(log *mockLog, store bpmn.StorageData) (*bpmn.BpnmBuilder, error) {
	os.Setenv("env", env.LocalTest)
	builder, e := bpmn.NewBpnmBuilderRawPath("../../../../function-data/dev/flows/draft/product_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	e = bpmn.IsError(
		builder.AddHandler("catnatIntegration", mock_catnatIntegration(log)),
		builder.AddHandler("catnatdownloadPolicy", funcTest("catnatdownloadPolicy", log)),
	)
	if e != nil {
		return nil, e
	}
	return builder, nil
}
