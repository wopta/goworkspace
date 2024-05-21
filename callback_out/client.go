package callback_out

import (
	"errors"
	"net/http"

	"github.com/wopta/goworkspace/callback_out/win"
	"github.com/wopta/goworkspace/models"
)

type CallbackClient interface {
	// Proposal(models.Policy) (*http.Request, *http.Response, error)
	Emit(models.Policy) (*http.Request, *http.Response, error)
	// Signed(models.Policy) (*http.Request, *http.Response, error)
	Paid(models.Policy) (*http.Request, *http.Response, error)
	RequestApproval(models.Policy) (*http.Request, *http.Response, error)
	// Approved(models.Policy) (*http.Request, *http.Response, error)
	// Rejected(models.Policy) (*http.Request, *http.Response, error)
}

var ErrCallbackClientNotSet = errors.New("callback client not set")

func newClient(node *models.NetworkNode) (CallbackClient, error) {
	switch node.CallbackConfig.Name {
	case "winClient":
		return win.NewClient(), nil
	default:
		return nil, ErrCallbackClientNotSet
	}
}