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
	Name     string `json:"name"`
	NameDesc string `json:"nameDesc"`
	Version  string `json:"version"`
	Company  string `json:"company"`
	Logo     string `json:"logo"`
}

func GetProductsListByEntitlementFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	var (
		err               error
		userRole, userUID string
		productsList      = make([]models.Product, 0)
	)
	log.Println("GetProductsListByEntitlement")

	origin := r.Header.Get("Origin")
	idToken := r.Header.Get("Authorization")

	if idToken == "" {
		productsList = getEcommerceProductsList()
		output, err := json.Marshal(productsList)
		lib.CheckError(err)
		return string(output), productsList, err
	}

	userRole, err = lib.GetUserRoleFromIdToken(idToken)
	lib.CheckError(err)
	userUID, err = lib.GetUserIdFromIdToken(idToken)
	lib.CheckError(err)

	switch userRole {
	case models.UserRoleAdmin, models.UserRoleManager:
		productsList = getMgaProductsList()
	case models.UserRoleAgency:
		productsList = getAgencyProductsList(userUID, origin)
	case models.UserRoleAgent:
		productsList = getAgentProductsList(userUID, origin)
	case models.UserRoleCustomer:
		productsList = getEcommerceProductsList()
	}

	list := make([]GetProductListResp, 0)
	for _, product := range productsList {
		for _, company := range product.Companies {
			list = append(list, GetProductListResp{
				Name:     product.Name,
				NameDesc: *product.NameDesc,
				Version:  product.Version,
				Company:  company.Name,
				Logo:     "",
			})
		}
	}

	jsonOut, err := json.Marshal(list)

	return string(jsonOut), list, err
}

func getMgaProductsList() []models.Product {
	productsList := make([]models.Product, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + models.UserRoleAgent)
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		productsList = append(productsList, product)
	}
	return productsList
}

func getEcommerceProductsList() []models.Product {
	productsList := make([]models.Product, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "e-commerce")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.IsEcommerceActive {
			productsList = append(productsList, product)
		}
	}
	return productsList
}

func getAgencyProductsList(agencyUid, origin string) []models.Product {
	var (
		agency       models.Agency
		productsList = make([]models.Product, 0)
	)

	log.Println("GetAgentProducts")

	fireAgency := lib.GetDatasetByEnv(origin, models.UserRoleAgency)
	docsnap, err := lib.GetFirestoreErr(fireAgency, agencyUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&agency)
	lib.CheckError(err)

	defaultAgencyProduct := getDefaultProductsByChannel(models.UserRoleAgency)

	for _, product := range agency.Products {
		if !product.IsAgencyActive {
			continue
		}
		for _, defaultProduct := range defaultAgencyProduct {
			isProductActive := product.Name == defaultProduct.Name && defaultProduct.IsAgentActive
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

	fireAgent := lib.GetDatasetByEnv(origin, models.UserRoleAgent)
	docsnap, err := lib.GetFirestoreErr(fireAgent, agentUid)
	lib.CheckError(err)
	err = docsnap.DataTo(&agent)
	lib.CheckError(err)

	defaultAgentProducts := getDefaultProductsByChannel(models.UserRoleAgent)

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
