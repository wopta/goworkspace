package quote

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

func mock_sellable(policy *models.Policy, product *models.Product, _ bool) (*sellable.SellableOutput, error) {
	return sellable.CatnatSellable(policy, product, true)
}

type mock_clientCatnat struct {
	withError         bool
	nameFileToCompare string
}

func (c *mock_clientCatnat) Quote(dto catnat.QuoteRequest, _ *models.Policy) (response catnat.QuoteResponse, err error) {
	if c.withError {
		return response, errors.New("quote error")
	}
	var bytes []byte
	var dtoExpected catnat.QuoteRequest
	bytes, e := lib.GetFilesByEnvV2("data/test/quote/catnat/" + c.nameFileToCompare)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(bytes, &dtoExpected)
	if !reflect.DeepEqual(dto, dtoExpected) {
		log.PrintStruct("Expected: ", dtoExpected)
		log.PrintStruct("Got: ", dto)
		return response, errors.New("dto != expected")
	}
	return response, nil
}

func getPolicyWithEverythingForTest() *models.Policy {
	var policy models.Policy
	bytes, e := lib.GetFilesByEnvV2("data/test/quote/catnat/input_policy.json")
	if e != nil {
		panic(e)
	}
	if len(bytes) == 0 {
		panic("error retrieving policy")
	}
	json.Unmarshal(bytes, &policy)
	return &policy
}

func TestQuoteCatnatWithEverything(t *testing.T) {
	t.Setenv("env", env.LocalTest)
	policy := getPolicyWithEverythingForTest()
	policy.QuoteQuestions["alreadyEarthquake"] = false
	policy.QuoteQuestions["alreadyFlood"] = false
	product := product.GetProductV2(policy.Name, "v1", "mga", nil, nil)
	if product == nil {
		t.Fatal("error retrieving product")
	}
	client := new(mock_clientCatnat)
	client.nameFileToCompare = "output_everything_alreadyfalse.json"
	_, err := catnatQuote(policy, product, mock_sellable, client.Quote)
	if err != nil {
		t.Fatal(err)
	}
}
func TestQuoteCatnatWithEverythingButEarthquake(t *testing.T) {
	t.Setenv("env", env.LocalTest)
	policy := getPolicyWithEverythingForTest()
	policy.QuoteQuestions["alreadyEarthquake"] = true
	policy.QuoteQuestions["wantEarthquake"] = false
	policy.QuoteQuestions["alreadyFlood"] = true
	policy.QuoteQuestions["wantFlood"] = true
	product := product.GetProductV2(policy.Name, "v1", "mga", nil, nil)
	if product == nil {
		t.Fatal("error retrieving product")
	}
	client := new(mock_clientCatnat)
	client.nameFileToCompare = "output_noearthquake.json"
	_, err := catnatQuote(policy, product, mock_sellable, client.Quote)
	if err != nil {
		t.Fatal(err)
	}
}
func TestQuoteCatnatWithEverything2(t *testing.T) {
	t.Setenv("env", env.LocalTest)
	policy := getPolicyWithEverythingForTest()
	policy.QuoteQuestions["alreadyEarthquake"] = true
	policy.QuoteQuestions["wantEarthquake"] = true
	policy.QuoteQuestions["alreadyFlood"] = true
	policy.QuoteQuestions["wantFlood"] = true
	product := product.GetProductV2(policy.Name, "v1", "mga", nil, nil)
	if product == nil {
		t.Fatal("error retrieving product")
	}
	client := new(mock_clientCatnat)
	client.nameFileToCompare = "output_everything_alreadytrue.json"

	_, err := catnatQuote(policy, product, mock_sellable, client.Quote)
	if err != nil {
		t.Fatal(err)
	}
}
