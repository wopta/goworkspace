package sellable

import (
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	env "github.com/wopta/goworkspace/lib/environment"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

func assertEqual[t comparable](got, exp t, test *testing.T, namefield string) {
	if got != exp {
		log.ErrorF("%v Expected:%v, got:%v", namefield, exp, got)
		test.Fail()
	}
}
func fromGuaranteeMapToSlice(mapG map[string]*models.Guarante) (res []models.Guarante) {
	for _, m := range mapG {
		res = append(res, *m)
	}
	return
}
func setLandslideGuarantee(guarantees []models.Guarante) {
	for i, m := range guarantees {
		if m.Slug == "landslides" {
			guarantees[i].Value = &models.GuaranteValue{
				SumInsured:                 2,
				SumInsuredLimitOfIndemnity: 2,
			}
		}
	}
}
func getPrePopulatedPolicyForCatnat() models.Policy {
	return models.Policy{
		Name:           "cat-nat",
		ProductVersion: "v1",
		Channel:        "mga",
		StartDate:      time.Now(),
		EndDate:        time.Now(),
		Assets: slices.Clone([]models.Asset{
			{
				Type: models.AssetTypeBuilding,
				Guarantees: []models.Guarante{
					{Slug: "landslides"},
				},
				Building: &models.Building{},
			},
		}),
		QuoteQuestions: maps.Clone(map[string]any{}),
	}
}

func TestCatnatSellableNoAssets(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": true,
		"alreadyFlood":      true,
		"wantEarthquake":    false,
	}
	os.Setenv("env", env.LocalTest)
	policy.Assets = []models.Asset{}
	_, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	assertEqual(err.Error(), "Wrong numbers of locations", t, "")
}

func TestCatnatSellableWrongQuoteQuestionsAnswer(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": true,
		"alreadyFlood":      true,
		"wantEarthquake":    false,
	}
	_, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	assertEqual(err.Error(), "have to select at least earthquake or flood", t, "")
}

func TestCatnatSellableEarthQuakeConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": false,
		"alreadyFlood":      true,
		"wantFlood":         false,
	}
	output, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
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
}
func TestCatnatSellableFlood(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": true,
		"alreadyFlood":      true,
		"wantFlood":         true,
	}
	output, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
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
func TestCatnatSellableFloodWithQuoteAndNOConf(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()

	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": true,
		"wantEarthquake":    false,
		"alreadyFlood":      true,
		"wantFlood":         true,
	}
	out, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	if err != nil {
		t.Fatal(err)
	}
	policy.Assets[0].Guarantees = fromGuaranteeMapToSlice(out.Product.Companies[0].GuaranteesMap)
	setLandslideGuarantee(policy.Assets[0].Guarantees)
	_, err = CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), true)
	assertEqual(err.Error(), "You need atleast fabricato and contenuto for flood", t, "")
}
func TestCatnatSellableEartquakeWithQuoteOnlyFabricato(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyEarthquake": false,
		"alreadyFlood":      true,
		"wantFlood":         false,
	}
	out, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	if err != nil {
		t.Fatal(err)
	}
	policy.Assets[0].Guarantees = fromGuaranteeMapToSlice(out.Product.Companies[0].GuaranteesMap)
	setLandslideGuarantee(policy.Assets[0].Guarantees)
	_, err = CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), true)
	assertEqual(err.Error(), "You need atleast fabricato and contenuto for earthquake", t, "")
}

func TestCatnatSellableFloodWithQuote(t *testing.T) {
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
		"alreadyFlood":      false,
		"alreadyEarthquake": true,
		"wantEartquake":     false,
	}
	out, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	if err != nil {
		t.Fatal(err)
	}
	policy.Assets[0].Guarantees = fromGuaranteeMapToSlice(out.Product.Companies[0].GuaranteesMap)
	setLandslideGuarantee(policy.Assets[0].Guarantees)
	_, err = CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), true)
	assertEqual(err.Error(), "You need atleast fabricato and contenuto for flood", t, "")
}

func TestCatnatSellable(t *testing.T) {
	var policy = getPrePopulatedPolicyForCatnat()
	os.Setenv("env", env.LocalTest)
	policy.QuoteQuestions = map[string]any{
		"alreadyFlood":      false,
		"alreadyEarthquake": true,
		"wantEartquake":     false,
	}
	out, err := CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), false)
	if err != nil {
		t.Fatal(err)
	}
	policy.Assets[0].Guarantees = fromGuaranteeMapToSlice(out.Product.Companies[0].GuaranteesMap)
	policy.Assets[0].Guarantees[0].Value = &models.GuaranteValue{
		SumInsured:                 2,
		SumInsuredLimitOfIndemnity: 2,
	}
	policy.Assets[0].Guarantees[1].Value = &models.GuaranteValue{
		SumInsured:                 2,
		SumInsuredLimitOfIndemnity: 2,
	}
	setLandslideGuarantee(policy.Assets[0].Guarantees)
	_, err = CatnatSellable(&policy, product.GetProductV2(policy.Name, "v1", "mga", nil, nil), true)
	assertEqual(err, nil, t, "")
}
