package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
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
	case models.UserRoleAgency:
		roleProducts = getAgencyProductsList(authToken.UserID, origin)
	case models.UserRoleAgent:
		roleProducts = getAgentProductsList(authToken.UserID, origin)
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

func getAgencyProductsList(agencyUid, origin string) []models.Product {
	var (
		agency       models.Agency
		productsList = make([]models.Product, 0)
	)

	log.Println("GetAgentProducts")

	fireAgency := lib.GetDatasetByEnv(origin, models.AgencyCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgency, agencyUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&agency)
	lib.CheckError(err)

	defaultAgencyProduct := getDefaultProductsByChannel(models.UserRoleAgency + "/")

	for _, product := range agency.Products {
		if !product.IsAgencyActive {
			continue
		}
		for _, defaultProduct := range defaultAgencyProduct {
			isProductActive := product.Name == defaultProduct.Name && defaultProduct.IsAgencyActive
			if isProductActive {
				productsList = append(productsList, product)
				break
			}
		}
	}

	return productsList
}

func getAgentProductsList(agentUid, origin string) []models.Product {
	var (
		agent models.Agent
		//agency       models.Agency
		productsList = make([]models.Product, 0)
	)

	log.Println("GetAgentProducts")

	fireAgent := lib.GetDatasetByEnv(origin, models.AgentCollection)
	docsnap, err := lib.GetFirestoreErr(fireAgent, agentUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&agent)
	lib.CheckError(err)

	defaultAgentProducts := getDefaultProductsByChannel(models.UserRoleAgent + "/")

	for _, product := range agent.Products {
		if !product.IsAgentActive {
			continue
		}
		for _, defaultProduct := range defaultAgentProducts {
			isProductActive := product.Name == defaultProduct.Name && defaultProduct.IsAgentActive
			if isProductActive {
				productsList = append(productsList, product)
				break
			}
		}
	}

	/*if agent.AgencyUid != "" {
		agencyProductsList := getAgencyProductsList(agent.AgencyUid, origin)
		err = docsnap.DataTo(&agency)
		for index, product := range productsList {
			for _, agencyProduct := range agencyProductsList {
				hasToBeRemoved := product.Name == agencyProduct.Name && !agency.IsActive
				if hasToBeRemoved {
					productsList = append(productsList[index:], productsList[index+1:]...)
				}
			}
		}
	}*/

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
