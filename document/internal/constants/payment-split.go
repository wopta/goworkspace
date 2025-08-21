package constants

import "gitlab.dev.wopta.it/goworkspace/models"

const (
	paymentSplitYearly     = "Annuale"
	paymentSplitMonthly    = "Mensile"
	paymentSplitSemestral  = "Semestrale"
	paymentSplitSingle     = "Singolo"
	paymentSplitTrimestral = "Trimestrale"
)

var (
	PaymentSplitMap = map[string]string{
		string(models.PaySplitYearly):            paymentSplitYearly,
		string(models.PaySplitYear):              paymentSplitYearly,
		string(models.PaySplitMonthly):           paymentSplitMonthly,
		string(models.PaySplitSemestral):         paymentSplitSemestral,
		string(models.PaySplitSingleInstallment): paymentSplitSingle,
		string(models.PaySplitTrimestral):        paymentSplitTrimestral,
	}
)
