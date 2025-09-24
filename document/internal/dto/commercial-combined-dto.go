package dto

import (
	"slices"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type CommercialCombinedDTO struct {
	Contract           *contractDTO
	Contractor         *contractorDTO
	Enterprise         *enterpriseDTO
	Buildings          []*buildingDTO
	Claims             map[string]*claimDTO
	Prices             *priceDTO
	PricesBySection    map[string]*section
	HasExcludedFormula bool
}

func NewCommercialCombinedDto() *CommercialCombinedDTO {
	return &CommercialCombinedDTO{
		Contract:        newContractDTO(),
		Contractor:      newContractorDTO(),
		Enterprise:      newEnterpriseDTO(),
		Buildings:       make([]*buildingDTO, 0),
		Claims:          make(map[string]*claimDTO),
		Prices:          newPriceDTO(),
		PricesBySection: make(map[string]*section),
	}
}

func (cc *CommercialCombinedDTO) FromPolicy(policy models.Policy, product models.Product, isProposal bool) {
	var numBuildings int64

	cc.Contract.fromPolicy(policy, isProposal)
	cc.Contractor.fromPolicy(policy.Contractor)

	productGuarantees := product.Companies[0].GuaranteesMap

	cc.Buildings = make([]*buildingDTO, 0, 5)
	for i := 0; i < 5; i++ {
		building := newBuildingDTO()
		for _, guarantee := range productGuarantees {
			if guarantee.Type == "building" {
				newGuarantee := newGuaranteeDTO()
				newGuarantee.Description = guarantee.CompanyName
				building.Guarantees[guarantee.Slug] = newGuarantee
			}
		}
		cc.Buildings = append(cc.Buildings, building)
	}

	cc.Enterprise = newEnterpriseDTO()
	for _, guarantee := range productGuarantees {
		//TO CHANGE, THIS IS SHIT
		if guarantee.Type == models.UserLegalEntity {
			newGuarantee := newGuaranteeDTO()
			newGuarantee.Description = guarantee.CompanyName
			if guarantee.Slug == "additional-compensation" {
				newGuarantee.Description = "Danni Indiretti - Formula"
			}
			cc.Enterprise.Guarantees[guarantee.Slug] = newGuarantee
		}
	}

	for index, asset := range policy.Assets {
		if asset.Building != nil {
			cc.Buildings[numBuildings].fromPolicy(*policy.Assets[index].Building, policy.Assets[index].Guarantees)
			numBuildings++
		}
		if asset.Enterprise != nil {
			cc.Enterprise.fromPolicy(*policy.Assets[index].Enterprise, policy.Assets[index].Guarantees)
			isExcluded := true
			for _, guarantee := range policy.Assets[index].Guarantees {
				if slices.Contains([]string{"daily-allowance", "increased-cost", "additional-compensation"}, guarantee.Slug) {
					isExcluded = false
					break
				}
			}
			if isExcluded {
				cc.HasExcludedFormula = true
			}
		}
	}

	claimsDescriptionsMap := map[string]string{
		"property":                "Danni ai beni (escluso Furto)",
		"third-party-liability":   "Responsabilità Civile",
		"theft":                   "Furto",
		"management-organization": "D&O",
		"cyber":                   "Cyber",
	}

	for slug, description := range claimsDescriptionsMap {
		cc.Claims[slug] = newClaimDTO()
		cc.Claims[slug].Description = description
	}

	for _, declaredClaim := range policy.DeclaredClaims {
		slug := declaredClaim.GuaranteeSlug
		if _, ok := cc.Claims[slug]; !ok {
			continue
		}
		for _, history := range declaredClaim.History {
			cc.Claims[slug].Quantity.FromValue(cc.Claims[slug].Quantity.ValueInt + int64(history.Quantity))
			cc.Claims[slug].Value.FromValue(cc.Claims[slug].Value.ValueFloat + history.Value)
		}
	}

	cc.Prices.Gross.FromValue(policy.PriceGross)
	cc.Prices.Net.FromValue(policy.PriceNett)
	cc.Prices.Taxes.FromValue(policy.TaxAmount)

	sectionMap := map[string]string{
		"A": "A - INCENDIO E \"TUTTI I RISCHI\"",
		"B": "B - DANNI INDIRETTI",
		"C": "C - FURTO",
		"D": "D - RESPONSABILITÀ CIVILE VERSO TERZI (RCT)",
		"E": "E - RESP. CIVILE VERSO PRESTATORI DI LAVORO (RCO)",
		"F": "F - RESP. CIVILE DA PRODOTTI DIFETTOSI (RCP)",
		"G": "G - RITIRO PRODOTTI",
		"H": "H - RESP. AMMINISTRATORI SINDACI DIRIGENTI (D&O)",
		"I": "I - CYBER RESPONSE E DATA SECURITY",
	}

	for sectionKey, description := range sectionMap {
		cc.PricesBySection[sectionKey] = newSection()
		cc.PricesBySection[sectionKey].Description = description
	}

	groupSectionMap := map[string]string{
		"Fabbricato":                                   "A",
		"Contenuto (Merci e Macchinari)":               "A",
		"Merci (aumento temporaneo)":                   "A",
		"Furto, rapina, estorsione (in aumento)":       "C",
		"Rischio locativo (in aumento)":                "A",
		"Altre garanzie su Contenuto":                  "A",
		"Ricorso terzi (in aumento)":                   "A",
		"Danni indiretti":                              "B",
		"Perdita Pigioni":                              "B",
		"Responsabilità civile terzi":                  "D",
		"Responsabilità civile prestatori lavoro":      "E",
		"Responsabilità civile prodotti":               "F",
		"Ritiro Prodotti":                              "G",
		"Resp. Amministratori Sindaci Dirigenti (D&O)": "H",
		"Cyber": "I",
	}

	for _, price := range policy.PriceGroup {
		sectionKey := groupSectionMap[price.Name]
		cc.PricesBySection[sectionKey].Price.Gross.FromValue(cc.PricesBySection[sectionKey].Price.Gross.ValueFloat + price.Gross)
		cc.PricesBySection[sectionKey].Price.Net.FromValue(cc.PricesBySection[sectionKey].Price.Net.ValueFloat + price.Net)
		cc.PricesBySection[sectionKey].Price.Taxes.FromValue(cc.PricesBySection[sectionKey].Price.Taxes.ValueFloat + price.Tax)
		if cc.PricesBySection[sectionKey].Active == constants.No && cc.PricesBySection[sectionKey].Price.Gross.ValueFloat > 0 {
			cc.PricesBySection[sectionKey].Active = constants.Yes
		}
	}
}

type section struct {
	Description string
	Active      string
	Price       *priceDTO
}

func newSection() *section {
	return &section{
		Description: constants.EmptyField,
		Active:      constants.No,
		Price:       newPriceDTO(),
	}
}
