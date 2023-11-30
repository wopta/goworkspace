package document

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type slugStruct struct {
	name  string
	order int64
}

func loadLifeBeneficiariesInfo(policy *models.Policy) ([]map[string]string, string, string) {
	legitimateSuccessorsChoice := ""
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
	if err != nil {
		return beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice
	}

	if (*deathGuarantee.Beneficiaries)[0].BeneficiaryType != "legalAndWillSuccessors" {
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
	} else {
		legitimateSuccessorsChoice = "X"
	}

	return beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice
}

func loadProponentInfo(networkNode *models.NetworkNode) map[string]string {
	policyProponent := make(map[string]string)

	if networkNode == nil || networkNode.IsMgaProponent {
		policyProponent["address"] = "Galleria del Corso, 1 - 20122 MILANO (MI)"
		policyProponent["phone"] = "02.91.24.03.46"
		policyProponent["email"] = "info@wopta.it"
		policyProponent["pec"] = "woptaassicurazioni@legalmail.it"
		policyProponent["website"] = "wopta.it"
	} else {
		proponentNode := network.GetNetworkNodeByUid(networkNode.WorksForUid)
		if proponentNode == nil {
			panic("could not find node for proponent with uid " + networkNode.WorksForUid)
		}

		policyProponent["address"] = proponentNode.GetAddress()
		policyProponent["phone"] = proponentNode.Agency.Phone
		policyProponent["email"] = proponentNode.Mail
		policyProponent["pec"] = proponentNode.Agency.Pec
		policyProponent["website"] = proponentNode.Agency.Website
	}

	jsonProponent, _ := json.Marshal(policyProponent)
	log.Printf("[loadProponentInfo] proponent info %s", string(jsonProponent))
	return policyProponent
}

func loadProducerInfo(origin string, networkNode *models.NetworkNode) map[string]string {
	policyProducer := map[string]string{
		"name":            "LOMAZZI MICHELE",
		"ruiSection":      "A",
		"ruiCode":         "A000703480",
		"ruiRegistration": "02.03.2022",
	}

	if networkNode == nil || strings.EqualFold(networkNode.Type, models.PartnershipNetworkNodeType) {
		jsonProducer, _ := json.Marshal(policyProducer)
		log.Printf("[loadProducerInfo] producer info %s", string(jsonProducer))
		return policyProducer
	}

	log.Printf("[loadProducerInfo] setting producer %s info, producerType %s", networkNode.Uid, networkNode.Type)

	switch networkNode.Type {
	case models.AgentNetworkNodeType:
		policyProducer["name"] = strings.ToUpper(networkNode.Agent.Surname) + " " + strings.ToUpper(networkNode.Agent.Name)
		policyProducer["ruiSection"] = networkNode.Agent.RuiSection
		policyProducer["ruiCode"] = networkNode.Agent.RuiCode
		policyProducer["ruiRegistration"] = networkNode.Agent.RuiRegistration.Format("02.01.2006")
	case models.AgencyNetworkNodeType:
		policyProducer["name"] = strings.ToUpper(networkNode.Agency.Name)
		policyProducer["ruiSection"] = networkNode.Agency.RuiSection
		policyProducer["ruiCode"] = networkNode.Agency.RuiCode
		policyProducer["ruiRegistration"] = networkNode.Agency.RuiRegistration.Format("02.01.2006")
	}

	jsonProducer, _ := json.Marshal(policyProducer)
	log.Printf("[loadProducerInfo] producer info %s", string(jsonProducer))

	return policyProducer
}

