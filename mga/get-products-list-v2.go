package mga

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

type GetProductsListByEntitlementResponse struct {
	Products []models.ProductInfo `json:"products"`
}

func GetProductsListByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetProductsListByChannelFx] Handler start ------------")

	var (
		err      error
		response GetProductsListByEntitlementResponse
	)

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("[GetProductsListByChannelFx] error extracting auth token: %s", err.Error())
		return "", "", err
	}

	// TODO use authToken.GetChannel so we only extract the node for networks
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)

	if authToken.Role == models.UserRoleAdmin {
		log.Println("[GetProductsListByChannelFx] getting mga products")
		response.Products = product.GetProductsByChannel(models.MgaChannel)
	} else if networkNode == nil {
		log.Println("[GetProductsListByChannelFx] getting e-commerce products")
		response.Products = product.GetProductsByChannel(models.ECommerceChannel)
	} else {
		log.Println("[GetProductsListByChannelFx] getting network products")
		warrant := network.GetWarrant(networkNode.Warrant)
		productList := lib.SliceMap[models.Product, string](warrant.Products, func(p models.Product) string { return p.Name })
		log.Printf("[GetProductsListByChannelFx] product list '%s'", productList)
		response.Products = product.GetNetworkNodeProducts(productList)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("[GetProductsListByChannelFx] error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Printf("[GetProductsListByChannelFx] found products: %s", string(responseBytes))

	return string(responseBytes), response, nil
}
