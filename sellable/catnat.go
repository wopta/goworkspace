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

type SellableOutput struct {
	Msg     string
	Product *models.Product
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

	if out.Msg == "" {
		return out.Product, nil
	}

	return out.Product, fmt.Errorf(out.Msg)
}

func getCatnatInputRules(p *models.Policy) ([]byte, error) {
	var res []byte
	out := make(map[string]any)
	locationlen := 0
	out["isEarthQuakeSelected"]=false
	out["isFloodSelected"]=false
	out["isSlideLandsSelected"]=false

	if g,err:=p.ExtractGuarantee("earthquake");err==nil{
		log.Println( "jfdsklfjds")
		out["isEarthQuakeSelected"]=g.IsSelected
	}else{
		log.Println(err)
	}
	if g,err:=p.ExtractGuarantee("flood");err==nil{
		out["isFloodSelected"]=g.IsSelected
	}
	if g,err:=p.ExtractGuarantee("slidelands");err==nil{
		out["isSlideLandsSelected"]=g.IsSelected
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
