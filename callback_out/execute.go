package callback_out

import (
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
)

type CallbackoutAction = string

var (
	RequestApproval CallbackoutAction = "RequestApproval"
	Emit            CallbackoutAction = "Emit"
	Paid            CallbackoutAction = "Paid"
)

func Execute(node *models.NetworkNode, policy models.Policy, action CallbackoutAction) {
	var (
		client CallbackClient
		err    error
		fx     func(models.Policy) (*http.Request, *http.Response, error)
	)

	if node == nil || node.CallbackConfig == nil {
		log.Println("no node or callback config available")
		return
	}

	if client, err = newClient(node); err != nil {
		log.Println(err)
		return
	}

	switch action {
	case RequestApproval:
		fx = client.RequestApproval
	case Emit:
		fx = client.Emit
	case Paid:
		fx = client.Paid
	default:
		log.Printf("unhandled callback action '%s'", action)
		return
	}

	log.Printf("executing action '%s'", action)

	req, res, err := fx(policy)
	log.Printf("Callback request: %v", req)
	log.Printf("Callback response: %v", res)
	log.Printf("Callback error: %s", err)

	saveAudit(node, action, req, res)
}
