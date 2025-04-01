package sellable

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

type SellableOutput struct {
	Msg     string
	Product *models.Product
}

func CatnatFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		policy *models.Policy
		err    error
	)
	log.SetPrefix("[CatnatFx] ")
	log.Println("Handler start -----------------------------------------------")

	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	policy.Normalize()

	var warrant *models.Warrant
	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)

	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if pr, err := Catnat(policy, policy.Channel, networkNode, warrant); err == nil {
		js, err := pr.Marshal()
		return string(js), err, nil
	}

	return "", nil, fmt.Errorf("policy not sellable by: %v", err)
}

func Catnat(p *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Product, error) {
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
	out = ruleOutput.(*SellableOutput)
	log.Println(out.Msg)
	return out.Product, nil
}

func getCatnatInputRules(p *models.Policy) ([]byte, error) {
	var res []byte
	out := make(map[string]any)
	locationlen := 0
	out["isEarthQuakeSelected"] = false
	out["isFloodSelected"] = false
	out["isLandSlidesSelected"] = false
	out["fabricato"] = 0
	out["contenuto"] = 0
	out["merci"] = 0

	setConfSetting := func(conf *models.GuaranteConfig) {
		if conf == nil {
			return
		}
		//both fabricato e contenuto == true or either fabricato or contenuto
		if val := conf.SumInsuredTextField; val != nil && len(val.Values) >= 1 && val.Values[0] > 0 {
			out["fabricato"] = conf.SumInsuredTextField.Values[0]
		}
		if val := conf.SumInsuredLimitOfIndemnityTextField; val != nil && len(val.Values) >= 1 && val.Values[0] > 0 {
			out["contenuto"] = conf.SumInsuredLimitOfIndemnityTextField.Values[0]
		}
		if val := conf.LimitOfIndemnityTextField; val != nil && len(val.Values) >= 1 && val.Values[0] > 0 {
			out["merci"] = conf.LimitOfIndemnityTextField.Values[0]
		}
	}
	if val, ok := p.QuoteQuestions["isEarthQuakeSelected"]; ok {
		out["isEarthQuakeSelected"] = val
	}
	if val, ok := p.QuoteQuestions["isFloodSelected"]; ok {
		out["isFloodSelected"] = val
	}

	if g, err := p.ExtractGuarantee("earthquake"); err == nil {
		setConfSetting(g.Config)
	}
	if g, err := p.ExtractGuarantee("flood"); err == nil {
		setConfSetting(g.Config)
	}
	if g, err := p.ExtractGuarantee("landslides"); err == nil {
		setConfSetting(g.Config)
		out["isLandSlidesSelected"] = g.IsSelected
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
