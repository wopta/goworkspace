package dto

import (
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

type LifeDTO struct {
	Contractor       contractorDTO
	Channel          string
	Prizes           prizeDTO
	PriceAnnuity     string
	ConsultancyValue consultancyDto
	ValidityDate     validitydate
	ProductorName    string
}

type validitydate struct {
	StartDate          string
	EndDate            string
	FirstAnnuityExpiry string
}

type prizeDTO struct {
	Gross float64
	Split string
}

func formatDate(t time.Time) string {
	location, _ := time.LoadLocation("Europe/Rome")
	time := t.In(location)
	return time.In(location).Format("02/01/2006")
}

func getSplit(split string) string {
	switch split {
	case string(models.PaySplitMonthly):
		return "mensile"
	case string(models.PaySplitYearly):
		return "annuale"
	case string(models.PaySplitSingleInstallment):
		return "singolo"
	}
	return ""
}

func NewLifeDto(policy *models.Policy) LifeDTO {
	dto := LifeDTO{}
	dto.Channel = policy.Channel
	(&dto.Contractor).fromPolicy(policy.Contractor)
	dto.Prizes = prizeDTO{
		Split: getSplit(policy.PaymentSplit),
		Gross: policy.PriceGross,
	}
	dto.ConsultancyValue.Price = lib.HumanaizePriceEuro(policy.ConsultancyValue.Price)
	dto.PriceAnnuity = lib.HumanaizePriceEuro(policy.PaymentComponents.PriceAnnuity.Total)

	dto.ValidityDate = validitydate{
		StartDate:          formatDate(policy.StartDate),
		EndDate:            formatDate(policy.EndDate),
		FirstAnnuityExpiry: formatDate(policy.StartDate.AddDate(1, 0, 0)),
	}

	networkModel := network.GetNetworkNodeByUid(policy.ProducerUid)
	if policy.Channel == models.ECommerceChannel {
		dto.ProductorName = "Michele Lomazzi"
	}
	if policy.Channel == models.NetworkChannel {
		dto.ProductorName = networkModel.GetName()
	}
	return dto
}

func (l *LifeDTO) GetAddressFirstPart() string {
	return l.Contractor.StreetName + ", " + l.Contractor.StreetNumber
}
func (l *LifeDTO) GetAddressSecondPart() string {
	return l.Contractor.PostalCode + " " + l.Contractor.City + " (" + l.Contractor.CityCode + ")"
}
func (l *LifeDTO) GetFullNameContractor() string {
	return l.Contractor.Name + " " + l.Contractor.Surname
}
