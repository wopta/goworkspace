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

	mgaProduct, err := product.GetMgaProduct(productName, "v1")
	lib.CheckError(err)

	switch authToken.Role {
	case models.UserRoleAdmin, models.UserRoleManager:
		responseProduct, err = getMgaProduct(&mgaProduct)
	case models.UserRoleAll:
		responseProduct, err = getEcommerceProduct(&mgaProduct, productName)
	case models.UserRoleAgency:
		responseProduct, err = getAgencyProduct(&mgaProduct, productName, authToken.UserID)
	case models.UserRoleAgent:
		responseProduct, err = getAgentProduct(&mgaProduct, productName, authToken.UserID)
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

func getMgaProduct(mgaProduct *models.Product) (*models.Product, error) {
	return mgaProduct, nil
}

func getEcommerceProduct(mgaProduct *models.Product, productName string) (*models.Product, error) {
	if !mgaProduct.IsEcommerceActive {
		return productNotActive()
	}

	ecomProduct, err := product.GetProduct(productName, "v1", "")

	return &ecomProduct, err
}

func getAgencyProduct(mgaProduct *models.Product, productName, agencyUid string) (*models.Product, error) {
	if !mgaProduct.IsAgencyActive {
		return productNotActive()
	}

	agencyDefaultProduct, err := product.GetProduct(productName, "v1", models.UserRoleAgency)
	lib.CheckError(err)
	responseProduct := &agencyDefaultProduct
	log.Printf("Agency Product Start: %v", responseProduct)
	agency, err := models.GetAgencyByAuthId(agencyUid)
	lib.CheckError(err)

	agencyProduct := getProductByName(agency.Products, productName)
	if agencyProduct != nil {
		if len(agencyProduct.Steps) > 0 {
			responseProduct.Steps = agencyProduct.Steps
		}

		for _, c := range agencyProduct.Companies {
			for _, c2 := range responseProduct.Companies {
				if c2.Name == c.Name {
					c2.Mandate = c.Mandate
				}
			}
		}
		log.Printf("Agency Product Modified: %v", responseProduct)
	}

	log.Printf("Agency Product Response: %v", responseProduct)
	return responseProduct, nil
}

func getAgentProduct(mgaProduct *models.Product, productName, agentUid string) (*models.Product, error) {
	if !mgaProduct.IsAgentActive {
		return productNotActive()
	}

	agentDefaultProduct, err := product.GetProduct(productName, "v1", models.UserRoleAgent)
	lib.CheckError(err)
	responseProduct := &agentDefaultProduct
	log.Printf("Agent Product Start: %v", responseProduct)
	agent, err := models.GetAgentByAuthId(agentUid)
	lib.CheckError(err)
	agency, err := models.GetAgencyByAuthId(agent.AgencyUid)
	lib.CheckError(err)

	// TODO: traverse network
	agencyProduct := getProductByName(agency.Products, productName)
	if agencyProduct != nil {
		if len(agencyProduct.Steps) > 0 {
			responseProduct.Steps = agencyProduct.Steps
		}

		for _, c := range agencyProduct.Companies {
			for _, c2 := range responseProduct.Companies {
				if c2.Name == c.Name {
					c2.Mandate = c.Mandate
				}
			}
		}
		log.Printf("Agent product modified by agency: %v", responseProduct)
	}

	agentProduct := getProductByName(agent.Products, productName)
	if agentProduct != nil {
		if len(agentProduct.Steps) > 0 {
			responseProduct.Steps = agentProduct.Steps
		}

		for _, c := range agentProduct.Companies {
			for _, c2 := range responseProduct.Companies {
				if c2.Name == c.Name {
					c2.Mandate = c.Mandate
				}
			}
		}
		log.Printf("Agent product modified by agent: %v", responseProduct)
	}

	log.Printf("Agent Product Response: %v", responseProduct)
	return responseProduct, nil
}
