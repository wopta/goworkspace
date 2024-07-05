package callback_out

import (
	"log"

	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/models"
)

var (
	Proposal        internal.CallbackoutAction = internal.ExtProposal
	RequestApproval internal.CallbackoutAction = internal.ExtRequestApproval
	Emit            internal.CallbackoutAction = internal.ExtEmit
	Paid            internal.CallbackoutAction = internal.ExtPaid
	EmitRemittance  internal.CallbackoutAction = internal.ExtEmitRemittance
)

func Execute(node *models.NetworkNode, policy models.Policy, rawAction internal.CallbackoutAction) {
	var (
		client CallbackClient
		err    error
		fx     func(models.Policy) internal.CallbackInfo
	)

	if node == nil || node.CallbackConfig == nil {
		log.Println("no node or callback config available")
		return
	}

	if client, err = newClient(node); err != nil {
		log.Println(err)
		return
	}

	actions := client.DecodeAction(rawAction)

	for _, action := range actions {
		switch action {
		case Proposal:
			fx = client.Proposal
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

		res := fx(policy)
		log.Printf("Callback request: %v", res.Request)
		log.Printf("Callback response: %v", res.Response)
		log.Printf("Callback error: %s", res.Error)

		saveAudit(node, action, res)
	}
}
