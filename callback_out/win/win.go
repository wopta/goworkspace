package win

import (
	"errors"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
)

var ErrUnhandledStatus = errors.New("unhandled callback status")

func CallbackHandler(policy models.Policy) error {
	log.Println("executing win callback handler...")

	var fx func(models.Policy) (*http.Request, *http.Response, error)

	switch policy.Status {
	case models.PolicyStatusProposal:
		fx = proposalCallback
	case models.PolicyStatusWaitForApproval, models.PolicyStatusWaitForApprovalMga:
		fx = approvalCallback
	}
	if policy.IsPay {
		fx = contractCallback
	}
	if policy.CompanyEmit {
		fx = emitCallback
	}

	if fx == nil {
		log.Printf("status '%s' not handled by callback", policy.Status)
		return ErrUnhandledStatus
	}

	req, res, err := fx(policy)
	if err != nil {
		// TODO: save error somewhere
		log.Printf("Callback error: %s", err.Error())
		log.Printf("Callback request: %v", req)
		log.Printf("Callback response: %v", res)
	}

	return err
}