func loadDesignation(networkNode *models.NetworkNode) string {
	var (
		designation                           string
		mgaProponentDirectDesignationFormat   = "%s %s"
		mgaRuiInfo                            = "Wopta Assicurazioni Srl, Società iscritta alla Sezione A del RUI con numero A000701923 in data 14/02/2022"
		designationDirectManager              = "Responsabile dell’attività di intermediazione assicurativa di"
		mgaProponentIndirectDesignationFormat = "%s di %s, iscritta in sezione E del RUI con numero %s in data %s, che opera per conto di %s"
		mgaEmitterDesignationFormat           = "%s dell’intermediario di %s iscritta alla sezione %s del RUI con numero %s in data %s"
	)

	if networkNode == nil || networkNode.Type == models.PartnershipNetworkNodeType {
		designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, designationDirectManager, mgaRuiInfo)
	} else if networkNode.IsMgaProponent {
		if networkNode.WorksForUid == models.WorksForMgaUid {
			designation = fmt.Sprintf(mgaProponentDirectDesignationFormat, networkNode.Designation, mgaRuiInfo)
		} else {
			worksForNode := network.GetNetworkNodeByUid(networkNode.WorksForUid)
			designation = fmt.Sprintf(
				mgaProponentIndirectDesignationFormat,
				networkNode.Designation,
				worksForNode.Agency.Name,
				worksForNode.Agency.RuiCode,
				worksForNode.Agency.RuiRegistration.Format(dateLayout),
				mgaRuiInfo,
			)
		}
	} else {
		worksForNode := network.GetNetworkNodeByUid(networkNode.WorksForUid)
		designation = fmt.Sprintf(
			mgaEmitterDesignationFormat,
			networkNode.Designation,
			worksForNode.Agency.Name,
			worksForNode.Agency.RuiSection,
			worksForNode.Agency.RuiCode,
			worksForNode.Agency.RuiRegistration.Format(dateLayout),
		)
	}

	log.Printf("[loadDesignation] designation info %s", designation)

	return designation
}

func loadLifeGuarantees(policy *models.Policy, product *models.Product) (map[string]map[string]string, []slugStruct) {
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

	guaranteesMap = make(map[string]map[string]string, 0)

	for guaranteeSlug, guarantee := range product.Companies[0].GuaranteesMap {
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

func loadPersonaGuarantees(policy *models.Policy, product *models.Product) (map[string]map[string]string, []slugStruct) {
	var (
		guaranteesMap map[string]map[string]string
		slugs         []slugStruct
	)

	guaranteesMap = make(map[string]map[string]string, 0)
	offerName := policy.OfferlName

	for guaranteeSlug, guarantee := range product.Companies[0].GuaranteesMap {
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
						details += product.Companies[0].GuaranteesMap["D"].BeneficiaryOptions[beneficiary.
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

func loadAnnex4Section1Info(policy *models.Policy, networkNode *models.NetworkNode) string {
	var (
		section1Info       string
		mgaProponentFormat = "Secondo quanto indicato nel modulo di proposta/polizza e documentazione " +
			"precontrattuale ricevuta, la distribuzione  relativamente a questa proposta/contratto è svolta per " +
			"conto della seguente impresa di assicurazione: %s"
		mgaEmitterFormat = "Il contratto viene intermediato da %s, in qualità di soggetto proponente, che opera in " +
			"virtù della collaborazione con Wopta Assicurazioni Srl (intermediario emittente dell'Impresa di " +
			"Assicurazione %s, iscritto al RUI sezione A nr A000701923 dal 14.02.2022, ai sensi dell’articolo 22, " +
			"comma 10, del decreto legge 18 ottobre 2012, n. 179, convertito nella legge 17 dicembre 2012, n. 221"
	)

	if policy.Channel != models.NetworkChannel || networkNode == nil || networkNode.IsMgaProponent {
		section1Info = fmt.Sprintf(
			mgaProponentFormat,
			companyMap[policy.Company],
		)
	} else {
		worksForNode := network.GetNetworkNodeByUid(networkNode.WorksForUid)
		section1Info = fmt.Sprintf(
			mgaEmitterFormat,
			worksForNode.Agency.Name,
			companyMap[policy.Company],
		)
	}

	log.Printf("[loadAnnex4Section1Info] section 1 info: %s", section1Info)

	return section1Info
}
