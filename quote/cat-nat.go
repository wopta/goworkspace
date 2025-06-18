package quote

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/sellable"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
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
	resp, err := CatnatQuote(reqPolicy, product, sellable.CatnatSellable, client.Quote)
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
