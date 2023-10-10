package mga

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	pathPrefix = "products/"
)

type GetProductListResp struct {
	Products []ProductInfo `json:"products"`
}

type ProductInfo struct {
	Name         string `json:"name"`
	NameTitle    string `json:"nameTitle"`
	NameSubtitle string `json:"nameSubtitle"`
	NameDesc     string `json:"nameDesc"`
	Version      string `json:"version"`
	Company      string `json:"company"`
	Logo         string `json:"logo"`
}

func GetProductsListByEntitlementFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	var (
		err          error
		roleProducts = make([]models.Product, 0)
	)
	log.Println("GetProductsListByEntitlement")

	origin := r.Header.Get("Origin")

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	log.Println("GetProductsListByEntitlement: loading products list")
	switch authToken.Role {
	case models.UserRoleAdmin, models.UserRoleManager:
		roleProducts = getMgaProductsList()
	case models.UserRoleAgency, models.UserRoleAgent:
		roleProducts = getNetworkNodeProductsList(authToken.UserID, origin)
	case models.UserRoleCustomer, models.UserRoleAll:
		roleProducts = getEcommerceProductsList()
	}
	log.Printf("GetProductsListByEntitlement: found %d products for %s", len(roleProducts), authToken.Role)

	resp := GetProductListResp{Products: make([]ProductInfo, 0)}
	for _, product := range roleProducts {
		for _, company := range product.Companies {
			resp.Products = append(resp.Products, ProductInfo{
				Name:         product.Name,
				NameTitle:    product.NameTitle,
				NameSubtitle: product.NameSubtitle,
				NameDesc:     *product.NameDesc,
				Version:      product.Version,
				Company:      company.Name,
				Logo:         product.Logo,
			})
		}
	}

	jsonResp, err := json.Marshal(resp)

	return string(jsonResp), resp, err
}

func getMgaProductsList() []models.Product {
	productsList := make([]models.Product, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "mga/")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.Name == "pmi" {
			continue
		}
		productsList = append(productsList, product)
	}
	return productsList
}

func getEcommerceProductsList() []models.Product {
	productsList := make([]models.Product, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "e-commerce/")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.Name == "pmi" || !product.IsEcommerceActive {
			continue
		}
		productsList = append(productsList, product)
	}
	return productsList
}

func 
getNetworkNodeProductsList(networkNodeUid, origin string) []models.Product {
	var (
		networkNode  models.NetworkNode
		channel      string
		productsList = make([]models.Product, 0)
	)

	log.Println("GetAgencyProducts")

	fireNetwork := lib.GetDatasetByEnv(origin, models.NetworkNodesCollection)
	docsnap, err := lib.GetFirestoreErr(fireNetwork, networkNodeUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&networkNode)
	lib.CheckError(err)

	switch networkNode.Type {
	case models.AgentNetworkNodeType:
		channel = models.AgentChannel
	case models.AgencyNetworkNodeType:
		channel = models.AgencyChannel
	}

	defaultChannelProduct := getDefaultProductsByChannel(channel + "/")

	for _, product := range networkNode.Products {
		for _, defaultProduct := range defaultChannelProduct {
			if product.Name == defaultProduct.Name && product.Version == defaultProduct.Version {
				productsList = append(productsList, defaultProduct)
				break
			}
		}
	}

	return productsList
}


func getDefaultProductsByChannel(channel string) []models.Product {
	products := make([]models.Product, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + channel)
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		products = append(products, product)
	}
	return products
}
