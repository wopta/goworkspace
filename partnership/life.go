package partnership

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/quote"
	"github.com/wopta/goworkspace/user"
)

func LifePartnershipFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		response        PartnershipResponse
		partnershipNode *models.NetworkNode
		policy          models.Policy
		productLife     *models.Product
		claims          models.LifeClaims
		err             error
	)

	log.SetPrefix("[LifePartnershipFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	partnershipUid := strings.ToLower(chi.URLParam(r, "partnershipUid"))
	jwtData := r.URL.Query().Get("jwt")
	key := lib.ToUpper(fmt.Sprintf("%s_SIGNING_KEY", partnershipUid))

	log.Printf("partnershipUid: %s jwt: %s", partnershipUid, jwtData)

	if partnershipNode, err = network.GetNodeByUid(partnershipUid); err != nil {
		log.Printf("error getting node: %s", err.Error())
		return "", nil, err
	}

	if partnershipNode == nil {
		log.Printf("no partnership found")
		return "", nil, err
	}

	if !partnershipNode.IsActive {
		log.Printf("partnership is not active")
		return "", nil, err
	}

	log.Printf("loading latest life product")
	productLife = product.GetLatestActiveProduct(models.LifeProduct, models.ECommerceChannel, partnershipNode, nil)
	if productLife == nil {
		log.Printf("no product found")
		return "", nil, fmt.Errorf("no product found")
	}
	policy = setPolicyPartnershipInfo(policy, productLife, partnershipNode)

	if claims, err = partnershipNode.Partnership.DecryptJwtClaims(
		jwtData, os.Getenv(key), lifeClaimsExtractor(partnershipNode.Partnership)); err != nil {
		log.Printf("could not validate partnership JWT - %s", err.Error())
		return "", nil, err
	}

	if !claims.IsEmpty() {
		policy, err = setClaimsIntoPolicy(policy, productLife, claims)
		if err != nil {
			log.Printf("error extracting data from claims: %s", err.Error())
			return "", nil, err
		}

		if policy.Contractor.BirthDate != "" {
			quotedPolicy, err := quote.Life(policy, models.ECommerceChannel, partnershipNode, nil, models.ECommerceFlow)
			if err != nil {
				log.Printf("error quoting for partnership: %s", err.Error())
				return "", nil, err
			}
			policy = quotedPolicy
		}
	}

	err = savePartnershipLead(&policy, partnershipNode, "")
	if err != nil {
		log.Printf("error saving lead: %s", err.Error())
		return "", nil, err
	}

	response.Policy = policy
	response.Product = *productLife
	response.Partnership = PartnershipNode{partnershipNode.Partnership.Name, partnershipNode.Partnership.Skin}

	responseJson, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}

func setPolicyPartnershipInfo(policy models.Policy, product *models.Product, node *models.NetworkNode) models.Policy {
	policy.Name = product.Name
	policy.NameDesc = *product.NameDesc
	policy.ProductVersion = product.Version
	policy.Company = product.Companies[0].Name
	policy.ProducerUid = node.Uid
	policy.ProducerCode = node.Partnership.Name
	policy.PartnershipName = node.Partnership.Name
	policy.ProducerType = node.Type

	return policy
}

func lifeClaimsExtractor(node *models.PartnershipNode) func([]byte) (models.LifeClaims, error) {
	switch node.Name {
	case models.PartnershipBeProf:
		return beprofLifeClaimsExtractor
	case models.PartnershipFacile:
		return facileLifeClaimsExtractor
	case models.PartnershipELeads:
		return eLeadsLifeClaimsExtractor
	}

	if node.IsJwtProtected() {
		return defaultLifeClaimsExtractor
	}

	return func(b []byte) (models.LifeClaims, error) {
		return models.LifeClaims{}, nil
	}
}

func defaultLifeClaimsExtractor(b []byte) (models.LifeClaims, error) {
	claims := &models.LifeClaims{}
	err := json.Unmarshal(b, claims)
	if err != nil {
		return models.LifeClaims{}, err
	}

	data := make(map[string]interface{})
	dataBytes, _ := json.Marshal(b)
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return models.LifeClaims{}, err
	}

	claims.Data = data
	return *claims, nil
}

func facileLifeClaimsExtractor(b []byte) (models.LifeClaims, error) {
	facileClaims := &FacileClaims{}
	err := json.Unmarshal(b, facileClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}

	adapter := FacileLifeClaimsAdapter{
		facileClaims: facileClaims,
	}
	return adapter.ExtractClaims()
}

func beprofLifeClaimsExtractor(b []byte) (models.LifeClaims, error) {
	beprofClaims := &BeprofClaims{}
	err := json.Unmarshal(b, beprofClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}

	adapter := BeprofLifeClaimsAdapter{
		beprofClaims: beprofClaims,
	}
	return adapter.ExtractClaims()
}

func eLeadsLifeClaimsExtractor(b []byte) (models.LifeClaims, error) {
	eLeadsClaims := &ELeadsClaims{}
	err := json.Unmarshal(b, eLeadsClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}

	adapter := ELeadsLifeClaimsAdapter{
		eLeadsClaims: eLeadsClaims,
	}
	return adapter.ExtractClaims()
}

func setClaimsIntoPolicy(policy models.Policy, product *models.Product, claims models.LifeClaims) (models.Policy, error) {
	var (
		person models.User
		asset  models.Asset
	)

	person.Name = claims.Name
	person.Surname = claims.Surname
	person.BirthDate = claims.BirthDate
	person.Gender = claims.Gender
	person.Mail = claims.Email
	person.FiscalCode = claims.FiscalCode
	person.Address = claims.Address
	person.PostalCode = claims.Postalcode
	person.City = claims.City
	person.CityCode = claims.CityCode
	person.Work = claims.Work
	person.VatCode = claims.VatCode

	if person.FiscalCode != "" {
		_, personData, err := user.ExtractUserDataFromFiscalCode(person)
		if err != nil {
			return models.Policy{}, err
		}
		person = personData
	}

	person.Normalize()

	policy.Contractor = *person.ToContractor()
	asset.Person = &person
	policy.OfferlName = "default"

	if claims.Guarantees != nil {
		asset.Guarantees = make([]models.Guarante, 0)
		for slug, value := range claims.Guarantees {
			g := product.Companies[0].GuaranteesMap[slug]
			g.Value = &models.GuaranteValue{
				Duration: &models.Duration{
					Year: value.Duration,
				},
				SumInsuredLimitOfIndemnity: value.SumInsuredLimitOfIndemnity,
			}
			asset.Guarantees = append(asset.Guarantees, *g)
		}
	}

	policy.Assets = append(policy.Assets, asset)
	policy.PartnershipData = claims.Data

	return policy, nil
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

	policyFire := lib.GetDatasetByEnv(origin, lib.PolicyCollection)

	policy.Channel = lib.ECommerceChannel
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
