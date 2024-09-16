package reserved

import (
	"log"
	"math"

	"github.com/wopta/goworkspace/models"
)

func personaReserved(policy *models.Policy) (bool, *models.ReservedInfo) {
	log.Println("[personaReserved]")

	const reason = "BMI index out of range"

	reservedInfo := &models.ReservedInfo{
		Reasons:       make([]string, 0),
		RequiredExams: make([]string, 0),
	}

	isReserved := false

	voidPolicy := len(policy.Assets) == 0 || policy.Assets[0].Person == nil
	voidData := policy.Assets[0].Person.Weight == 0 || policy.Assets[0].Person.Height == 0
	if voidPolicy || voidData {
		return isReserved, reservedInfo
	}

	_, isOutOfRange := checkOutOfRangeBMI(policy.Assets[0].Person.Weight, policy.Assets[0].Person.Height)

	if isOutOfRange {
		isReserved = true
		reservedInfo.Reasons = append(reservedInfo.Reasons, reason)
	}

	return isReserved, reservedInfo
}

func checkOutOfRangeBMI(weight int, height int) (float64, bool) {
	/*
		Numbers have no exact representation in binary floating point.
		So instead of comparing values, we compare their difference against a very small
		constant (called epsilon here)
	*/
	const epsilon = 1e-9
	const bmiLowerLimit = 16
	const bmiUpperLimit = 40
	w := float64(weight)
	h := float64(height)
	bmi := w / math.Pow(h/100, 2)
	isOutOfRange := false
	if (bmi-bmiLowerLimit) < epsilon || (bmi-bmiUpperLimit) > epsilon {
		isOutOfRange = true
	}

	return bmi, isOutOfRange
}
