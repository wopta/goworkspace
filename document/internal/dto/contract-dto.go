package dto

import (
	"strconv"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
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
	Clause       string
}

func newContractDTO() *contractDTO {
	return &contractDTO{
		CodeHeading:  constants.EmptyField,
		Code:         constants.EmptyField,
		StartDate:    constants.EmptyField,
		EndDate:      constants.EmptyField,
		PaymentSplit: constants.EmptyField,
		NextPay:      constants.EmptyField,
		HasBond:      constants.No,
		BondText:     constants.EmptyField,
		Clause:       constants.EmptyField,
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
			c.Code = constants.EmptyField
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
		c.HasBond = constants.Yes
		c.BondText = policy.Bond
	}

	if policy.Clause != "" {
		c.Clause = policy.Clause
	}
}
