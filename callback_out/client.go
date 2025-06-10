package callback_out

import (
	"errors"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/callback_out/win"
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

	// This method is temporary while we do not settle on the config for the node
	DecodeAction(base.CallbackoutAction) []base.CallbackoutAction
}

var ErrCallbackClientNotSet = errors.New("callback client not set")

func newClient(node *models.NetworkNode) (CallbackClient, error) {
	switch node.CallbackConfig.Name {
	case "winClient":
		return win.NewClient(node.ExternalNetworkCode), nil
	case "facileBrokerClient":
		return base.NewClient(node, "facile_broker"), nil
	default:
		return nil, ErrCallbackClientNotSet
	}
}
