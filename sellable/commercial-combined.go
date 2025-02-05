package sellable

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type Out struct {
	Msg string
}

func CommercialCombinedFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		policy *models.Policy
		err    error
	)

	log.Println("[CommercialCombinedFx] handler start ----------- ")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[CommercialCombinedFx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return "", nil, err
	}

	policy.Normalize()

	in, err := getCommercialCombinedInputData(policy)
	if err != nil {
		log.Printf("[Commercial-Combined] error getting input data: %s", err.Error())
		return "", nil, err
	}

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)

	fx := new(models.Fx)

	var out = new(Out)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, out, in, nil)
	out = ruleOutput.(*Out)

	log.Println("[CommercialCombinedFx] handler end ----------------")

	if out.Msg == "" {
		return http.StatusText(http.StatusOK), nil, nil
	}

	return "", nil, fmt.Errorf("policy not sellable by: %v", out.Msg)
}

func getCommercialCombinedInputData(policy *models.Policy) ([]byte, error) {
	var numEmp = 0
	var numBuild = 0
	var missingMandatoryWarrant = false
	var mandatoryThirdParty = false
	var buildingAndRental = false
	var mandatoryWarrantList = []string{"building", "rental-risk", "machinery", "stock"} // "third-party-liability-work-providers"

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
				if !slices.Contains(mandatoryWarrantList, g.Slug) {
					missingMandatoryWarrant = true
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
	out["missingMandatoryWarrant"] = missingMandatoryWarrant
	out["buildingAndRental"] = buildingAndRental
	out["mandatoryThirdParty"] = mandatoryThirdParty

	output, err := json.Marshal(out)

	return output, err
}
