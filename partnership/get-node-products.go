package partnership

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
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

	log.SetPrefix("[GetPartnershipNodeAndProductsFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	partnershipUid := chi.URLParam(r, "partnershipUid")
	jwtData := r.URL.Query().Get("jwt")

	if node, err = network.GetNodeByUid(partnershipUid); err != nil {
		log.Printf("error getting node '%s': %s", partnershipUid, err.Error())
		return "", nil, err
	}

	if _, err := node.DecryptJwt(jwtData); err != nil {
		log.Printf("error decoding jwt: %s", err.Error())
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
