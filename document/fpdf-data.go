package document

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

type slugStruct struct {
	name  string
	order int64
}

func loadLifeBeneficiariesInfo(policy *models.Policy) ([]map[string]string, string, string) {
	legitimateSuccessorsChoice := "X"
	designatedSuccessorsChoice := ""
	beneficiaries := []map[string]string{
		{
			"name":           "=====",
			"fiscalCode":     "=====",
			"address":        "=====",
			"mail":           "=====",
			"phone":          "=====",
			"relation":       "=====",
			"contactConsent": "=====",
		},
		{
			"name":           "=====",
			"fiscalCode":     "=====",
			"address":        "=====",
			"mail":           "=====",
			"phone":          "=====",
			"relation":       "=====",
			"contactConsent": "=====",
		},
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")
	lib.CheckError(err)

	if deathGuarantee.Beneficiaries != nil && (*deathGuarantee.Beneficiaries)[0].BeneficiaryType != "legalAndWillSuccessors" {
		legitimateSuccessorsChoice = ""
		designatedSuccessorsChoice = "X"

		for index, beneficiary := range *deathGuarantee.Beneficiaries {
			address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
				" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
				" (" + beneficiary.Residence.CityCode + ")")
			beneficiaries[index]["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
			beneficiaries[index]["fiscalCode"] = strings.ToUpper(beneficiary.FiscalCode)
			beneficiaries[index]["address"] = address
			beneficiaries[index]["mail"] = beneficiary.Mail
			beneficiaries[index]["phone"] = beneficiary.Phone
			if beneficiary.IsFamilyMember {
				beneficiaries[index]["relation"] = "Nucleo familiare (rapporto di parentela, coniuge, unione civile, " +
					"convivenza more uxorio)"
			} else {
				beneficiaries[index]["relation"] = "Altro (no rapporto parentela)"
			}
			if beneficiary.IsContactable {
				beneficiaries[index]["contactConsent"] = "SI"
			} else {
				beneficiaries[index]["contactConsent"] = "NO"
			}
		}
	}

	return beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice
}

func loadProducerInfo(origin string, policy *models.Policy) map[string]string {
	policyProducer := map[string]string{
		"name":            "LOMAZZI MICHELE",
		"ruiSection":      "A",
		"ruiCode":         "A000703480",
		"ruiRegistration": "02.03.2022",
	}

	if policy.ProducerUid == "" || strings.EqualFold(policy.ProducerType, models.PartnershipProducerType) {
		jsonProducer, _ := json.Marshal(policyProducer)
		log.Printf("[loadProducerInfo] producer info %s", string(jsonProducer))
		return policyProducer
	}

	log.Printf("[loadProducerInfo] loading producer %s data from Firestore...", policy.ProducerUid)

	networkNode, err := network.GetNodeByUid(policy.ProducerUid)
	lib.CheckError(err)

	log.Printf("[loadProducerInfo] setting producer %s info, producerType %s", policy.ProducerUid, policy.ProducerType)

	switch networkNode.Type {
	case models.AgentProducerType:
		policyProducer["name"] = strings.ToUpper(networkNode.Agent.Surname) + " " + strings.ToUpper(networkNode.Agent.Name)
		policyProducer["ruiSection"] = networkNode.Agent.RuiSection
		policyProducer["ruiCode"] = networkNode.Agent.RuiCode
		policyProducer["ruiRegistration"] = networkNode.Agent.RuiRegistration.Format("02.01.2006")
	case models.AgencyProducerType:
		policyProducer["name"] = strings.ToUpper(networkNode.Agency.Name)
		policyProducer["ruiSection"] = networkNode.Agency.RuiSection
		policyProducer["ruiCode"] = networkNode.Agency.RuiCode
		policyProducer["ruiRegistration"] = networkNode.Agency.RuiRegistration.Format("02.01.2006")
	}

	jsonProducer, _ := json.Marshal(policyProducer)
	log.Printf("[loadProducerInfo] producer info %s", string(jsonProducer))

	return policyProducer
}

func loadLifeGuarantees(policy *models.Policy) (map[string]map[string]string, []slugStruct) {
	const (
		death               = "death"
		permanentDisability = "permanent-disability"
		temporaryDisability = "temporary-disability"
		seriousIll          = "serious-ill"
	)
	var (
		guaranteesMap map[string]map[string]string
		slugs         []slugStruct
	)
	lifeProduct, err := product.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	lib.CheckError(err)

	guaranteesMap = make(map[string]map[string]string, 0)

	for guaranteeSlug, guarantee := range lifeProduct.Companies[0].GuaranteesMap {
		guaranteesMap[guaranteeSlug] = make(map[string]string, 0)

		guaranteesMap[guaranteeSlug]["name"] = guarantee.CompanyName
		guaranteesMap[guaranteeSlug]["sumInsuredLimitOfIndemnity"] = "====="
		guaranteesMap[guaranteeSlug]["duration"] = "=="
		guaranteesMap[guaranteeSlug]["endDate"] = "==="
		guaranteesMap[guaranteeSlug]["price"] = "===="
		if guaranteeSlug != death {
			guaranteesMap[guaranteeSlug]["price"] += " (*)"
		}
		slugs = append(slugs, slugStruct{name: guaranteeSlug, order: guarantee.Order})
	}

	sort.Slice(slugs, func(i, j int) bool {
		return slugs[i].order < slugs[j].order
	})

	for _, guarantee := range policy.GuaranteesToMap() {
		var price float64
		guaranteesMap[guarantee.Slug]["sumInsuredLimitOfIndemnity"] = humanize.FormatFloat("#.###,",
			guarantee.Value.SumInsuredLimitOfIndemnity) + " €"
		guaranteesMap[guarantee.Slug]["duration"] = strconv.Itoa(guarantee.Value.Duration.Year)
		guaranteesMap[guarantee.Slug]["endDate"] = policy.StartDate.AddDate(guarantee.Value.Duration.Year, 0, 0).Format(dateLayout)
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			price = guarantee.Value.PremiumGrossMonthly * 12
		} else {
			price = guarantee.Value.PremiumGrossYearly
		}
		guaranteesMap[guarantee.Slug]["price"] = humanize.FormatFloat("#.###,##", price) + " €"
		if guarantee.Slug != death {
			guaranteesMap[guarantee.Slug]["price"] += " (*)"
		}
	}
	return guaranteesMap, slugs
}

