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
	Name    string `json:"name"`
	Company string `json:"company"`
	Logo    string `json:"logo"`
}

func GetProductsListByEntitlementFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	var (
		err          error
		userRole     string
		productsList = make([]GetProductListResp, 0)
	)
	log.Println("GetProductsListByEntitlement")

	idToken := r.Header.Get("Authorization")

	if idToken == "" {
		productsList = getEcommerceProductsList()
		output, err := json.Marshal(productsList)
		return string(output), productsList, err
	}

	userRole, err = lib.GetUserRoleFromIdToken(idToken)
	lib.CheckError(err)

	switch userRole {
	case models.UserRoleAdmin, models.UserRoleManager:
		productsList = getMgaProductsList()
	case models.UserRoleAgency:
		productsList = getAgencyProductsList()
	case models.UserRoleAgent:
		productsList = getAgentProductsList()
	}

	jsonOut, err := json.Marshal(productsList)

	return string(jsonOut), productsList, err
}

func getMgaProductsList() []GetProductListResp {
	productsList := make([]GetProductListResp, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "agent")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		for _, company := range product.Companies {
			productsList = append(productsList, GetProductListResp{
				Name:    product.Name,
				Company: company.Name,
				Logo:    "",
			})
		}
	}
	return productsList
}

func getAgencyProductsList() []GetProductListResp {
	productsList := make([]GetProductListResp, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "agent")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.IsAgentActive {
			for _, company := range product.Companies {
				if company.IsAgencyActive {
					productsList = append(productsList, GetProductListResp{
						Name:    product.Name,
						Company: company.Name,
						Logo:    "",
					})
				}
			}
		}
	}
	return productsList
}

func getAgentProductsList() []GetProductListResp {
	productsList := make([]GetProductListResp, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "agent")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.IsAgentActive {
			for _, company := range product.Companies {
				if company.IsAgentActive {
					productsList = append(productsList, GetProductListResp{
						Name:    product.Name,
						Company: company.Name,
						Logo:    "",
					})
				}
			}
		}
	}
	return productsList
}

func getEcommerceProductsList() []GetProductListResp {
	productsList := make([]GetProductListResp, 0)
	res := lib.GetFolderContentByEnv(pathPrefix + "e-commerce")
	for _, file := range res {
		var product models.Product
		err := json.Unmarshal(file, &product)
		lib.CheckError(err)
		if product.IsEcommerceActive {
			for _, company := range product.Companies {
				if company.IsEcommerceActive {
					productsList = append(productsList, GetProductListResp{
						Name:    product.Name,
						Company: company.Name,
						Logo:    "",
					})
				}
			}
		}
	}
	return productsList
}
