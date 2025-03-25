package dto

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type numeric struct {
	ValueFloat float64
	ValueInt   int64
	Text       string
}

func (n *numeric) FromValue(v any) {
	switch v := v.(type) {
	case float64:
		n.ValueFloat = v
		n.Text = lib.HumanaizePriceEuro(v)
	case int64:
		n.ValueInt = v
		n.Text = strconv.FormatInt(v, 10)
	}
}

func newNumeric() numeric {
	return numeric{
		ValueFloat: 0,
		ValueInt:   0,
		Text:       constants.EmptyField,
	}
}

type guaranteeDTO struct {
	Description                string
	SumInsuredLimitOfIndemnity numeric
	LimitOfIndemnity           numeric
	SumInsured                 numeric
	StartDate                  string
	Duration                   numeric
	RetroactiveDate            string
	RetroactiveDateUsa         string
}

func newGuaranteeDTO() *guaranteeDTO {
	return &guaranteeDTO{
		Description:                constants.EmptyField,
		SumInsuredLimitOfIndemnity: newNumeric(),
		LimitOfIndemnity:           newNumeric(),
		SumInsured:                 newNumeric(),
		StartDate:                  constants.EmptyField,
		Duration:                   newNumeric(),
		RetroactiveDate:            constants.EmptyField,
		RetroactiveDateUsa:         constants.EmptyField,
	}
}

func (g *guaranteeDTO) fromPolicy(guarantee models.Guarante) {
	g.SumInsuredLimitOfIndemnity.ValueFloat = guarantee.Value.SumInsuredLimitOfIndemnity
	if g.SumInsuredLimitOfIndemnity.ValueFloat != 0 {
		g.SumInsuredLimitOfIndemnity.Text = lib.HumanaizePriceEuro(g.SumInsuredLimitOfIndemnity.ValueFloat)
	}
	g.LimitOfIndemnity.ValueFloat = guarantee.Value.LimitOfIndemnity
	if g.LimitOfIndemnity.ValueFloat != 0 {
		g.LimitOfIndemnity.Text = lib.HumanaizePriceEuro(g.LimitOfIndemnity.ValueFloat)
	}
	g.SumInsured.ValueFloat = guarantee.Value.SumInsured
	if g.SumInsured.ValueFloat != 0 {
		g.SumInsured.Text = lib.HumanaizePriceEuro(g.SumInsured.ValueFloat)
	}
	if guarantee.Value.StartDate != nil && !guarantee.Value.StartDate.IsZero() {
		g.StartDate = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
	}
	if guarantee.Value.Duration != nil {
		g.Duration.ValueInt = int64(guarantee.Value.Duration.Day)
		g.Duration.Text = fmt.Sprintf("%d", g.Duration.ValueInt)
	}
	if guarantee.Value.RetroactiveDate != nil && !guarantee.Value.RetroactiveDate.IsZero() {
		g.RetroactiveDate = guarantee.Value.RetroactiveDate.Format(constants.DayMonthYearFormat)
	}
	if guarantee.Value.RetroactiveUsaCanDate != nil && !guarantee.Value.RetroactiveDate.IsZero() {
		g.RetroactiveDateUsa = guarantee.Value.RetroactiveUsaCanDate.Format(constants.DayMonthYearFormat)
	}
}

type quoteGuaranteeDTO struct {
	Description                string
	SumInsuredLimitOfIndemnity numeric
	ExpiryDate                 string
	Duration                   numeric
	PremiumGrossYearly         numeric
}

func newQuoteGuaranteeDTO() *quoteGuaranteeDTO {
	return &quoteGuaranteeDTO{
		Description:                constants.EmptyField,
		SumInsuredLimitOfIndemnity: newNumeric(),
		Duration:                   newNumeric(),
		PremiumGrossYearly:         newNumeric(),
		ExpiryDate:                 constants.EmptyField,
	}
}

func (g *quoteGuaranteeDTO) fromData(guarantee models.Guarante, startDate time.Time) {
	g.SumInsuredLimitOfIndemnity.ValueFloat = guarantee.Value.SumInsuredLimitOfIndemnity
	if g.SumInsuredLimitOfIndemnity.ValueFloat != 0 {
		g.SumInsuredLimitOfIndemnity.Text = lib.HumanaizePriceEuro(g.SumInsuredLimitOfIndemnity.ValueFloat)
	}

	if guarantee.Value.Duration != nil {
		g.Duration.ValueInt = int64(guarantee.Value.Duration.Year)
		g.Duration.Text = fmt.Sprintf("%d", g.Duration.ValueInt)
	}

	g.PremiumGrossYearly.ValueFloat = guarantee.Value.PremiumGrossYearly
	if g.PremiumGrossYearly.ValueFloat != 0 {
		g.PremiumGrossYearly.Text = lib.HumanaizePriceEuro(g.PremiumGrossYearly.ValueFloat)
	}

	g.ExpiryDate = startDate.AddDate(int(g.Duration.ValueInt), 0,0).Format(constants.DayMonthYearFormat)
}
