package sellable

import (
	"fmt"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

const (
	minPriceValue        = 4000
	maxPriceValue        = 120000
	minPriceOnlyComplete = 100000
	maxAgeAtStartDate    = 5
	maxAgeAtEndDate      = 8
	maxAgeFullCoverage   = 3
	maxCoverage          = 5
)

// Given a policy that should contain the Gap and the Persona assets, then it returns:
//   - the product or parts of it depending on the sellability rules
//   - and an eventual error
func Gap(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Product, error) {
	log.AddPrefix("Gap")
	defer log.PopPrefix()
	log.Println("function start ---------------")

	log.Println("validating policy")

	if err := validatePolicy(policy); err != nil {
		log.ErrorF("error validating policy: %s", err.Error())
		return nil, fmt.Errorf("the policy did not pass validation: %v", err)
	}

	log.Println("loading product file")

	product, err := getProduct(policy, channel, networkNode, warrant)
	if err != nil {
		log.ErrorF("error loading product: %s", err.Error())
		return nil, fmt.Errorf("no products for this vehicle: %v", err)
	}

	log.Println("check policy vendibility")

	if err := isVehicleSellable(policy); err != nil {
		log.ErrorF("error check policy vendility: %s", err.Error())
		return nil, fmt.Errorf("vehicle not sellable: %v", err)
	}

	log.Println("function end ---------------")

	return product, nil
}

func getProduct(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Product, error) {
	product := prd.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		return nil, fmt.Errorf("no product found")
	}

	vehiclePrice := policy.Assets[0].Vehicle.PriceValue
	if vehiclePrice > minPriceOnlyComplete && vehiclePrice <= maxPriceValue {
		delete(product.Offers, "base")
	}

	return product, nil
}

// Returns nil if the policy is eligible for GAP, otherwise returns an error describing why the vehicle is not sellable
func isVehicleSellable(policy *models.Policy) error {
	vehicle := policy.Assets[0].Vehicle
	if !vehicle.IsFireTheftCovered {
		return fmt.Errorf("fire and theft is not covered")
	}
	if vehicle.MainUse != "private" {
		return fmt.Errorf("the vehicle is not private")
	}

	vehicleTypes := []string{"car", "truck", "camper"}
	if !lib.SliceContains(vehicleTypes, strings.ToLower(vehicle.VehicleType)) {
		return fmt.Errorf("The vehicle type is not in: %v", vehicleTypes)
	}

	anniversary := vehicle.RegistrationDate.AddDate(maxAgeAtStartDate, 0, 0)
	if policy.StartDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the start date")
	}

	anniversary = vehicle.RegistrationDate.AddDate(maxAgeAtEndDate, 0, 0)
	if policy.EndDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the end date")
	}

	vehiclePrice := policy.Assets[0].Vehicle.PriceValue
	if vehiclePrice < minPriceValue || vehiclePrice > maxPriceValue {
		return fmt.Errorf("the value is not within the accepted range")
	}

	return nil
}

func validatePolicy(policy *models.Policy) error {
	if len(policy.Assets) == 0 {
		return fmt.Errorf("no assets found")
	}

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
			return fmt.Errorf("The coverage has duration 0")
		}
		if policyDuration > coverage {
			return fmt.Errorf("wrong policy duration! it should be at most %d, we've got %d", coverage, policyDuration)
		}
	}
	return nil
}
