package sellable

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/product"
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

	product := product.GetLatestActiveProduct(policy.Name, policy.Channel, networkNode, warrant)
	if product == nil {
		return "", nil, errors.New("Error getting catnat product")
	}

	pr, err := CatnatSellable(policy, product, false)
	if err != nil {
		return "", nil, err
	}
	js, err := pr.Product.Marshal()
	return string(js), nil, nil
}

type fxForCatnat struct {
	*models.Fx
}

func (fx *fxForCatnat) RemoveGuaranteeGroup(product *models.Product, groupKey string) {
	for key, value := range product.Companies[0].GuaranteesMap {
		if value.Group == groupKey {
			fx.RemoveGuarantee(product.Companies[0].GuaranteesMap, key)
		}
	}
}

func (fx *fxForCatnat) SetAsSelected(product *models.Product, groupKey string) {
	for _, value := range product.Companies[0].GuaranteesMap {
		if value.Group == groupKey {
			value.IsMandatory = true
			value.IsSelected = true
			value.IsSellable = true
			value.IsConfigurable = false
		}
	}
}

func (fx *fxForCatnat) SetAsNoSelected(product *models.Product, groupKey string) {
	for _, value := range product.Companies[0].GuaranteesMap {
		if value.Group == groupKey {
			value.IsMandatory = false
			value.IsSelected = false
			value.IsSellable = false
			value.IsConfigurable = false
		}
	}
}

func CatnatSellable(policy *models.Policy, product *models.Product, isValidationForQuote bool) (*SellableOutput, error) {
	log.AddPrefix("CatnatSellalble")
	defer log.PopPrefix()

	in, err := getCatnatInputRules(policy)
	if err != nil {
		return nil, err
	}

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)
	fx := new(fxForCatnat)
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
	//you must have both 'building' and 'content'
	//if i have alreadyEarthquake and alreadyflood and tenant, fabricato is mandatory
	isContenutoAndFabricato := func(types []string) error {
		isContent := slices.Contains(types, "content")
		isBuilding := slices.Contains(types, "building")

		if !isContent {
			return errors.New("Contenuto é obbligatorio")
		}
		if alreadyEarthquake && alreadyFlood && useType == "tenant" {
			for _, guarantee := range product.Companies[0].GuaranteesMap {
				if strings.HasSuffix(guarantee.SynchronizeSlug, "building") {
					guarantee.Config.SumInsuredLimitOfIndemnityTextField.Min = 0
				}
			}
			return nil
		}

		if !isBuilding {
			return errors.New("Fabbricato é obbligatorio")
		}
		return nil
	}
	if policy.StartDate.IsZero() {
		return nil, errors.New("Start date can't be 0")
	}
	if policy.EndDate.IsZero() {
		return nil, errors.New("End date can't be 0")
	}

	guaranteeExist := func(policy *models.Policy, groupName string) (isSelected bool, error error) {
		var types []string
		var companyName string
		for _, guarantee := range policy.Assets[0].Guarantees {

			companyName = guarantee.CompanyName
			if guarantee.Group == groupName {
				isSelected = guarantee.IsSelected
				_, typeName, _ := strings.Cut(guarantee.Slug, "-")
				if guarantee.Value.SumInsuredLimitOfIndemnity > 0 {
					types = append(types, typeName)
				}
			}
		}
		err := isContenutoAndFabricato(types)
		if err != nil {
			return isSelected, fmt.Errorf("Per %s %w", companyName, err)
		}
		return isSelected, nil
	}

	exist, err := guaranteeExist(policy, "LANDSLIDE")
	if !exist {
		return nil, errors.New("You need to have landslide")
	}
	if err != nil {
		return nil, err
	}

	_, err = guaranteeExist(policy, "EARTHQUAKE")
	if err != nil {
		return nil, err
	}
	_, err = guaranteeExist(policy, "FLOOD")
	if err != nil {
		return nil, err
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
