package dto

import (
	"strconv"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type contractDTO struct {
	CodeHeading  string
	Code         string
	StartDate    string
	EndDate      string
	PaymentSplit string
	NextPay      string
	HasBond      string
	BondText     string
}

func newContractDTO() *contractDTO {
	return &contractDTO{
		CodeHeading:  emptyField,
		Code:         emptyField,
		StartDate:    emptyField,
		EndDate:      emptyField,
		PaymentSplit: emptyField,
		NextPay:      emptyField,
		HasBond:      no,
	}

}

func (c *contractDTO) fromPolicy(policy models.Policy, isProposal bool) {
	c.CodeHeading = "I dati della tua Proposta nr."
	c.Code = strconv.Itoa(policy.ProposalNumber)

	if !isProposal {
		c.CodeHeading = "I dati della tua Polizza nr."
		if policy.CodeCompany != "" {
			c.Code = policy.CodeCompany
		} else {
			c.Code = emptyField
		}
	}

	if !policy.StartDate.IsZero() {
		c.StartDate = policy.StartDate.Format(constants.DayMonthYearFormat)
	}

	if !policy.EndDate.IsZero() {
		c.EndDate = policy.EndDate.Format(constants.DayMonthYearFormat)
		c.NextPay = c.EndDate
	}

	if _, ok := constants.PaymentSplitMap[policy.PaymentSplit]; ok {
		c.PaymentSplit = constants.PaymentSplitMap[policy.PaymentSplit]
	}

	nextPayDate := lib.AddMonths(policy.StartDate, 12*policy.Annuity)
	if !nextPayDate.After(policy.EndDate) {
		c.NextPay = nextPayDate.Format(constants.DayMonthYearFormat)
	}

	if policy.HasBond {
		c.HasBond = yes
		c.BondText = policy.Bond
	}
}
