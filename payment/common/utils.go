package common

import (
	"fmt"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CheckPaymentModes(policy models.Policy) error {
	var allowedModes []string

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		allowedModes = models.GetAllowedMonthlyModes()
	case string(models.PaySplitYearly):
		allowedModes = models.GetAllowedYearlyModes()
	case string(models.PaySplitSingleInstallment):
		allowedModes = models.GetAllowedSingleInstallmentModes()
	}

	if !lib.SliceContains(allowedModes, policy.PaymentMode) {
		return fmt.Errorf("mode '%s' is incompatible with split '%s'", policy.PaymentMode, policy.PaymentSplit)
	}

	return nil
}
