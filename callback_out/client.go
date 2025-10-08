package callback_out

import (
	"encoding/json"
	"errors"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CallbackClient interface {
	Proposal(models.Policy) base.CallbackInfo
	Emit(models.Policy) base.CallbackInfo
	Signed(models.Policy) base.CallbackInfo
	Paid(models.Policy) base.CallbackInfo
	RequestApproval(models.Policy) base.CallbackInfo
	Approved(models.Policy) base.CallbackInfo
	Rejected(models.Policy) base.CallbackInfo
}

type CallbackConfig struct {
	Proposal        bool `json:"proposal"`
	RequestApproval bool `json:"requestApproval"`
	Signed          bool `json:"signed"`
	Paid            bool `json:"paid"`
	Emit            bool `json:"emit"`
	EmitRemittance  bool `json:"emitRemittance"`
	Approved        bool `json:"approved"`
	Rejected        bool `json:"rejected"`
}

var ErrCallbackClientNotSet = errors.New("callback client not set")

func NewClient(node *models.NetworkNode) (client CallbackClient, conf CallbackConfig, err error) {
	var bytes []byte
	log.WarningF("The callbacks accept only policy with name life")
	switch node.CallbackConfig.Name {
	case "winClient":
		client = win.NewClient(node.ExternalNetworkCode)
		bytes, err = lib.GetFilesByEnvV2("flows/callback/win.json")
	case "facileBrokerClient":
		client = base.NewClient(node, "facile_broker")
		bytes, err = lib.GetFilesByEnvV2("flows/callback/base.json")
	default:
		log.WarningF("Callback client isn't correct %s", node.CallbackConfig.Name)
		err = ErrCallbackClientNotSet
	}
	if err != nil {
		return client, conf, err
	}
	err = json.Unmarshal(bytes, &conf)
	return client, conf, err
}
