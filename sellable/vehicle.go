package sellable

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
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

// Function for exposing the sellability function of Gap (vehicles).
// It takes a request in which the body contains the policy in JSON format,
// which should contain a vehicle asset and a person asset
// and then it returns the appropriate product.
func GapFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
		err    error
	)

	log.Println("[GapFx] Handler start")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "{}", nil, fmt.Errorf("cannot read body: %v", err)
	}
	if err = json.Unmarshal(bytes, &policy); err != nil {
		return "{}", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	product, err := Gap(authToken.Role, &policy)
	if err != nil {
		return "", models.Product{}, fmt.Errorf("cannot retrieve the product: %v", err)
	}

	jsonProduct, err := json.Marshal(product)
	if err != nil {
		return "", models.Product{}, fmt.Errorf("cannot generate the JSON: %v", err)
	}

	return string(jsonProduct), product, err
}

// Given a policy that should contain the Gap and the Person assets, then it returns:
//   - the product or parts of it depending on the sellable rules
//   - and an eventual error
func Gap(role string, p *models.Policy) (models.Product, error) {
	if err := validatePolicy(p); err != nil {
		return models.Product{}, fmt.Errorf("the policy did not pass validation: %v", err)
	}

	if err := isVehicleSellable(p); err != nil {
		return models.Product{}, fmt.Errorf("vehicle not sellable: %v", err)
	}

	product, err := productForVehicle(p, role)
	if err != nil {
		return models.Product{}, fmt.Errorf("no products for this vehicle: %v", err)
	}

	return product, nil
}

func productForVehicle(p *models.Policy, r string) (models.Product, error) {
	product, err := prd.GetProduct("gap", "v1", r)
	if err != nil {
		return models.Product{}, fmt.Errorf("error in getting the product: %v", err)
	}

	vp := p.Assets[0].Vehicle.PriceValue
	if vp > minPriceOnlyComplete && vp <= maxPriceValue {
		delete(product.Offers, "base")
	}

	return product, nil
}

// Returns true if the policy is conforming to the sellability rules for GAP
func isVehicleSellable(p *models.Policy) error {
	v := p.Assets[0].Vehicle
	if !v.IsFireTheftCovered {
		return fmt.Errorf("fire and theft is not covered")
	}
	if v.MainUse != "private" {
		return fmt.Errorf("the vehicle is not private")
	}

	carTypes := []string{"auto", "autocarro", "suv"}
	if !lib.SliceContains(carTypes, v.VehicleType) {
		return fmt.Errorf("The vehicle type is not in: %v", carTypes)
	}

	anniversary := v.RegistrationDate.AddDate(maxAgeAtStartDate, 0, 0)
	if p.StartDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the start date")
	}

	anniversary = v.RegistrationDate.AddDate(maxAgeAtEndDate, 0, 0)
	if p.EndDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the end date")
	}

	vp := p.Assets[0].Vehicle.PriceValue
	if vp < minPriceValue || vp > maxPriceValue {
		return fmt.Errorf("the value is not within the accepted range")
	}

	return nil
}

func validatePolicy(p *models.Policy) error {
	if len(p.Assets) == 0 {
		return fmt.Errorf("no assets found")
	}

	if p.Assets[0].Person == nil {
		return fmt.Errorf("no person found")
	}

	if p.Assets[0].Vehicle == nil {
		return fmt.Errorf("no vehicle found")
	}

	v := p.Assets[0].Vehicle
	pd := lib.ElapsedYears(p.StartDate, p.EndDate)

	maxAgeFullCoverageBD := v.RegistrationDate.AddDate(maxAgeFullCoverage, 0, 0)
	if time.Now().Before(maxAgeFullCoverageBD) {
		if pd <= maxCoverage {
			return fmt.Errorf(
				"wrong policy duration! It should be at maximum %d, we've got %d",
				maxCoverage,
				pd,
			)
		}
	} else {
		decrease := lib.ElapsedYears(maxAgeFullCoverageBD, time.Now())
		coverage := maxCoverage - decrease

		if coverage <= 0 {
			return fmt.Errorf("The coverage has duration 0")
		}
		if pd <= coverage {
			return fmt.Errorf("wrong policy duration! it should be at most %d, we've got %d", coverage, pd)
		}
	}
	return nil
}
