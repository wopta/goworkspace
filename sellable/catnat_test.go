package sellable

import (
	"os"
	"testing"

	env "github.com/wopta/goworkspace/lib/environment"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func assertEqual[t comparable](a, exp t, test *testing.T, namefield string) {
	if a != exp {
		log.ErrorF("%v Expected:%v, got:%v", namefield, exp, a)
		test.Fail()
	}
}
func getPrePopulatedPolicyForCatnat() models.Policy {
	return models.Policy{
		Name:           "cat-nat",
		ProductVersion: "v1",
		Channel:        "mga",
		Assets: []models.Asset{
			{
				Type: models.AssetTypeBuilding,
				Guarantees: []models.Guarante{
					{Slug: "earthquake"},
					{Slug: "flood"},
					{Slug: "landslides",
						Config: &models.GuaranteConfig{},
					},
				}},
		},
		QuoteQuestions: map[string]any{
			"isEarthQuakeSelected": false,
			"isFloodSelected":      true,
		},
	}
}

func TestCatnatSellableNoAssets(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.Assets = []models.Asset{}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "Wrong numbers of locations", t, "")
}

func TestCatnatSellableWrongQuoteQuestionsAnswer(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      false,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "have to select at least earthQuake or flood", t, "")
}

func TestCatnatSellableEarthQuakeConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": true,
		"isFloodSelected":      false,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsSelected, true, t, "isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsSellable, true, t, "isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsMandatory, true, t, "isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsConfigurable, false, t, "isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSelected, false, t, "isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSellable, false, t, "isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsMandatory, false, t, "isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsConfigurable, false, t, "isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsSelected, false, t, "isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsSellable, true, t, "isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsMandatory, false, t, "isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsConfigurable, true, t, "isConfigurable")
}
func TestCatnatSellableFlood(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSelected, true, t, "isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSellable, true, t, "isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsMandatory, true, t, "isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsConfigurable, false, t, "isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsSelected, false, t, "isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsSellable, false, t, "isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsMandatory, false, t, "isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthQuake"].IsConfigurable, false, t, "isConfigurable")
}
func TestCatnatSellableLandSlideWithQuoteAndNOConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for flood", t, "")
}
func TestCatnatSellableLandSlideWithQuoteOnlyFabricato(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees[2].Config = &models.GuaranteConfig{ //3 is landslides, hard coded for test
		SumInsuredTextField: &models.GuaranteFieldConfig{Values: []float64{2}},
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for flood", t, "")
}

func TestCatnatSellableLandSlideWithQuote(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees[2].Config = &models.GuaranteConfig{ //3 is landslides, hard coded for test
		SumInsuredTextField:                 &models.GuaranteFieldConfig{Values: []float64{2}},
		SumInsuredLimitOfIndemnityTextField: &models.GuaranteFieldConfig{Values: []float64{2}},
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for flood", t, "")
}

func TestCatnatSellable(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees[2].Config = &models.GuaranteConfig{ //2 is landslides, hard coded for test
		SumInsuredTextField:                 &models.GuaranteFieldConfig{Values: []float64{2}},
		SumInsuredLimitOfIndemnityTextField: &models.GuaranteFieldConfig{Values: []float64{2}},
	}
	policy.Assets[0].Guarantees[1].Config = &models.GuaranteConfig{ //1 is flood, hard coded for test
		SumInsuredTextField:                 &models.GuaranteFieldConfig{Values: []float64{2}},
		SumInsuredLimitOfIndemnityTextField: &models.GuaranteFieldConfig{Values: []float64{2}},
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthQuakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
}
