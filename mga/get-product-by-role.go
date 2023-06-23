package mga

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

func GetProductByRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("GetProductByRoleFx")
	var (
		response models.Product
		err      error
	)

	productName := r.Header.Get("product")
	idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	// how to handle version ?
	response, err = GetProductByRole(productName, "v1", authToken)
	if err != nil {
		return "", response, err
	}
	jsonOut, err := json.Marshal(response)

	return string(jsonOut), response, err
}

func GetProductByRole(productName, version string, authToken models.AuthToken) (models.Product, error) {
	log.Println("GetProductByRole")
	var (
		responseProduct *models.Product
		err             error
	)

	switch authToken.Role {
	case models.UserRoleAdmin, models.UserRoleManager:
		responseProduct, err = getMgaProduct(productName)
	case models.UserRoleAll:
		responseProduct, err = getEcommerceProduct(productName)
	case models.UserRoleAgency:
		responseProduct, err = getAgencyProduct(productName, authToken.UserID)
	case models.UserRoleAgent:
		responseProduct, err = getAgentProduct(productName, authToken.UserID)
	default:
		responseProduct, err = productNotFound()
	}

	return *responseProduct, err
}

func getProductByName(products []models.Product, productName string) *models.Product {
	mapProduct := map[string]models.Product{}
	for _, p := range products {
		mapProduct[p.Name] = p
	}
	if p, ok := mapProduct[productName]; ok {
		return &p
	}
	return nil
}

func productNotActive() (*models.Product, error) {
	return nil, errors.New("product not active")
}

func productNotFound() (*models.Product, error) {
	return nil, errors.New("product not found")
}

func getMgaProduct(productName string) (*models.Product, error) {
	mgaProduct, err := product.GetMgaProduct(productName, "v1")
	lib.CheckError(err)

	return &mgaProduct, nil
}

func getEcommerceProduct(productName string) (*models.Product, error) {
	ecomProduct, err := product.GetProduct(productName, "v1", "")

	if !ecomProduct.IsEcommerceActive {
		return productNotActive()
	}

	return &ecomProduct, err
}

func getAgencyProduct(productName, agencyUid string) (*models.Product, error) {
	agencyDefaultProduct, err := product.GetProduct(productName, "v1", models.UserRoleAgency)
	lib.CheckError(err)

	if !agencyDefaultProduct.IsAgencyActive {
		return productNotActive()
	}

	responseProduct := &agencyDefaultProduct
	log.Printf("Agency Product Start: %v", responseProduct)
	agency, err := models.GetAgencyByAuthId(agencyUid)
	lib.CheckError(err)

	agencyProduct := getProductByName(agency.Products, productName)
	if agencyProduct == nil {
		return nil, errors.New("agency does not have product")
	}

	if !agencyProduct.IsAgencyActive {
		return productNotActive()
	}

	overrideProduct(responseProduct, agencyProduct)

	log.Printf("Agency Product Response: %v", responseProduct)
	return responseProduct, nil
}

func getAgentProduct(productName, agentUid string) (*models.Product, error) {
	agentDefaultProduct, err := product.GetProduct(productName, "v1", models.UserRoleAgent)
	lib.CheckError(err)

	if !agentDefaultProduct.IsAgentActive {
		return productNotActive()
	}

	responseProduct := &agentDefaultProduct
	log.Printf("Agent Product Start: %v", responseProduct)
	agent, err := models.GetAgentByAuthId(agentUid)
	lib.CheckError(err)
	agency, err := models.GetAgencyByAuthId(agent.AgencyUid)
	lib.CheckError(err)

	agentProduct := getProductByName(agent.Products, productName)
	if agentProduct == nil {
		return nil, errors.New("agent does not have product")
	}

	if !agentProduct.IsAgentActive {
		return productNotActive()
	}

	// TODO: traverse network
	agencyProduct := getProductByName(agency.Products, productName)
	if agencyProduct != nil {
		overrideProduct(responseProduct, agencyProduct)
		log.Printf("Agent product modified by agency: %v", responseProduct)
	}

	overrideProduct(responseProduct, agentProduct)
	log.Printf("Agent product modified by agent: %v", responseProduct)

	log.Printf("Agent Product Response: %v", responseProduct)
	return responseProduct, nil
}

func overrideProduct(baseProduct *models.Product, insertedProduct *models.Product) {
	if len(insertedProduct.Steps) > 0 {
		baseProduct.Steps = insertedProduct.Steps
	}

	for _, c := range insertedProduct.Companies {
		for _, c2 := range baseProduct.Companies {
			if c2.Name == c.Name {
				c2.Mandate = c.Mandate
			}
		}
	}
}
