package mga

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
)

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", Mga)
}

func Mga(w http.ResponseWriter, r *http.Request) {
	log.Println("Mga")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/products",
				Handler: func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) { return "", nil, nil },
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/journey/:product",
				Handler: GetProductJourneyByEntitlementFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}

	route.Router(w, r)
}

const (
	channelEcommerce string = "e-commerce"
	channelAgency    string = "agency"
	channelAgent     string = "agent"
	channelMga       string = "mga"
)

func GetProductJourneyByEntitlementFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		response models.Product
		err      error
	)

	idToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	authToken, err := GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	productName := r.Header.Get("product")

	// how to handle version ?
	response, err = GetProductJourneyByEntitlement(productName, "v1", authToken)
	lib.CheckError(err)

	jsonOut, err := json.Marshal(response)

	return string(jsonOut), response, err
}

func GetProductJourneyByEntitlement(productName, version string, authToken AuthToken) (models.Product, error) {
	log.Println("GetProductJourneyByEntitlement")
	var (
		responseProduct models.Product
	)

	mgaProduct, err := product.GetMgaProduct(productName, "v1")
	lib.CheckError(err)

	if authToken.Role == models.UserRoleAdmin || authToken.Role == models.UserRoleManager {
		return mgaProduct, nil
	}

	if authToken.Role == models.UserRoleAll {
		if mgaProduct.IsEcommerceActive {
			ecomProduct, err := product.GetProduct(productName, "v1", "")
			lib.CheckError(err)

			return ecomProduct, nil
		}

		return responseProduct, errors.New("product not active")
	}

	// Nesting hell!
	if authToken.Role == models.UserRoleAgency {
		if mgaProduct.IsAgencyActive {
			responseProduct = mgaProduct
			// get user product and check for overrides
			agency, err := GetAgencyByAuthId(authToken.UserID)
			lib.CheckError(err)

			// agencyProduct := sliceFind(
			// 	agency.Products,
			// 	func(p models.Product) bool { return p.Name == productName },
			// )
			agencyProduct := GetProductByName(agency.Products, productName)
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
			}
			return responseProduct, nil
		}
		return responseProduct, errors.New("Deactivated")
	}

	if authToken.Role == models.UserRoleAgent {
		if mgaProduct.IsAgentActive {
			responseProduct = mgaProduct
			// get user product and check for overrides
			agent, err := GetAgentByAuthId(authToken.UserID)
			lib.CheckError(err)
			agency, err := GetAgencyByAuthId(agent.AgencyUid)

			agencyProduct := GetProductByName(agency.Products, productName)
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
			}

			agentProduct := GetProductByName(agent.Products, productName)
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
			}

			return responseProduct, nil
		}
		return responseProduct, errors.New("Deactivated")
	}

	return responseProduct, errors.New("no product found")
}

type AuthToken struct {
	Role   string `json:"role"`
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func GetAuthTokenFromIdToken(idToken string) (AuthToken, error) {
	if idToken == "" {
		return AuthToken{
			Role:   models.UserRoleAll,
			UserID: "",
			Email:  "",
		}, nil
	}

	token, err := lib.VerifyUserIdToken(idToken)
	if err != nil {
		log.Println("GetUserRoleFromToken: token err")
		return AuthToken{}, err
	}

	return AuthToken{
		Role:   token.Claims["role"].(string),
		UserID: token.Claims["user_id"].(string),
		Email:  token.Claims["email"].(string),
	}, nil
}

func GetProductChannelFromRole(role string) string {
	var channel string

	switch role {
	case models.UserRoleAgency:
		channel = channelAgency
	case models.UserRoleAgent:
		channel = channelAgent
	case models.UserRoleAdmin, models.UserRoleManager:
		channel = channelMga
	case models.UserRoleCustomer, models.UserRoleAll:
		channel = ""
	}

	return channel
}

const (
	agentCollection  string = "agents"
	agencyCollection string = "agencies"
)

// move to correct domain - models?
func GetAgentByAuthId(authId string) (models.Agent, error) {
	var agent models.Agent

	agentFirebase := lib.WhereLimitFirestore(agentCollection, "user.authId", "==", authId, 1)
	agent, err := FirestoreDocumentToAgent(agentFirebase)

	return agent, err
}

func GetAgencyByAuthId(authId string) (models.Agency, error) {
	var agent models.Agency

	agencyFirebase := lib.WhereLimitFirestore(agencyCollection, "uid", "==", authId, 1)
	agent, err := FirestoreDocumentToAgency(agencyFirebase)

	return agent, err
}

func FirestoreDocumentToAgent(query *firestore.DocumentIterator) (models.Agent, error) {
	var result models.Agent
	agentDocumentSnapshot, err := query.Next()

	if err == iterator.Done && agentDocumentSnapshot == nil {
		log.Println("Agent not found in firebase DB")
		return result, fmt.Errorf("no agent found")
	}

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get agent`)
		return result, err
	}

	e := agentDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.User.Uid = agentDocumentSnapshot.Ref.ID
	}

	return result, e
}

func FirestoreDocumentToAgency(query *firestore.DocumentIterator) (models.Agency, error) {
	var result models.Agency
	agencyDocumentSnapshot, err := query.Next()

	if err == iterator.Done && agencyDocumentSnapshot == nil {
		log.Println("Agency not found in firebase DB")
		return result, fmt.Errorf("no agent found")
	}

	if err != iterator.Done && err != nil {
		log.Println(`error happened while trying to get agency`)
		return result, err
	}

	e := agencyDocumentSnapshot.DataTo(&result)
	if len(result.Uid) == 0 {
		result.Uid = agencyDocumentSnapshot.Ref.ID
	}

	return result, e
}

func sliceFind[T any](slice []T, cb func(T) bool) (ret T) {
	for _, item := range slice {
		if cb(item) {
			ret = item
		}
	}
	return
}

func GetProductByName(products []models.Product, productName string) *models.Product {
	mapProduct := map[string]models.Product{}
	for _, p := range products {
		mapProduct[p.Name] = p
	}
	if p, ok := mapProduct[productName]; ok {
		return &p
	}
	return nil
}
