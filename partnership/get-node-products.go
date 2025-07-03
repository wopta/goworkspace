package partnership

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/product"
)

type GetPartnershipNodeAndProductsResp struct {
	Partnership PartnershipNode      `json:"partnership"`
	Products    []models.ProductInfo `json:"products"`
}

func GetPartnershipNodeAndProductsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		response GetPartnershipNodeAndProductsResp
		node     *models.NetworkNode
		err      error
	)

	log.AddPrefix("GetPartnershipNodeAndProductsFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	partnershipUid := chi.URLParam(r, "partnershipUid")
	jwtData := r.URL.Query().Get("jwt")

	if node, err = network.GetNodeByUidErr(partnershipUid); err != nil {
		log.ErrorF("error getting node '%s': %s", partnershipUid, err.Error())
		return "", nil, err
	}

	if node == nil {
		response.Partnership.Name = partnershipUid
		response.Products = product.GetAllProductsByChannel(lib.ECommerceChannel)
		responseJson, err := json.Marshal(response)
		return string(responseJson), nil, err
	}

	if _, err := node.DecryptJwt(jwtData); err != nil {
		log.ErrorF("error decoding jwt: %s", err.Error())
		return "", nil, err
	}

	productList := lib.SliceMap(node.Products, func(p models.Product) string { return p.Name })
	productInfos := product.GetProductsByChannel(productList, lib.ECommerceChannel)

	response.Partnership = PartnershipNode{node.Partnership.Name, node.Partnership.Skin}
	response.Products = productInfos

	responseJson, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}
