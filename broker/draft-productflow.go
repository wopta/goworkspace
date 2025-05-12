package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models/dto/net"
)

func getProductFlow() (*bpmn.BpnmBuilder, error) {
	store := bpmn.NewStorageBpnm()
	builder, e := bpmn.NewBpnmBuilder("flows/draft/node_flows.json")
	if e != nil {
		return nil, e
	}
	builder.SetStorage(store)
	err := bpmn.IsError(
		builder.AddHandler("catnatIntegration", callBackApproved),
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
	var url string
	var cnReq net.RequestDTO
	err = cnReq.FromPolicy(policy.Policy, false) //TODO:to se
	if err != nil {
		return err
	}
	rBuff := new(bytes.Buffer)
	err = json.NewEncoder(rBuff).Encode(cnReq)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, url, rBuff)
	req.Header.Set("Content-Type", "application/json")
	res, err := lib.RetryDo(req, 5, 100)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("Error integration with net-ensure")
	}
	return nil
}
