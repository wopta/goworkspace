package callback_out

import (
	"log"
	"net/http"
	"slices"

	"github.com/wopta/goworkspace/models"
)

func Execute(node *models.NetworkNode, policy models.Policy) {
	var (
		client CallbackClient
		err    error
		fx     func(models.Policy) (*http.Request, *http.Response, error)
	)

	if node == nil || node.CallbackConfig == nil {
		return
	}

	if client, err = newClient(node); err != nil {
		return
	}

	if slices.Contains([]string{models.PolicyStatusWaitForApproval, models.PolicyStatusWaitForApprovalMga}, policy.Status) {
		fx = client.RequestApproval
	}
	if policy.IsPay {
		fx = client.Paid
	}
	if policy.CompanyEmit {
		fx = client.Emit
	}

	if fx == nil {
		log.Printf("status '%s' not handled by callback", policy.Status)
		return
	}

	req, res, err := fx(policy)
	log.Printf("Callback request: %v", req)
	log.Printf("Callback response: %v", res)
	log.Printf("Callback error: %s", err)
	// TODO: save
}
