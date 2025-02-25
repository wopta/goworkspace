package quote_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/quote"
)

type InputData struct {
	HasDependants bool   `json:"hasDependants"`
	LifeRisk      int    `json:"lifeRisk"`
	FinancialRisk int    `json:"financialRisk"`
	Age           int    `json:"age"`
	WorkType      string `json:"workType"`
	RiskClass     string `json:"riskClass"`
}

type OutputGuaranteeData struct {
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
	Deductible                 string  `json:"deductible"`
	DeductibleType             string  `json:"deductibleType"`
	PriceGross                 float64 `json:"priceGross"`
}

type OutputOfferData struct {
	IPI        *OutputGuaranteeData `json:"IPI,omitempty"`
	D          *OutputGuaranteeData `json:"D,omitempty"`
	DRG        *OutputGuaranteeData `json:"DRG,omitempty"`
	DC         *OutputGuaranteeData `json:"DC,omitempty"`
	RSC        *OutputGuaranteeData `json:"RSC,omitempty"`
	ITI        *OutputGuaranteeData `json:"ITI,omitempty"`
	PriceGross float64              `json:"priceGross"`
}

type OutputData = map[string]OutputOfferData

type TestData struct {
	Name   string     `json:"name"`
	Input  InputData  `json:"input"`
	Output OutputData `json:"output"`
}

const filename = "data/test/quote/persona.json"

