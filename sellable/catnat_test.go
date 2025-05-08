package sellable

import (
	"os"
	"testing"

	env "github.com/wopta/goworkspace/lib/environment"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func assertEqual[t comparable](got, exp t, test *testing.T, namefield string) {
	if got != exp {
		log.ErrorF("%v Expected:%v, got:%v", namefield, exp, got)
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
					{Slug: "landslides"},
				}},
		},
		QuoteQuestions: map[string]any{
			"isEarthquakeSelected": false,
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
		"isEarthquakeSelected": false,
		"isFloodSelected":      false,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "have to select at least earthquake or flood", t, "")
}

func TestCatnatSellableEarthQuakeConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": true,
		"isFloodSelected":      false,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsSelected, true, t, "earthquake_isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsSellable, true, t, "earthquake_isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsMandatory, true, t, "earthquake_isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsConfigurable, false, t, "earthquake_isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSelected, false, t, "flood_isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSellable, false, t, "flood_isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsMandatory, false, t, "flood_isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsConfigurable, false, t, "flood_isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsSelected, false, t, "landslides_isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsSellable, true, t, "landslides_isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsMandatory, false, t, "landslides_isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["landslides"].IsConfigurable, true, t, "landslides_isConfigurable")
}
func TestCatnatSellableFlood(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSelected, true, t, "flood_isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsSellable, true, t, "flood_isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsMandatory, true, t, "flood_isMandatory")
	assertEqual(output.Product.Companies[0].GuaranteesMap["flood"].IsConfigurable, false, t, "flood_isConfigurable")

	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsSelected, false, t, "earthquake_isSelected")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsSellable, false, t, "earthquake_isSellable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsConfigurable, false, t, "earthquake_isConfigurable")
	assertEqual(output.Product.Companies[0].GuaranteesMap["earthquake"].IsMandatory, false, t, "earthquake_isMandatory")
}
func TestCatnatSellableLandSlideWithQuoteAndNOConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for landslides,", t, "")
}
func TestCatnatSellableLandSlideWithQuoteOnlyFabricato(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees = append(policy.Assets[0].Guarantees,
		models.Guarante{Slug: "flood"},
	)
	policy.Assets[0].Guarantees[1].Value = &models.GuaranteValue{
		SumInsured: 2,
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for landslides,You need atleast fabricato and contenuto for flood,", t, "")
}

func TestCatnatSellableLandSlideWithQuote(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees = append(policy.Assets[0].Guarantees,
		models.Guarante{Slug: "flood"},
	)
	policy.Assets[0].Guarantees[0].Value = &models.GuaranteValue{
		SumInsured:                 2,
		SumInsuredLimitOfIndemnity: 2,
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "You need atleast fabricato and contenuto for flood,", t, "")
}

func TestCatnatSellable(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	policy.Assets[0].Guarantees = append(policy.Assets[0].Guarantees,
		models.Guarante{Slug: "flood"},
	)
	policy.Assets[0].Guarantees[1].Value = &models.GuaranteValue{
		SumInsured:                 2,
		SumInsuredLimitOfIndemnity: 2,
	}
	policy.Assets[0].Guarantees[0].Value = &models.GuaranteValue{
		SumInsured:                 2,
		SumInsuredLimitOfIndemnity: 2,
	}
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"isEarthquakeSelected": false,
		"isFloodSelected":      true,
	}
	output, err := CatnatSellable(&policy, policy.Channel, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(output.Msg, "", t, "")
}
