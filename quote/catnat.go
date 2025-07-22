package quote

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"
	"gitlab.dev.wopta.it/goworkspace/sellable"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type sellableCatnat = func(policy *models.Policy, product *models.Product, isValidationForQuote bool) (*sellable.SellableOutput, error)

type clientQuote = func(dto catnat.QuoteRequest, policy *models.Policy) (response catnat.QuoteResponse, err error)

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

	auth, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("error getting authToken")
		return "", nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.ErrorF("error decoding request body")
		return "", nil, err
	}
	client := catnat.NewNetClient()

	nodeUid := reqPolicy.PartnershipName
	if strings.EqualFold(reqPolicy.Channel, models.NetworkChannel) {
		nodeUid = auth.UserID
	}
	networkNode := network.GetNetworkNodeByUid(nodeUid)
	var warrant *models.Warrant
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	product := prd.GetProductV2(reqPolicy.Name, reqPolicy.ProductVersion, reqPolicy.Channel, networkNode, warrant)
	resp, err := catnatQuote(reqPolicy, product, sellable.CatnatSellable, client.Quote)
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

func catnatQuote(policy *models.Policy, product *models.Product, sellable sellableCatnat, clientQuote clientQuote) (resp catnat.QuoteResponse, err error) {
	outSellable, err := sellable(policy, product, true)
	if err != nil {
		return resp, err
	}
	internal.AddGuaranteesSettingsFromProduct(policy, outSellable.Product)

	var cnReq catnat.QuoteRequest
	err = cnReq.FromPolicyForQuote(policy)
	if err != nil {
		log.ErrorF("error building NetInsurance DTO: %s", err.Error())
		return resp, err
	}

	resp, err = clientQuote(cnReq, policy)
	log.PrintStruct("response quote", resp)
	if err != nil {
		return resp, err
	}
	internal.AddConsultacyPrice(policy, product)
	return
}
