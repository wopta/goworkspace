package sellable

import (
	"encoding/json"
	"errors"
	"net/http"

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
	log.AddPrefix("CatnatFx")
	defer log.PopPrefix()
	log.Println("Handler start -----------------------------------------------")

	defer func() {
		r.Body.Close()
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end ---------------------------------------------")
	}()

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

	pr, err := CatnatSellable(policy, policy.Channel, networkNode, warrant, false)
	if err != nil {
		return "", nil, err
	}
	js, err := pr.Product.Marshal()
	return string(js), nil, nil
}

func CatnatSellable(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, isValidationForQuote bool) (*SellableOutput, error) {
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

	if len(out.Msg) != 0 {
		out = ruleOutput.(*SellableOutput)
		return nil, errors.New(out.Msg)
	}

	if !isValidationForQuote {
		out = ruleOutput.(*SellableOutput)
		log.InfoF(out.Msg)
		return out, nil
	}

	//you must have both SumInsuredTextField(Fabricato) and SumInsuredLimitOfIndemnityTextField(Contenuto)
	isContenutoAndFabricato := func(value *models.GuaranteValue) bool {
		val := value.SumInsured
		if val == 0 {
			return false
		}
		val = value.SumInsuredLimitOfIndemnity
		if val == 0 {
			return false
		}
		return true
	}

	if g, err := policy.ExtractGuarantee("landslide"); err == nil {
		if g.Value == nil || !isContenutoAndFabricato(g.Value) {
			out.Msg += "You need atleast fabricato and contenuto for landSlide,"
		}
	}

	if g, err := policy.ExtractGuarantee("earthquake"); err == nil {
		if g.Value == nil || !isContenutoAndFabricato(g.Value) {
			out.Msg += "You need atleast fabricato and contenuto for earthquake,"
		}
	}

	if g, err := policy.ExtractGuarantee("flood"); err == nil {
		if g.Value == nil || !isContenutoAndFabricato(g.Value) {
			out.Msg += "You need atleast fabricato and contenuto for flood,"
		}
	}
	out = ruleOutput.(*SellableOutput)
	return out, nil
}

func getCatnatInputRules(p *models.Policy) ([]byte, error) {
	var res []byte
	out := make(map[string]any)
	out["isEarthquakeSelected"] = false
	out["isFloodSelected"] = false
	locationlen := 0

	alreadyEarthquake := p.QuoteQuestions["alreadyEarthquake"]
	if alreadyEarthquake == nil {
		alreadyEarthquake = false
	}
	alreadyFlood := p.QuoteQuestions["alreadyFlood"]
	if alreadyFlood == nil {
		alreadyFlood = false
	}
	wantEarthquake := p.QuoteQuestions["wantEarthquake"]
	if wantEarthquake == nil {
		wantEarthquake = false
	}
	wantFlood := p.QuoteQuestions["wantFlood"]
	if wantFlood == nil {
		wantFlood = false
	}
	out["isEarthquakeSelected"] = (alreadyEarthquake.(bool) && wantEarthquake.(bool)) || !alreadyEarthquake.(bool)
	out["isFloodSelected"] = ((alreadyFlood).(bool) && wantFlood.(bool)) || !alreadyFlood.(bool)

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
