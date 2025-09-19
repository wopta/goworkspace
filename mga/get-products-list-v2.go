package mga

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/product"
)

type GetProductsListByEntitlementResponse struct {
	Products []models.ProductInfo `json:"products"`
}

func getProductsListByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		response GetProductsListByEntitlementResponse
	)

	log.AddPrefix("GetProductsListByChannelFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.ErrorF("error extracting auth token: %s", err.Error())
		return "", "", err
	}

	channel := authToken.GetChannelByRoleV2()
	switch channel {
	case models.MgaChannel:
		log.Println("getting mga products")
		response.Products = product.GetAllProductsByChannel(models.MgaChannel)
	case models.ECommerceChannel:
		log.Println("getting e-commerce products")
		response.Products = product.GetAllProductsByChannel(models.ECommerceChannel)
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
		retrievedProducts := product.GetProductsByChannel(productList, channel)
		for index, prd := range retrievedProducts {
			for _, warrantProduct := range warrant.Products {
				if prd.Name == warrantProduct.Name {
					retrievedProducts[index].IsActive = warrantProduct.IsActive
				}
			}
		}
		response.Products = retrievedProducts
	default:
		log.ErrorF("error channel %s unaavailable", channel)
		return "", "", fmt.Errorf("unavailable channel")
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "", "", err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(responseBytes), response, nil
}
