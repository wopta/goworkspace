package mga

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

type GetProductByRoleRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Company string `json:"company"`
}

func GetProductByRoleFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("GetProductByRoleFx")
	var (
		resp       models.Product
		respString string
		request    GetProductByRoleRequest
		err        error
	)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &request)
	lib.CheckError(err)
	log.Printf("GetProductByRoleFx body: %s", string(body))

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	log.Printf("GetProductByRoleFx authToken: %s", authToken)

	resp, err = GetProductByRole(request.Name, request.Version, request.Company, authToken)
	if err != nil {
		return "", resp, err
	}
	jsonResp, err := json.Marshal(resp)

	respString = string(jsonResp)
	switch request.Name {
	case "persona":
		respString, resp, err = product.ReplaceDatesInProduct(resp, 75)
	case "life":
		respString, resp, err = product.ReplaceDatesInProduct(resp, 55)
	}

	log.Printf("GetProductByRoleFx response: %s", respString)
	return respString, resp, err
}

func GetProductByRole(productName, version, company string, authToken models.AuthToken) (models.Product, error) {
	log.Println("GetProductByRole")
	var (
		responseProduct *models.Product
		err             error
	)

	switch authToken.Role {
	case models.UserRoleAdmin, models.UserRoleManager:
		responseProduct, err = getMgaProduct(productName, version, company)
	case models.UserRoleAll, models.UserRoleCustomer:
		responseProduct, err = getEcommerceProduct(productName, version, company)
	case models.UserRoleAgency:
		responseProduct, err = getAgencyProduct(productName, version, company, authToken.UserID)
	case models.UserRoleAgent:
		responseProduct, err = getAgentProduct(productName, version, company, authToken.UserID)
	default:
		responseProduct, err = productNotFound()
	}

	return *responseProduct, err
}

func getProductByName(products []models.Product, productName string) *models.Product {
	log.Println("getProductByName")
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

func getMgaProduct(productName, version, company string) (*models.Product, error) {
	log.Println("getMgaProduct")
	mgaProduct, err := product.GetMgaProduct(productName, version)
	lib.CheckError(err)

	return &mgaProduct, nil
}

func getEcommerceProduct(productName, version, company string) (*models.Product, error) {
	log.Println("getEcommerceProduct")
	ecomProduct, err := product.GetProduct(productName, version, "")

	if !ecomProduct.IsEcommerceActive {
		return productNotActive()
	}

	return &ecomProduct, err
}

func getAgencyProduct(productName, version, company, agencyUid string) (*models.Product, error) {
	log.Println("getAgencyProduct")
	agencyDefaultProduct, err := product.GetProduct(productName, version, models.UserRoleAgency)
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

func getAgentProduct(productName, version, company, agentUid string) (*models.Product, error) {
	log.Println("getAgentProduct")
	agentDefaultProduct, err := product.GetProduct(productName, version, models.UserRoleAgent)
	lib.CheckError(err)

	if !agentDefaultProduct.IsAgentActive {
		return productNotActive()
	}

	responseProduct := &agentDefaultProduct
	log.Printf("Agent Product Start: %v", responseProduct)
	agent, err := models.GetAgentByAuthId(agentUid)
	lib.CheckError(err)
	agency, _ := models.GetAgencyByAuthId(agent.AgencyUid)

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
	log.Println("overrideProduct")
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
