package partnership

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/golang-jwt/jwt/v4"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/quote"
	"github.com/wopta/goworkspace/user"
)

func init() {
	log.Println("INIT Partnership")
	functions.HTTP("Partnership", Partnership)
}

func Partnership(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/life/:partnershipUid",
				Handler: LifePartnershipFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}

func LifePartnershipFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response PartnershipResponse
	resp.Header().Set("Access-Control-Allow-Methods", "GET")

	partnershipUid := strings.ToLower(r.Header.Get("partnershipUid"))
	jwtData := r.URL.Query().Get("jwt")

	policy, product, node, err := LifePartnership(partnershipUid, jwtData)

	response.Policy = policy
	response.Product = product
	response.Partnership = *node.Partnership

	responseJson, err := json.Marshal(response)

	return string(responseJson), response, err
}

func LifePartnership(partnershipUid, jwtData string) (models.Policy, models.Product, models.NetworkNode, error) {
	var (
		policy      models.Policy
		productLife models.Product
		node        models.NetworkNode
	)

	partnershipNode, err := network.GetNodeByUid(partnershipUid)

	if err != nil {
		return policy, productLife, node, err
	}

	partnershipName := partnershipNode.Partnership.Name

	products := partnershipNode.Products
	productLife = getLatestLifeProduct(products)

	ecommerceProducts := product.GetAllProductsByChannel(models.ECommerceChannel)
	latestLifeProduct := getLatestLifeProduct(ecommerceProducts)

	policy.Name = productLife.Name
	policy.NameDesc = *productLife.NameDesc
	policy.ProductVersion = productLife.Version
	policy.Company = productLife.Companies[0].Name
	policy.ProducerUid = partnershipUid
	policy.ProducerCode = partnershipName
	policy.PartnershipName = partnershipName
	policy.ProducerType = partnershipNode.Type

	switch partnershipName {
	case models.PartnershipBeProf:
		err = beProfPartnership(jwtData, &policy, &latestLifeProduct)
	case models.PartnershipFacile:
		err = facilePartnership(jwtData, &policy, &latestLifeProduct)
	}

	if err != nil {
		return policy, productLife, node, err
	}

	policy, err = quote.Life(models.UserRoleCustomer, policy)

	if err != nil {
		return policy, productLife, node, err
	}

	return policy, productLife, node, err
}

func getLatestLifeProduct(products []models.Product) models.Product {
	products = lib.SliceFilter(products, func(product models.Product) bool {
		return product.Name == "life"
	})
	sort.Slice(products, func(i, j int) bool {
		return products[i].Version > products[j].Version
	})
	latestLifeProduct := products[0]
	return latestLifeProduct
}

func addDefaultGuarantees(asset *models.Asset, guarantees *map[string]*models.Guarante, offerName string) {
	for _, guarantee := range *guarantees {
		guarantee.IsSelected = guarantee.IsSelected || guarantee.IsMandatory

		if !guarantee.IsSelected {
			continue
		}

		guarantee.Value = &models.GuaranteValue{
			SumInsuredLimitOfIndemnity: guarantee.Offer[offerName].SumInsuredLimitOfIndemnity,
			Duration:                   guarantee.Offer[offerName].Duration,
		}
		asset.Guarantees = append(asset.Guarantees, *guarantee)
	}
}

func beProfPartnership(jwtData string, policy *models.Policy, product *models.Product) error {
	var (
		person models.User
		asset  models.Asset
	)

	token, err := jwt.ParseWithClaims(jwtData, &BeprofClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key, e := b64.StdEncoding.DecodeString(os.Getenv("BEPROF_SIGNING_KEY"))

		return []byte(key), e
	})

	if claims, ok := token.Claims.(*BeprofClaims); ok && token.Valid {
		person.Name = claims.UserFirstname
		person.Surname = claims.UserLastname
		person.Mail = claims.UserEmail
		person.FiscalCode = claims.UserFiscalcode
		person.Address = claims.UserAddress
		person.PostalCode = claims.UserPostalcode
		person.City = claims.UserCity
		person.CityCode = claims.UserMunicipalityCode
		person.Work = claims.UserEmploymentSector
		person.VatCode = claims.UserPiva

		if _, personData, err := user.ExtractUserDataFromFiscalCode(person); err == nil {
			person = personData
		}

		policy.Contractor = person
		asset.Person = &person
		policy.OfferlName = "default"

		addDefaultGuarantees(&asset, &product.Companies[0].GuaranteesMap, policy.OfferlName)

		policy.Assets = append(policy.Assets, asset)
		policy.PartnershipData = claims.ToMap()

		return nil
	}

	log.Printf("[beProfPartnership] could not validate beprof partnership JWT - %s", err.Error())
	return err
}

