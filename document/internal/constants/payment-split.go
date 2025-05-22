package constants

import "gitlab.dev.wopta.it/goworkspace/models"

const (
	paymentSplitYearly  = "Annuale"
	paymentSplitMonthly = "Mensile"
)

var (
	PaymentSplitMap = map[string]string{
		string(models.PaySplitYearly):  paymentSplitYearly,
		string(models.PaySplitMonthly): paymentSplitMonthly,
	}
)
