package document

import (
	"github.com/dustin/go-humanize"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"sort"
	"strconv"
	"strings"
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
			"name":       "=====",
			"fiscalCode": "=====",
			"address":    "=====",
			"mail":       "=====",
			"phone":      "=====",
			"relation":   "=====",
			"consent":    "=====",
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

	if deathGuarantee.Beneficiaries != nil && !(*deathGuarantee.Beneficiaries)[0].IsLegitimateSuccessors {
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

	if policy.AgentUid != "" {
		var agent models.Agent
		fireAgent := lib.GetDatasetByEnv(origin, models.AgentCollection)
		docsnap, err := lib.GetFirestoreErr(fireAgent, policy.AgentUid)
		lib.CheckError(err)
		err = docsnap.DataTo(&agent)
		lib.CheckError(err)
		policyProducer["name"] = strings.ToUpper(agent.Surname) + " " + strings.ToUpper(agent.Name)
		policyProducer["ruiSection"] = agent.RuiSection
		policyProducer["ruiCode"] = agent.RuiCode
		policyProducer["ruiRegistration"] = agent.RuiRegistration.Format("02.01.2006")
	} else if policy.AgencyUid != "" {
		var agency models.Agency
		fireAgency := lib.GetDatasetByEnv(origin, models.AgencyCollection)
		docsnap, err := lib.GetFirestoreErr(fireAgency, policy.AgencyUid)
		lib.CheckError(err)
		err = docsnap.DataTo(&agency)
		lib.CheckError(err)
		policyProducer["name"] = strings.ToUpper(agency.Name)
		policyProducer["ruiSection"] = agency.RuiSection
		policyProducer["ruiCode"] = agency.RuiCode
		policyProducer["ruiRegistration"] = agency.RuiRegistration.Format("02.01.2006")
	}
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
				if guarantee.Beneficiaries == nil || (*guarantee.Beneficiaries)[0].IsLegitimateSuccessors {
					details += "Eredi leggitimi e/o testamentari"
				} else {
					for _, beneficiary := range *guarantee.Beneficiaries {
						details += beneficiary.Name + " " + beneficiary.Surname + " " + beneficiary.FiscalCode + "\n"
					}
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