func TestPersona(t *testing.T) {
	env := os.Getenv("ENV")
	folder := "../../function-data/dev/"

	if env == "ci" {
		dir, _ := os.Getwd()
		folder = dir + "/" + folder
	}

	t.Setenv("env", "local-test")

	fileReader, err := os.Open(folder + filename)
	if err != nil {
		t.Fatalf("unable to load data from %s: %s", folder, err)
	}

	testData := make([]TestData, 0)

	if err := json.NewDecoder(fileReader).Decode(&testData); err != nil {
		t.Fatalf("unable to decode data: %s", err)
	}

	for idx, data := range testData {
		if len(data.Output) == 0 {
			continue
		}

		p := buildPolicy(data.Input)
		if err := quote.Persona(&p, models.ECommerceChannel, nil, nil, models.ECommerceFlow); err != nil {
			t.Fatalf("error quoting test %d: %s", idx+1, err)
		}

		var (
			numOffersExpected int
			numOffersGot      = len(p.OffersPrices)
		)

		for offerName, value := range data.Output {
			numOffersExpected++
			var (
				mismatchedPrice    = false
				offerPrice         = 0.0
				guaranteesExpected = make([]string, 0)
				guaranteesGot      = make([]string, 0)
			)

			if a, ok := p.OffersPrices[offerName]; ok {
				if v, ok := a["yearly"]; ok {
					offerPrice = v.Gross
				}
			}
			if offerPrice != value.PriceGross {
				mismatchedPrice = true
				t.Errorf("%s - mismatched offer price for %s. Expected %.2f - Got %.2f", data.Name, offerName, value.PriceGross, offerPrice)
			}
			if value.IPI != nil {
				g, _ := p.ExtractGuarantee("IPI")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for IPI guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.IPI.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - IPI sum. Expected %.2f - Got %.2f", data.Name, offerName, value.IPI.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if g.Offer[offerName].Deductible != value.IPI.Deductible {
					t.Errorf("%s - mismatched offer %s - IPI deductible. Expected %s - Got %s", data.Name, offerName, value.IPI.Deductible, g.Offer[offerName].Deductible)
				}
				if g.Offer[offerName].DeductibleType != value.IPI.DeductibleType {
					t.Errorf("%s - mismatched offer %s - IPI deductibleType. Expected %s - Got %s", data.Name, offerName, value.IPI.DeductibleType, g.Offer[offerName].DeductibleType)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.IPI.PriceGross {
					t.Errorf("%s - mismatched offer %s - IPI price. Expected %.2f - Got %.2f", data.Name, offerName, value.IPI.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}
			if value.D != nil {
				g, _ := p.ExtractGuarantee("D")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for D guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.D.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - D sum. Expected %.2f - Got %.2f", data.Name, offerName, value.D.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.D.PriceGross {
					t.Errorf("%s - mismatched offer %s - D price. Expected %.2f - Got %.2f", data.Name, offerName, value.D.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}
			if value.DRG != nil {
				g, _ := p.ExtractGuarantee("DRG")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for DRG guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.DRG.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - DRG sum. Expected %.2f - Got %.2f", data.Name, offerName, value.DRG.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.DRG.PriceGross {
					t.Errorf("%s - mismatched offer %s - DRG price. Expected %.2f - Got %.2f", data.Name, offerName, value.DRG.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}
			if value.DC != nil {
				g, _ := p.ExtractGuarantee("DC")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for DC guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.DC.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - DC sum. Expected %.2f - Got %.2f", data.Name, offerName, value.DC.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.DC.PriceGross {
					t.Errorf("%s - mismatched offer %s - DC price. Expected %.2f - Got %.2f", data.Name, offerName, value.DC.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}
			if value.RSC != nil {
				g, _ := p.ExtractGuarantee("RSC")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for RSC guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.RSC.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - RSC sum. Expected %.2f - Got %.2f", data.Name, offerName, value.RSC.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.RSC.PriceGross {
					t.Errorf("%s - mismatched offer %s - RSC price. Expected %.2f - Got %.2f", data.Name, offerName, value.RSC.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}
			if value.ITI != nil {
				g, _ := p.ExtractGuarantee("ITI")
				if _, ok := g.Offer[offerName]; !ok {
					t.Errorf("%s - offer %s for ITI guarantee not found", data.Name, offerName)
					continue
				}
				guaranteesExpected = append(guaranteesExpected, g.Slug)
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.ITI.SumInsuredLimitOfIndemnity {
					t.Errorf("%s - mismatched offer %s - ITI sum. Expected %.2f - Got %.2f", data.Name, offerName, value.ITI.SumInsuredLimitOfIndemnity, g.Offer[offerName].SumInsuredLimitOfIndemnity)
				}
				if g.Offer[offerName].Deductible != value.ITI.Deductible {
					t.Errorf("%s - mismatched offer %s - ITI deductible. Expected %s - Got %s", data.Name, offerName, value.ITI.Deductible, g.Offer[offerName].Deductible)
				}
				if mismatchedPrice && g.Offer[offerName].PremiumGrossYearly != value.ITI.PriceGross {
					t.Errorf("%s - mismatched offer %s - ITI price. Expected %.2f - Got %.2f", data.Name, offerName, value.ITI.PriceGross, g.Offer[offerName].PremiumGrossYearly)
				}
			}

			for _, g := range p.Assets[0].Guarantees {
				for name := range g.Offer {
					if name == offerName {
						guaranteesGot = append(guaranteesGot, g.Slug)
					}
				}
			}

			if len(guaranteesExpected) != len(guaranteesGot) {
				t.Errorf("%s - mismatched offer %s - number of guarantees. Expected %+v - Got %+v", data.Name, offerName, guaranteesExpected, guaranteesGot)
			}
		}
		if numOffersExpected != numOffersGot {
			t.Errorf("%s - mismatched number of offers. Expected %d - Got %d", data.Name, numOffersExpected, numOffersGot)
			for name, o := range p.OffersPrices {
				t.Errorf("%s - Policy Offer %s: %+v", data.Name, name, o["yearly"])
			}
		}
	}
}

func buildPolicy(in InputData) models.Policy {
	assets := []models.Asset{{
		Person:     nil,
		Guarantees: make([]models.Guarante, 0),
	}}

	return models.Policy{
		ProductVersion: "v1",
		Name:           "persona",
		Company:        "global",
		Assets:         assets,
		QuoteQuestions: map[string]interface{}{
			"hasDependants": in.HasDependants,
			"lifeRisk":      in.LifeRisk,
			"financialRisk": in.FinancialRisk,
		},
		Contractor: models.Contractor{
			BirthDate: lib.AddMonths(time.Now().UTC().Truncate(time.Hour*24), in.Age*-12).Format(time.RFC3339),
			WorkType:  in.WorkType,
			RiskClass: in.RiskClass,
		},
	}
}
