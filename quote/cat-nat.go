package quote

import (
	"encoding/json"
	"net/http"

	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/quote/catnat"
	"github.com/wopta/goworkspace/sellable"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CatNatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		reqPolicy *models.Policy
	)

	log.AddPrefix("CatNatFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	_, err = lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("error getting authToken")
		return "", nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.ErrorF("error decoding request body")
		return "", nil, err
	}
	client := catnat.NewNetClient()

	networkNode := network.GetNetworkNodeByUid(reqPolicy.ProducerUid)
	var warrant *models.Warrant
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product := prd.GetProductV2(reqPolicy.Name, reqPolicy.ProductVersion, reqPolicy.Channel, networkNode, warrant)
	resp, err := catnat.CatnatQuote(reqPolicy, product, sellable.CatnatSellable, client)
	if err != nil {
		return "", nil, err
	}
	cnReqStr, err := json.Marshal(resp)
	if err != nil {
		return "", nil, err
	}
	log.InfoF(string(cnReqStr))
	var out []byte
	out, err = json.Marshal(reqPolicy)
	if err != nil {
		log.ErrorF("error encoding response %v", err.Error())
		return "", nil, err
	}

	return string(out), out, err
}
