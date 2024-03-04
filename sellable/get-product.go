package sellable

import (
	"fmt"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
)

type SellableWrapper struct {
	input       *models.Policy
	baseProduct *models.Product
	params      sellableParameters
	fx          func(*SellableWrapper) (*models.Product, error)
}

func (w *SellableWrapper) Evaluate() (*models.Product, error) {
	return w.fx(w)
}

type sellableParameters interface {
	ExtractParams(*SellableWrapper) ([]byte, error)
	Validate(*SellableWrapper) error
}

func GetProduct(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (product *models.Product, err error) {
	var wrapper *SellableWrapper

	// load base product
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		return nil, fmt.Errorf("no product found")
	}

	// build strategy by product
	switch policy.Name {
	case models.LifeProduct:
		wrapper = &SellableWrapper{
			input:       policy,
			params:      &ByContractorAge{},
			baseProduct: product,
			fx:          sellableLife,
		}
	case models.GapProduct:
		wrapper = &SellableWrapper{
			input:       policy,
			params:      &ByVehicleRegistrationDate{},
			baseProduct: product,
			fx:          sellableGap,
		}
	default:
		return nil, fmt.Errorf("product not implemented: '%s'", policy.Name)
	}

	return wrapper.Evaluate()
}

func sellableLife(w *SellableWrapper) (*models.Product, error) {
	rulesFile := lib.GetRulesFileV2(w.input.Name, w.input.ProductVersion, rulesFilename)

	in, err := w.params.ExtractParams(w)
	if err != nil {
		return nil, err
	}

	fx := new(models.Fx)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, w.baseProduct, in, nil)

	return ruleOutput.(*models.Product), nil
}

func sellableGap(w *SellableWrapper) (*models.Product, error) {
	err := isGapPolicySellable(w.input)
	if err != nil {
		return nil, err
	}

	vehiclePrice := w.input.Assets[0].Vehicle.PriceValue
	if vehiclePrice > minPriceOnlyComplete && vehiclePrice <= maxPriceValue {
		delete(w.baseProduct.Offers, "base")
	}

	return w.baseProduct, nil
}

func isGapPolicySellable(policy *models.Policy) error {
	vehicle := policy.Assets[0].Vehicle
	if !vehicle.IsFireTheftCovered {
		return fmt.Errorf("fire and theft is not covered")
	}
	if vehicle.MainUse != "private" {
		return fmt.Errorf("the vehicle is not private")
	}

	vehicleTypes := []string{"car", "truck", "camper"}
	if !lib.SliceContains(vehicleTypes, strings.ToLower(vehicle.VehicleType)) {
		return fmt.Errorf("the vehicle type is not in: %v", vehicleTypes)
	}

	anniversary := vehicle.RegistrationDate.AddDate(maxAgeAtStartDate, 0, 0)
	if policy.StartDate.After(anniversary) {
		return fmt.Errorf("the registration is too old, exceeded the start date")
	}

	anniversary = vehicle.RegistrationDate.AddDate(maxAgeAtEndDate, 0, 0)
	if policy.EndDate.After(anniversary) {
		return fmt.Errorf("the registration is too old, exceeded the end date")
	}

	vehiclePrice := policy.Assets[0].Vehicle.PriceValue
	if vehiclePrice < minPriceValue || vehiclePrice > maxPriceValue {
		return fmt.Errorf("the value is not within the accepted range")
	}

	return nil
}
