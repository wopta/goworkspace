package sellable

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wopta/goworkspace/lib"
)

// Strategy for SellableWrapper that takes in account the policy's
// contractor birth date
type ByContractorAge struct{}

func (*ByContractorAge) ExtractParams(w *SellableWrapper) ([]byte, error) {
	age, err := w.input.CalculateContractorAge()
	if err != nil {
		return nil, err
	}

	out := make(map[string]int)
	out["age"] = age

	output, err := json.Marshal(out)

	return output, err
}

// Not implemented - age verification is handled elsewhere (life/persona)
func (*ByContractorAge) Validate(*SellableWrapper) error {
	return fmt.Errorf("ByContractorAge.Validate not implemented")
}

// Strategy for SellableWrapper that takes in account the policy's
// vehicle registration date
type ByVehicleRegistrationDate struct{}

// Not implemented - strategy does not use rules engine
func (*ByVehicleRegistrationDate) ExtractParams(w *SellableWrapper) ([]byte, error) {
	return nil, fmt.Errorf("ByVehicleRegistrationDate.ExtractParams not implemented")
}

func (*ByVehicleRegistrationDate) Validate(w *SellableWrapper) error {
	policy := w.input
	if len(policy.Assets) == 0 {
		return fmt.Errorf("no assets found")
	}

	// TODO: move me - not on scope
	if policy.Assets[0].Person == nil {
		return fmt.Errorf("no person found")
	}

	if policy.Assets[0].Vehicle == nil {
		return fmt.Errorf("no vehicle found")
	}

	vehicle := policy.Assets[0].Vehicle
	policyDuration := lib.ElapsedYears(policy.StartDate, policy.EndDate)

	maxAgeFullCoverageBD := vehicle.RegistrationDate.AddDate(maxAgeFullCoverage, 0, 0)
	if time.Now().Before(maxAgeFullCoverageBD) {
		if policyDuration > maxCoverage {
			return fmt.Errorf(
				"wrong policy duration! It should be at maximum %d, we've got %d",
				maxCoverage,
				policyDuration,
			)
		}
	} else {
		decrease := lib.ElapsedYears(maxAgeFullCoverageBD, time.Now())
		coverage := maxCoverage - decrease

		if coverage <= 0 {
			return fmt.Errorf("the coverage has duration 0")
		}
		if policyDuration > coverage {
			return fmt.Errorf("wrong policy duration! it should be at most %d, we've got %d", coverage, policyDuration)
		}
	}
	return nil
}
