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

// Function for exposing the sellability function.
// It takes a request in which the body contains the policy in JSON format,
// and then it returns the appropriate product.
func VehicleHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
		err    error
	)

	log.Println("Vehicle")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "{}", nil, fmt.Errorf("cannot read body: %v", err)
	}
	if err = json.Unmarshal(bytes, &policy); err != nil {
		return "{}", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	product, productJson, err := Vehicle(authToken.Role, policy)

	return productJson, product, err
}

// Given a policy that shoudl contain the Vehicle asset, then it returns:
//   - the product or parts of it depending on the sellable rules
//   - its representation in json
//   - and an eventual error
func Vehicle(role string, p models.Policy) (models.Product, string, error) {
	if err := validatePolicy(p); err != nil {
		return models.Product{}, "", fmt.Errorf("the policy did not pass validation: %v", err)
	}

	if err := isVehicleSellable(p); err != nil {
		return models.Product{}, "", fmt.Errorf("vehicle not sellable: %v", err)
	}

	product, err := productForVehicle(p, role)
	if err != nil {
		return models.Product{}, "", fmt.Errorf("no products for this vehicle: %v", err)
	}

	jsonProduct, err := json.Marshal(product)
	if err != nil {
		return models.Product{}, "", fmt.Errorf("cannot generate the JSON: %v", err)
	}

	return product, string(jsonProduct), nil
}

func productForVehicle(p models.Policy, r string) (models.Product, error) {
	product, err := prd.GetProduct("gap", "v1", r)
	if err != nil {
		return models.Product{}, fmt.Errorf("error in getting the product: %v", err)
	}

	vp := p.Assets[0].Vehicle.PriceValue
	if vp > 100000 && vp <= 120000 {
		delete(product.Offers, "base")
	}

	return product, nil
}

// Returns true if the policy is conforming to the sellability rules
func isVehicleSellable(p models.Policy) error {
	v := p.Assets[0].Vehicle
	if !v.IsFireTheftCovered {
		return fmt.Errorf("fire and theft is not covered")
	}
	if v.MainUse != "private" {
		return fmt.Errorf("the vehicle is not private")
	}

	car_types := []string{"auto", "autocarro", "suv"}
	if !sliceContains(car_types, v.VehicleType) {
		return fmt.Errorf("The vehicle type is not in: %v", car_types)
	}

	anniversary := v.RegistrationDate.AddDate(5, 0, 0)
	if p.StartDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the start date")
	}

	anniversary = v.RegistrationDate.AddDate(8, 0, 0)
	if p.EndDate.After(anniversary) {
		return fmt.Errorf("The registration is too old, exceeded the end date")
	}

	vp := p.Assets[0].Vehicle.PriceValue
	if vp < 4000 || vp > 120000 {
		return fmt.Errorf("the value is not within the accepted range")
	}

	return nil
}

func validatePolicy(p models.Policy) error {
	if len(p.Assets) == 0 {
		return fmt.Errorf("no assets found")
	}
	if p.Assets[0].Vehicle == nil {
		return fmt.Errorf("no vehicle found")
	}

	v := p.Assets[0].Vehicle
	pd := elapsedYears(p.StartDate, p.EndDate)

	thirdBD := v.RegistrationDate.AddDate(3, 0, 0)
	if time.Now().Before(thirdBD) {
		if pd != 5 {
			return fmt.Errorf("wrong policy duration! it should be 5, we've got %d", pd)
		}
	} else {
		decrease := elapsedYears(thirdBD, time.Now())
		coverage := 5 - decrease

		if coverage <= 0 {
			return fmt.Errorf("The coverage has duration 0")
		}
		if pd != coverage {
			return fmt.Errorf("wrong policy duration! it should be %d, we've got %d", coverage, pd)
		}
	}
	return nil
}

func sliceContains[T comparable](s []T, o T) bool {
	for _, e := range s {
		if o == e {
			return true
		}
	}
	return false
}

// computes the age/elapsed years between t1, and t2.
func elapsedYears(t1 time.Time, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}

	t1y, t1m, t1d := t1.Date()
	date1 := time.Date(t1y, t1m, t1d, 0, 0, 0, 0, time.UTC)

	t2y, t2m, t2d := t2.Date()
	date2 := time.Date(t2y, t2m, t2d, 0, 0, 0, 0, time.UTC)

	if date2.Before(date1) {
		return 0
	}

	years := t2y - t1y
	anniversary := date1.AddDate(years, 0, 0)
	if anniversary.After(date2) {
		years--
	}

	return years
}
