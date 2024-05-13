package callback_out

import (
	"errors"
	"log"

	"github.com/wopta/goworkspace/models"
)

var ErrUnhandledStatus = errors.New("unhandled callback status")

func winCallbackHandler(policy models.Policy) error {
	log.Println("executing win callback handler...")
	switch policy.Status {
	case models.PolicyStatusProposal:
		return winProposalCallback(policy)
	case models.PolicyStatusWaitForApproval, models.PolicyStatusWaitForApprovalMga:
		return winWaitForApprovalCallback(policy)
	}

	if policy.IsPay {
		return winContractCallback(policy)
	}

	if policy.CompanyEmit {
		return winEmitCallback(policy)
	}

	log.Printf("status '%s' not handled by callback", policy.Status)
	return ErrUnhandledStatus
}

func winProposalCallback(policy models.Policy) error {
	log.Println("win proposal calback...")
	return nil
}

func winWaitForApprovalCallback(policy models.Policy) error {
	log.Println("win wait for approval calback...")
	return nil
}

func winEmitCallback(policy models.Policy) error {
	log.Println("win emit calback...")
	return nil
}

func winContractCallback(policy models.Policy) error {
	log.Println("win contract calback...")
	return nil
}
