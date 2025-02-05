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
	log.SetPrefix("[CommercialCombinedFx]")
	log.Println("handler start ----------- ")

	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ----------------------------------------------")
		log.SetPrefix("")
	}()

	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.Printf("error decoding request body: %s", err)
		return "", nil, err
	}

	policy.Normalize()

	in, err := getCommercialCombinedInputData(policy)
	if err != nil {
		log.Printf("error getting input data: %s", err.Error())
		return "", nil, err
	}

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)

	fx := new(models.Fx)

	type Out struct {
		Msg string
	}
	var out = new(Out)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, out, in, nil)
	out = ruleOutput.(*Out)

	if out.Msg == "" {
		return http.StatusText(http.StatusOK), nil, nil
	}

	return "", nil, fmt.Errorf("policy not sellable by: %v", out.Msg)
}

func getCommercialCombinedInputData(policy *models.Policy) ([]byte, error) {
	var numEmp = 0
	var numBuild = 0
	var mandatoryWarrant = false
	var mandatoryThirdParty = false
	var buildingAndRental = false
	var mandatoryWarrantList = []string{"building", "rental-risk", "machinery", "stock"}

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
	out["numEmp"] = numEmp
	out["numBuild"] = numBuild
	out["mandatoryWarrant"] = mandatoryWarrant
	out["buildingAndRental"] = buildingAndRental
	out["mandatoryThirdParty"] = mandatoryThirdParty

	output, err := json.Marshal(out)

	return output, err
}
