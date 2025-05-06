package sellable

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
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
	)
	log.AddPrefix("[CatnatFx]")
	defer log.PopPrefix()
	log.Println("Handler start -----------------------------------------------")

	step := chi.URLParam(r, "step")

	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
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

	if pr, err := catnatSellable(policy, policy.Channel, networkNode, warrant, step); err == nil {
		js, err := pr.Product.Marshal()
		return string(js), err, nil
	}
	return "", nil, fmt.Errorf("policy not sellable by: %v", err)
}

func catnatSellable(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, step string) (*SellableOutput, error) {
	log.AddPrefix("CatnatSellalble")
	defer log.PopPrefix()

	in, err := getCatnatInputRules(policy)
	if err != nil {
		return nil, err
	}

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)
	fx := new(models.Fx)
	product := product.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		return nil, errors.New("Error getting catnat product")
	}
	out := &SellableOutput{
		Msg:     "",
		Product: product,
	}
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, out, in, nil)

	if step != quoteStep {
		out = ruleOutput.(*SellableOutput)
		log.InfoF(out.Msg)
		return out, nil
	}

	//you must have both SumInsuredTextField(Fabricato) and SumInsuredLimitOfIndemnityTextField(Contenuto)
	isContenutoAndFabricato := func(conf *models.GuaranteConfig) bool {
		val := conf.SumInsuredTextField
		if val == nil || len(val.Values) == 0 || val.Values[0] == 0 {
			return false
		}
		val = conf.SumInsuredLimitOfIndemnityTextField
		if val == nil || len(val.Values) == 0 || val.Values[0] == 0 {
			return false
		}
		return true
	}
	if g, err := policy.ExtractGuarantee("landslides"); err == nil {
		if product.Companies[0].GuaranteesMap["landslides"].IsSelected && (g.Config == nil || !isContenutoAndFabricato(g.Config)) {
			out.Msg += "You need atleast fabricato and contenuto for landslide"
		}
	} else {
		out.Msg += ("you must have landslides")
		return out, err
	}

	if g, err := policy.ExtractGuarantee("earthQuake"); err == nil {
		if product.Companies[0].GuaranteesMap["earthQuake"].IsSelected && (g.Config == nil || !isContenutoAndFabricato(g.Config)) {
			out.Msg += "You need atleast fabricato and contenuto for earthquake"
		}
	}
	if g, err := policy.ExtractGuarantee("flood"); err == nil {
		if product.Companies[0].GuaranteesMap["flood"].IsSelected && (g.Config == nil || !isContenutoAndFabricato(g.Config)) {
			out.Msg += "You need atleast fabricato and contenuto for flood"
		}
	}
	out = ruleOutput.(*SellableOutput)
	return out, nil
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
