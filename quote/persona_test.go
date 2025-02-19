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
	Work          string `json:"work"`
	WorkType      string `json:"workType"`
	RiskClass     string `json:"riskClass"`
}

type OutputGuaranteeData struct {
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
	Deductible                 string  `json:"deductible"`
	DeductibleType             string  `json:"deductibleType"`
}

type OutputOfferData struct {
	IPI        *OutputGuaranteeData `json:"IPI,omitempty"`
	D          *OutputGuaranteeData `json:"D,omitempty"`
	DRG        *OutputGuaranteeData `json:"DRG,omitempty"`
	DC         *OutputGuaranteeData `json:"DC,omitempty"`
	RSC        *OutputGuaranteeData `json:"RSC,omitempty"`
	PriceGross float64              `json:"priceGross"`
}

type OutputData = map[string]OutputOfferData

type TestData struct {
	Input  InputData  `json:"input"`
	Output OutputData `json:"output"`
}

func TestPersona(t *testing.T) {
	os.Setenv("env", "local-test")

	fileReader, err := os.Open("../../function-data/dev/" + "data/test/quote/persona.json")
	if err != nil {
		t.Fatalf("unable to load data: %s", err)
	}

	testData := make([]TestData, 0)

	if err := json.NewDecoder(fileReader).Decode(&testData); err != nil {
		t.Fatalf("unable to decode data: %s", err)
	}

	for idx, data := range testData {
		p := buildPolicy(data.Input)
		if err := quote.Persona(&p, models.ECommerceChannel, nil, nil, models.ECommerceFlow); err != nil {
			t.Fatalf("error quoting test %d: %s", idx+1, err)
		}
		for offerName, value := range data.Output {
			offerPrice := p.OffersPrices[offerName]["yearly"].Gross
			if offerPrice != value.PriceGross {
				t.Fatalf("mismatched offer price for %s. Expected %.2f - Got %.2f", offerName, value.PriceGross, offerPrice)
			}
			if value.IPI != nil {
				g, _ := p.ExtractGuarantee("IPI")
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.IPI.SumInsuredLimitOfIndemnity {
					t.Fatalf("mismatched offer sum for %s. Expected %.2f - Got %.2f", offerName, g.Offer[offerName].SumInsuredLimitOfIndemnity, value.IPI.SumInsuredLimitOfIndemnity)
				}
			}
			if value.D != nil {
				g, _ := p.ExtractGuarantee("D")
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.D.SumInsuredLimitOfIndemnity {
					t.Fatalf("mismatched offer sum for %s. Expected %.2f - Got %.2f", offerName, g.Offer[offerName].SumInsuredLimitOfIndemnity, value.D.SumInsuredLimitOfIndemnity)
				}
			}
			if value.DRG != nil {
				g, _ := p.ExtractGuarantee("DRG")
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.DRG.SumInsuredLimitOfIndemnity {
					t.Fatalf("mismatched offer sum for %s. Expected %.2f - Got %.2f", offerName, g.Offer[offerName].SumInsuredLimitOfIndemnity, value.DRG.SumInsuredLimitOfIndemnity)
				}
			}
			if value.DC != nil {
				g, _ := p.ExtractGuarantee("DC")
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.DC.SumInsuredLimitOfIndemnity {
					t.Fatalf("mismatched offer sum for %s. Expected %.2f - Got %.2f", offerName, g.Offer[offerName].SumInsuredLimitOfIndemnity, value.DC.SumInsuredLimitOfIndemnity)
				}
			}
			if value.RSC != nil {
				g, _ := p.ExtractGuarantee("RSC")
				if g.Offer[offerName].SumInsuredLimitOfIndemnity != value.RSC.SumInsuredLimitOfIndemnity {
					t.Fatalf("mismatched offer sum for %s. Expected %.2f - Got %.2f", offerName, g.Offer[offerName].SumInsuredLimitOfIndemnity, value.RSC.SumInsuredLimitOfIndemnity)
				}
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
			BirthDate: lib.AddMonths(time.Now(), in.Age*-12).Format(time.RFC3339),
			Work:      in.Work,
			WorkType:  in.WorkType,
			RiskClass: in.RiskClass,
		},
	}
}
