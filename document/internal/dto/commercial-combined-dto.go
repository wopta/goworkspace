package dto

import (
	"slices"
	"strconv"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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
		if guarantee.Type == "enterprise" {
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
			cc.Claims[slug].Quantity.ValueInt += int64(history.Quantity)
			cc.Claims[slug].Quantity.Text = strconv.FormatInt(cc.Claims[slug].Quantity.ValueInt, 10)
			cc.Claims[slug].Value.ValueFloat += history.Value
			cc.Claims[slug].Value.Text = lib.HumanaizePriceEuro(cc.Claims[slug].Value.ValueFloat)
		}
	}

	cc.Prices.Gross = policy.PriceGross
	cc.Prices.GrossText = lib.HumanaizePriceEuro(cc.Prices.Gross)
	cc.Prices.Net = policy.PriceNett
	cc.Prices.NetText = lib.HumanaizePriceEuro(cc.Prices.Net)
	cc.Prices.Taxes = policy.TaxAmount
	cc.Prices.TaxesText = lib.HumanaizePriceEuro(cc.Prices.Taxes)

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
		cc.PricesBySection[sectionKey].Price.Gross += price.Gross
		cc.PricesBySection[sectionKey].Price.GrossText = lib.HumanaizePriceEuro(cc.PricesBySection[sectionKey].Price.Gross)
		cc.PricesBySection[sectionKey].Price.Net += price.Net
		cc.PricesBySection[sectionKey].Price.NetText = lib.HumanaizePriceEuro(cc.PricesBySection[sectionKey].Price.Net)
		cc.PricesBySection[sectionKey].Price.Taxes += price.Tax
		cc.PricesBySection[sectionKey].Price.TaxesText = lib.HumanaizePriceEuro(cc.PricesBySection[sectionKey].Price.Taxes)
		if cc.PricesBySection[sectionKey].Active == constants.No && cc.PricesBySection[sectionKey].Price.Gross > 0 {
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
