package sellable

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"log"
	"net/http"
)

const quoteStep = "quote"

type SellableOutput struct {
	Msg     string
	Product *models.Product
}

func CatnatFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		policy *models.Policy
		err    error
		pr     *models.Product
	)
	log.SetPrefix("[CatnatFx]")
	log.Println("Handler start -----------------------------------------------")
	log.Prefix()

	step := chi.URLParam(r, "step")

	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	log.Println("error decoding request body")
	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		return "", nil, err
	}

	policy.Normalize()

	var warrant *models.Warrant
	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)

	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if pr, err = Catnat(policy, policy.Channel, networkNode, warrant, step); err == nil {
		js, err := pr.Marshal()
		return string(js), err, nil
	}

	return "", nil, fmt.Errorf("policy not sellable by: %v", err)
}

func Catnat(p *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, step string) (*models.Product, error) {
	log.SetPrefix("[CatnatSellalble]")
	defer log.SetPrefix("")

	in, err := getCatnatInputRules(p)
	if err != nil {
		return nil, err
	}

	rulesFile := lib.GetRulesFileV2(p.Name, p.ProductVersion, rulesFilename)
	fx := new(models.Fx)
	product := product.GetProductV2(p.Name, p.ProductVersion, channel, networkNode, warrant)
	out := &SellableOutput{
		Msg:     "",
		Product: product,
	}

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, out, in, nil)

	isGuarantConfigValid := func(conf *models.GuaranteConfig) bool {
		var (
			fabricato float64
			contenuto float64
		)

		//both fabricato e contenuto == true or either fabricato or contenuto
		if val := conf.SumInsuredTextField; val != nil && len(val.Values) >= 1 && val.Values[0] > 0 {
			fabricato = conf.SumInsuredTextField.Values[0]
		}
		if val := conf.SumInsuredLimitOfIndemnityTextField; val != nil && len(val.Values) >= 1 && val.Values[0] > 0 {
			contenuto = conf.SumInsuredLimitOfIndemnityTextField.Values[0]
		}
		return fabricato+contenuto != 0
	}
	if step != quoteStep {
		out = ruleOutput.(*SellableOutput)
		return out.Product, nil
	}
	var isValidResult bool
	if g, err := p.ExtractGuarantee("earthquake"); err == nil {
		isValidResult = isValidResult && isGuarantConfigValid(g.Config)
	}
	if g, err := p.ExtractGuarantee("flood"); err == nil {
		isValidResult = isValidResult && isGuarantConfigValid(g.Config)
	}
	if g, err := p.ExtractGuarantee("landlides"); err == nil {
		isValidResult = isValidResult && isGuarantConfigValid(g.Config)
	}
	if !isValidResult {
		return nil, fmt.Errorf("you need atleast fabricato or contenuto")
	}
	out = ruleOutput.(*SellableOutput)
	return out.Product, nil
}

func getCatnatInputRules(p *models.Policy) ([]byte, error) {
	var res []byte
	out := make(map[string]any)
	locationlen := 0
	out["isEarthQuakeSelected"] = false
	out["isFloodSelected"] = false

	if val, ok := p.QuoteQuestions["isEarthQuakeSelected"]; ok {
		out["isEarthQuakeSelected"] = val
	}
	if val, ok := p.QuoteQuestions["isFloodSelected"]; ok {
		out["isFloodSelected"] = val
	}

	for _, a := range p.Assets {
		if a.Type == models.AssetTypeBuilding {
			locationlen += 1
		}
	}

	out["locationlen"] = locationlen
	res, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return res, nil
}
