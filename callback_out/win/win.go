package win

import (
	"errors"
	"log"

	"github.com/wopta/goworkspace/models"
)

var ErrUnhandledStatus = errors.New("unhandled callback status")

func CallbackHandler(policy models.Policy) error {
	log.Println("executing win callback handler...")
	switch policy.Status {
	case models.PolicyStatusProposal:
		return proposalCallback(policy)
	case models.PolicyStatusWaitForApproval, models.PolicyStatusWaitForApprovalMga:
		return approvalCallback(policy)
	}

	if policy.IsPay {
		return contractCallback(policy)
	}

	if policy.CompanyEmit {
		return emitCallback(policy)
	}

	log.Printf("status '%s' not handled by callback", policy.Status)
	return ErrUnhandledStatus
}
