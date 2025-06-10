package callback_out

import (
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func Execute(node *models.NetworkNode, policy models.Policy, rawAction base.CallbackoutAction) {
	var (
		client CallbackClient
		err    error
		fx     func(models.Policy) base.CallbackInfo
	)

	if node == nil || node.CallbackConfig == nil {
		log.Println("no node or callback config available")
		return
	}

	if client, err = newClient(node); err != nil {
		log.Error(err)
		return
	}

	if client == nil {
		log.ErrorF("client not found")
		return
	}

	actions := client.DecodeAction(rawAction)

	if len(actions) == 0 {
		log.Printf("action '%s' not implemented for client", rawAction)
		return
	}

	for _, action := range actions {
		switch action {
		case base.Proposal:
			fx = client.Proposal
		case base.RequestApproval:
			fx = client.RequestApproval
		case base.Emit:
			fx = client.Emit
		case base.Signed:
			fx = client.Signed
		case base.Paid:
			fx = client.Paid
		case base.Approved:
			fx = client.Approved
		case base.Rejected:
			fx = client.Rejected
		default:
			log.Printf("unhandled callback action '%s'", action)
			return
		}

		log.Printf("executing action '%s'", action)

		res := fx(policy)
		log.Printf("Callback error: %s", res.Error)

		saveAudit(node, res)
	}
}
