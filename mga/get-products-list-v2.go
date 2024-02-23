package mga

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

type GetProductsListByEntitlementResponse struct {
	Products []models.ProductInfo `json:"products"`
}

func GetProductsListByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		response GetProductsListByEntitlementResponse
	)

	log.SetPrefix("[GetProductsListByChannelFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error extracting auth token: %s", err.Error())
		return "", "", err
	}

	channel := authToken.GetChannelByRoleV2()
	switch channel {
	case models.MgaChannel:
		log.Println("getting mga products")
		response.Products = product.GetProductsByChannel(models.MgaChannel)
	case models.ECommerceChannel:
		log.Println("getting e-commerce products")
		response.Products = product.GetProductsByChannel(models.ECommerceChannel)
	case models.NetworkChannel:
		log.Println("getting network products")
		networkNode := network.GetNetworkNodeByUid(authToken.UserID)
		if networkNode == nil {
			log.Println("node not found")
			return "", "", fmt.Errorf("no node set for authToken")
		}
		warrant := networkNode.GetWarrant()
		if warrant == nil {
			log.Println("warrant not found")
			return "", "", fmt.Errorf("no warrant set for node")
		}
		productList := lib.SliceMap[models.Product, string](warrant.Products, func(p models.Product) string { return p.Name })
		log.Printf("product list '%s'", productList)
		response.Products = product.GetNetworkNodeProducts(productList)
	default:
		log.Printf("error channel %s unaavailable", channel)
		return "", "", fmt.Errorf("unavailable channel")
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Printf("found products: %s", string(responseBytes))

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}