func loadPersonaGuarantees(policy *models.Policy) (map[string]map[string]string, []slugStruct) {
	var (
		guaranteesMap map[string]map[string]string
		slugs         []slugStruct
	)
	personaProduct, err := product.GetProduct(policy.Name, policy.ProductVersion, models.UserRoleAdmin)
	lib.CheckError(err)

	guaranteesMap = make(map[string]map[string]string, 0)
	offerName := policy.OfferlName

	for guaranteeSlug, guarantee := range personaProduct.Companies[0].GuaranteesMap {
		guaranteesMap[guaranteeSlug] = make(map[string]string, 0)

		guaranteesMap[guaranteeSlug]["name"] = guarantee.CompanyName
		guaranteesMap[guaranteeSlug]["sumInsuredLimitOfIndemnity"] = "====="
		guaranteesMap[guaranteeSlug]["details"] = "====="
		guaranteesMap[guaranteeSlug]["price"] = "====="
		slugs = append(slugs, slugStruct{name: guaranteeSlug, order: guarantee.Order})
	}

	sort.Slice(slugs, func(i, j int) bool {
		return slugs[i].order < slugs[j].order
	})

	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			var price float64
			var details string

			guaranteesMap[guarantee.Slug]["sumInsuredLimitOfIndemnity"] = humanize.FormatFloat("#.###,", guarantee.Offer[offerName].SumInsuredLimitOfIndemnity) + " €"
			if policy.PaymentSplit == string(models.PaySplitMonthly) {
				price = guarantee.Value.PremiumGrossMonthly * 12
			} else {
				price = guarantee.Value.PremiumGrossYearly
			}
			guaranteesMap[guarantee.Slug]["price"] = humanize.FormatFloat("#.###,##", price) + " €"

			switch guarantee.Slug {
			case "IPI":
				details = "Franchigia " + guarantee.Value.Deductible + guarantee.Value.DeductibleUnit
				if guarantee.Value.DeductibleType == "absolute" {
					details += " Assoluta"
				} else {
					details += " Assorbibile"
				}
			case "D":
				details = "Beneficiari:\n"
				for _, beneficiary := range *guarantee.Beneficiaries {
					if beneficiary.BeneficiaryType != "chosenBeneficiary" {
						details += personaProduct.Companies[0].GuaranteesMap["D"].BeneficiaryOptions[beneficiary.
							BeneficiaryType]
						break
					}
					details += beneficiary.Name + " " + beneficiary.Surname + " " + beneficiary.FiscalCode + "\n"

				}
			case "ITI":
				details = "Franchigia " + guarantee.Value.Deductible + " " + guarantee.Offer[offerName].DeductibleUnit
			default:
				details = "====="
			}
			guaranteesMap[guarantee.Slug]["details"] = details
		}
	}

	return guaranteesMap, slugs
}
