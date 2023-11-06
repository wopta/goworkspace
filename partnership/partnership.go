package partnership

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mohae/deepcopy"
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
	resp.Header().Set("Content-Type", "application/json")

	log.Println("[LifePartnershipFx] handler start ----------------------------")

	partnershipUid := strings.ToLower(r.Header.Get("partnershipUid"))
	jwtData := r.URL.Query().Get("jwt")

	log.Printf("[LifePartnershipFx] partnershipUid: %s jwt: %s", partnershipUid, jwtData)

	policy, product, node, err := LifePartnership(partnershipUid, jwtData, r.Header.Get("Origin"))

	response.Policy = policy
	response.Product = product
	response.Partnership = *node.Partnership

	responseJson, err := json.Marshal(response)

	log.Printf("[LifePartnershipFx] response: %s", string(responseJson))

	return string(responseJson), response, err
}

func LifePartnership(partnershipUid, jwtData, origin string) (models.Policy, models.Product, *models.NetworkNode, error) {
	var (
		policy          models.Policy
		productLife     *models.Product
		partnershipNode *models.NetworkNode
		err             error
	)

	log.Printf("[LifePartnership]")

	if partnershipNode, err = network.GetNodeByUid(partnershipUid); err != nil {
		return policy, *productLife, partnershipNode, err
	}

	if partnershipNode == nil {
		log.Printf("[LifePartnership] no partnership found")
		return policy, *productLife, partnershipNode, fmt.Errorf("no partnership found")
	}

	if !partnershipNode.IsActive {
		log.Printf("[LifePartnership] partnership is not active")
		return policy, *productLife, partnershipNode, fmt.Errorf("partnership is not active")
	}

	partnershipName := partnershipNode.Partnership.Name

	log.Printf("[LifePartnership] loading latest life product")

	warrant := partnershipNode.GetWarrant()
	productLife = product.GetLatestActiveProduct(models.LifeProduct, models.ECommerceChannel, partnershipNode, warrant)
	if productLife == nil {
		log.Printf("[LifePartnership] no product found")
		return policy, models.Product{}, partnershipNode, fmt.Errorf("no product found")
	}

	log.Printf("[LifePartnership] setting policy basic info")

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
		log.Println("[LifePartnership] call beProfPartnership function")
		err = beProfPartnership(jwtData, &policy, productLife)
	case models.PartnershipFacile:
		log.Println("[LifePartnership] call facilePartnership function")
		err = facilePartnership(jwtData, &policy, productLife)
	default:
		log.Printf("[LifePartnership] could not find partnership with name %s", partnershipName)
		err = fmt.Errorf("invalid partnership name: %s", partnershipName)
	}

	if err != nil {
		return policy, *productLife, partnershipNode, err
	}

	policy, err = quote.Life(policy, models.ECommerceChannel, partnershipNode, warrant)

	if err != nil {
		return policy, *productLife, partnershipNode, err
	}

	err = savePartnershipLead(&policy, partnershipNode, origin)

	return policy, *productLife, partnershipNode, err
}

func removeUnselectedGuarantees(policy *models.Policy) models.Policy {
	policyCopy := deepcopy.Copy(*policy).(models.Policy)
	for i, asset := range policy.Assets {
		policyCopy.Assets[i].Guarantees = lib.SliceFilter(asset.Guarantees, func(guarantee models.Guarante) bool {
			return guarantee.IsSelected
		})
	}
	return policyCopy
}

func savePartnershipLead(policy *models.Policy, node *models.NetworkNode, origin string) error {
	var err error

	log.Println("[savePartnershipLead] start --------------------------------------------")

	policyFire := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	policy.Channel = models.ECommerceChannel
	now := time.Now().UTC()

	policy.CreationDate = now
	policy.Status = models.PolicyStatusPartnershipLead
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	log.Printf("[savePartnershipLead] policy status %s", policy.Status)

	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = now

	log.Println("[savePartnershipLead] saving lead to firestore...")
	policyUid := lib.NewDoc(policyFire)
	policy.Uid = policyUid

	policyToSave := removeUnselectedGuarantees(policy)

	if err = lib.SetFirestoreErr(policyFire, policyUid, policyToSave); err != nil {
		return err
	}

	log.Println("[savePartnershipLead] saving lead to bigquery...")
	policyToSave.BigquerySave(origin)

	log.Println("[savePartnershipLead] end ----------------------------------------------")
	return err
}

func beProfPartnership(jwtData string, policy *models.Policy, product *models.Product) error {
	var (
		person models.User
		asset  models.Asset
	)

	log.Println("[beProfPartnership] decoding jwt")

	token, err := jwt.ParseWithClaims(jwtData, &BeprofClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key, e := b64.StdEncoding.DecodeString(os.Getenv("BEPROF_SIGNING_KEY"))

		return []byte(key), e
	})

	if claims, ok := token.Claims.(*BeprofClaims); ok && token.Valid {
		log.Println("[beProfPartnership] setting person info")
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

	log.Println("[facilePartnership] decoding jwt")

	token, err := jwt.ParseWithClaims(jwtData, &FacileClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := os.Getenv("FACILE_SIGNING_KEY")

		return []byte(key), nil
	})

	if claims, ok := token.Claims.(*FacileClaims); ok && token.Valid {
		log.Println("[facilePartnership] setting person info")
		person.Name = claims.CustomerName
		person.Surname = claims.CustomerFamilyName
		person.Mail = claims.Email
		birthDate, _ := time.Parse(models.TimeDateOnly, claims.CustomerBirthDate)
		person.BirthDate = birthDate.Format(time.RFC3339)
		person.Phone = fmt.Sprintf("+39%s", claims.Mobile)
		person.Gender = claims.Gender
		policy.Contractor = person
		asset.Person = &person
		policy.OfferlName = "default"

		log.Println("[facilePartnership] setting death guarantee info")

		deathGuarantee := product.Companies[0].GuaranteesMap["death"]
		deathGuarantee.Value = &models.GuaranteValue{
			Duration: &models.Duration{
				Year: claims.Duration,
			},
			SumInsuredLimitOfIndemnity: float64(claims.InsuredCapital),
		}
		asset.Guarantees = make([]models.Guarante, 0)
		asset.Guarantees = append(asset.Guarantees, *deathGuarantee)

		policy.Assets = append(policy.Assets, asset)
		policy.PartnershipData = claims.ToMap()
		return err
	}

	log.Printf("[facilePartnership] could not validate facile partnership JWT - %s", err.Error())
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

func (beprofClaims BeprofClaims) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(beprofClaims)
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
