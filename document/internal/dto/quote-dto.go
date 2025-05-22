package dto

import (
	"slices"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type QuoteBaseDTO struct {
	Price     priceDTO
	StartDate time.Time
}

type QuoteLifeDTO struct {
	*QuoteBaseDTO
	Contractor contractorDTO
	Guarantees map[string]*quoteGuaranteeDTO
}

func NewQuoteLifeDTO() *QuoteLifeDTO {
	return &QuoteLifeDTO{
		QuoteBaseDTO: &QuoteBaseDTO{
			StartDate: time.Time{},
			Price:     *newPriceDTO(),
		},
		Contractor: *newContractorDTO(),
		Guarantees: make(map[string]*quoteGuaranteeDTO),
	}
}
func (q *QuoteBaseDTO) fromData(p *models.Policy) {
	q.StartDate = p.StartDate

	q.Price.Gross.FromValue(p.PriceGross)
	q.Price.Consultancy.FromValue(p.ConsultancyValue.Price)
	q.Price.Total.FromValue(p.PriceGross + p.ConsultancyValue.Price)
}

func (q *QuoteLifeDTO) fromData(p *models.Policy, prd *models.Product) {
	birthDate := "INVALID DATE"
	if rawBirthDate, err := time.Parse(time.RFC3339, p.Contractor.BirthDate); err == nil {
		birthDate = rawBirthDate.Format(constants.DayMonthYearFormat)
	}
	q.Contractor.BirthDate = birthDate

	gMap := p.GuaranteesToMap()
	companyIdx := slices.IndexFunc(prd.Companies, func(c models.Company) bool {
		return c.Name == p.Company
	})
	if companyIdx < 0 {
		return
	}
	for slug, productGuarantee := range prd.Companies[companyIdx].GuaranteesMap {
		dto := newQuoteGuaranteeDTO()
		dto.Description = productGuarantee.CompanyName
		if g, ok := gMap[slug]; ok {
			dto.fromData(g, q.StartDate)
		}
		q.Guarantees[slug] = dto
	}
}

func (q *QuoteLifeDTO) FromData(p *models.Policy, prd *models.Product) {
	q.QuoteBaseDTO.fromData(p)
	q.fromData(p, prd)
}