func facilePartnership(jwtData string, policy *models.Policy, product *models.Product) error {
	var (
		person models.User
		asset  models.Asset
	)

	token, err := jwt.ParseWithClaims(jwtData, &FacileClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := os.Getenv("FACILE_SIGNING_KEY")

		return []byte(key), nil
	})

	if claims, ok := token.Claims.(*FacileClaims); ok && token.Valid {
		person.Name = claims.CustomerName
		person.Surname = claims.CustomerFamilyName
		person.Mail = claims.Email
		birthDate, _ := time.Parse(models.TimeDateOnly, claims.CustomerBirthDate)
		person.BirthDate = birthDate.Format(time.RFC3339)
		person.Phone = claims.Mobile
		person.Gender = claims.Gender
		policy.Contractor = person
		asset.Person = &person
		policy.OfferlName = "default"

		addDefaultGuarantees(&asset, &product.Companies[0].GuaranteesMap, policy.OfferlName)

		for index, guarantee := range asset.Guarantees {
			if guarantee.Slug == "death" {
				asset.Guarantees[index].Value.Duration.Year = claims.Duration
				asset.Guarantees[index].Value.SumInsuredLimitOfIndemnity = float64(claims.InsuredCapital)
			}
		}

		policy.Assets = append(policy.Assets, asset)
		policy.PartnershipData = claims.ToMap()
	}

	log.Printf("[beProfPartnership] could not validate beprof partnership JWT - %s", err.Error())
	return err
}

type BeprofClaims struct {
	UserBeprofid         int    `json:"user.beprofid"`
	UserFirstname        string `json:"user.firstname"`
	UserLastname         string `json:"user.lastname"`
	UserEmail            string `json:"user.email"`
	UserMobile           string `json:"user.mobile"`
	UserFiscalcode       string `json:"user.fiscalcode"`
	UserPiva             string `json:"user.piva"`
	UserProvince         string `json:"user.province"`
	UserCity             string `json:"user.city"`
	UserPostalcode       string `json:"user.postalcode"`
	UserAddress          string `json:"user.address"`
	UserMunicipalityCode string `json:"user.municipality_code"`
	UserEmploymentSector string `json:"user.employment_sector"`
	ProductCode          string `json:"product.code"`
	ProductPurchaseid    string `json:"product.purchaseid"`
	Price                string `json:"price"`
	jwt.RegisteredClaims
}

func (facileClaims BeprofClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(facileClaims)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m

}

type FacileClaims struct {
	CustomerName       string `json:"customerName"`
	CustomerFamilyName string `json:"customerFamilyName"`
	CustomerBirthDate  string `json:"customerBirthDate"`
	Gender             string `json:"gender"`
	Email              string `json:"email"`
	Mobile             string `json:"mobile"`
	IsSmoker           bool   `json:"isSmoker"`
	InsuredCapital     int    `json:"insuredCapital"`
	Duration           int    `json:"duration"`
	jwt.RegisteredClaims
}

func (facileClaims FacileClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(facileClaims)
	lib.CheckError(err)

	err = json.Unmarshal(b, &m)
	lib.CheckError(err)

	return m
}

type PartnershipResponse struct {
	Policy      models.Policy          `json:"policy"`
	Partnership models.PartnershipNode `json:"partnership"`
	Product     models.Product         `json:"product"`
}
