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

	alreadyEarthquake := policy.QuoteQuestions["alreadyEarthquake"].(bool)
	alreadyFlood := policy.QuoteQuestions["alreadyFlood"].(bool)
	useType := policy.Assets[0].Building.UseType
	//you must have both SumInsuredTextField(Fabricato) and SumInsuredLimitOfIndemnityTextField(Contenuto)
	//if i have alreadyEarthquake and alreadyflood and tenant, fabricato is mandatory
	isContenutoAndFabricato := func(value *models.GuaranteValue) bool {
		if value == nil {
			return false
		}
		if alreadyEarthquake && alreadyFlood && useType == "tenant" {
			val := value.SumInsuredLimitOfIndemnity
			if val == 0 {
				return false
			}
		}
		val := value.SumInsured
		if val == 0 {
			return false
		}
		return true
	}

	if g, err := policy.ExtractGuarantee("landslides"); err == nil {
		if !isContenutoAndFabricato(g.Value) && g.IsSelected {
			return nil, errors.New("You need atleast fabricato and contenuto for landSlide")
		}
	} else {
		return nil, errors.New("You need to select landslides")
	}

	if g, err := policy.ExtractGuarantee("earthquake"); err == nil {
		if !isContenutoAndFabricato(g.Value) && g.IsSelected {
			return nil, errors.New("You need atleast fabricato and contenuto for earthquake")
		}
	}

	if g, err := policy.ExtractGuarantee("flood"); err == nil {
		if !isContenutoAndFabricato(g.Value) && g.IsSelected {
			return nil, errors.New("You need atleast fabricato and contenuto for flood")
		}
	}
	out = ruleOutput.(*SellableOutput)
	return out, nil
}

func getCatnatInputRules(p *models.Policy) ([]byte, error) {
	var res []byte
	in := make(map[string]any)
	in["isEarthquakeSelected"] = false
	in["isFloodSelected"] = false
	locationlen := 0

	alreadyEarthquake := p.QuoteQuestions["alreadyEarthquake"]
	if alreadyEarthquake == nil {
		return nil, errors.New("missing field alreadyEarthquake")
	}
	alreadyFlood := p.QuoteQuestions["alreadyFlood"]
	if alreadyFlood == nil {
		return nil, errors.New("missing field alreadyFlood")
	}
	wantEarthquake := p.QuoteQuestions["wantEarthquake"]
	if wantEarthquake == nil {
		wantEarthquake = false
	}
	wantFlood := p.QuoteQuestions["wantFlood"]
	if wantFlood == nil {
		wantFlood = false
	}
	in["alreadyFlood"] = alreadyFlood
	in["wantFlood"] = wantFlood

	in["alreadyEarthquake"] = alreadyEarthquake
	in["wantEarthquake"] = wantEarthquake

	for _, a := range p.Assets {
		if a.Type == models.AssetTypeBuilding {
			locationlen += 1
		}
	}

	in["locationlen"] = locationlen
	res, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return res, nil
}
