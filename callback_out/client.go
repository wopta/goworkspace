package callback_out

import (
	"errors"

	"github.com/wopta/goworkspace/callback_out/base"
	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/callback_out/win"
	"github.com/wopta/goworkspace/models"
)

type CallbackClient interface {
	Proposal(models.Policy) internal.CallbackInfo
	Emit(models.Policy) internal.CallbackInfo
	Signed(models.Policy) internal.CallbackInfo
	Paid(models.Policy) internal.CallbackInfo
	RequestApproval(models.Policy) internal.CallbackInfo
	Approved(models.Policy) internal.CallbackInfo
	Rejected(models.Policy) internal.CallbackInfo

	// This method is temporary while we do not settle on the config for the node
	DecodeAction(string) []string
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
