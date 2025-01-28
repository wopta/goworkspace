package constants

import "github.com/wopta/goworkspace/models"

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
