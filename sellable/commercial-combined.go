package sellable

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CommercialCombinedFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		policy *models.Policy
		err    error
	)
	log.SetPrefix("[CommercialCombinedFx] ")
	log.Println("Handler start -----------------------------------------------")

	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.Println("error decoding request body",)
		return "", nil, err
	}

	policy.Normalize()

	if err = CommercialCombined(policy); err == nil {
		return "{}", nil, nil
	}

	return "", nil, fmt.Errorf("policy not sellable by: %v", err)
}

type SellableError struct {
	Msg string
}
func (e *SellableError) Error() string {
	return e.Msg
}

func CommercialCombined(p *models.Policy) error {
	in, err := getCommercialCombinedInputData(p)
	if err != nil {
		return err
	}

	rulesFile := lib.GetRulesFileV2(p.Name, p.ProductVersion, rulesFilename)
	fx := new(models.Fx)
	out := new(SellableError)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, out, in, nil)
	out = ruleOutput.(*SellableError)

	if out.Msg == "" {
		return nil
	}

	return out
}

func getCommercialCombinedInputData(policy *models.Policy) ([]byte, error) {
	var numEmp = 0
	var numBuild = 0
	var mandatoryWarrant = false
	var mandatoryThirdParty = false
	var buildingAndRental = false
	var mandatoryWarrantList = []string{"building", "rental-risk", "machinery", "stock"}

	step, err := fromStepStringToInt(policy.Step)
	if err != nil {
		return nil, err
	}

	for _, v := range policy.Assets {
		if v.Type == models.AssetTypeEnterprise {
			numEmp = int(v.Enterprise.Employer)
			for _, g := range v.Guarantees {
				if g.Slug == "third-party-liability-work-providers" {
					mandatoryThirdParty = true
				}
			}
		}
		if v.Type == models.AssetTypeBuilding {
			numBuild++
			var building, rentalRisk = false, false
			for _, g := range v.Guarantees {
				if slices.Contains(mandatoryWarrantList, g.Slug) {
					mandatoryWarrant = true
				}
				if g.Slug == "building" {
					building = true
				}
				if g.Slug == "rental-risk" {
					rentalRisk = true
				}
			}
			if building && rentalRisk {
				buildingAndRental = true
			}
		}
	}

	out := make(map[string]any)
	out["step"] = step
	out["numEmp"] = numEmp
	out["numBuild"] = numBuild
	out["mandatoryWarrant"] = mandatoryWarrant
	out["buildingAndRental"] = buildingAndRental
	out["mandatoryThirdParty"] = mandatoryThirdParty

	output, err := json.Marshal(out)

	return output, err
}

func fromStepStringToInt(step string) (int, error) {
	var (
		intStep int
		err     error = nil
	)

	switch step {
	case "quotereffectivedate":
		intStep = 0
	case "quoterenterprisedata":
		intStep = 1
	case "quoterbuildingdata":
		intStep = 2
	case "quoterclaimshistory":
		intStep = 3
	case "qbeguaranteestep":
		intStep = 4
	case "quoterbondsandclauses":
		intStep = 5
	case "quotersignatorydata":
		intStep = 6
	case "quoterstatements":
		intStep = 7
	default:
		err = fmt.Errorf("unable to parse step string %s", step)
	}

	return intStep, err
}
