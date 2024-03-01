package partnership

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

func LifePartnershipFx(resp http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response PartnershipResponse
	resp.Header().Set("Access-Control-Allow-Methods", "GET")
	resp.Header().Set("Content-Type", "application/json")

	log.Println("[LifePartnershipFx] handler start ----------------------------")

	partnershipUid := strings.ToLower(r.Header.Get("partnershipUid"))
	jwtData := r.URL.Query().Get("jwt")

	log.Printf("[LifePartnershipFx] partnershipUid: %s jwt: %s", partnershipUid, jwtData)

	policy, product, node, err := LifePartnership(partnershipUid, jwtData, r.Header.Get("Origin"))
	if err != nil {
		log.Printf("[LifePartnershipFx] error: %s", err.Error())
		return "", response, err
	}

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
		log.Println("[LifePartnership] call beProfLifePartnership function")
		err = beProfLifePartnership(jwtData, &policy, productLife, partnershipNode)
	case models.PartnershipFacile:
		log.Println("[LifePartnership] call facileLifePartnership function")
		err = facileLifePartnership(jwtData, &policy, productLife, partnershipNode)
	case models.PartnershipFpinsurance:
		log.Println("[LifePartnership] call fpinsuranceLifePartnership function")
		err = fpinsuranceLifePartnership(jwtData, &policy, productLife)
	default:
		log.Printf("[LifePartnership] could not find partnership with name %s", partnershipName)
		err = fmt.Errorf("invalid partnership name: %s", partnershipName)
	}

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
