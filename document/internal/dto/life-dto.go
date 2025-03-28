package dto

import (
	"fmt"
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"time"
)

var splitPayment map[string]string = map[string]string{
	string(models.PaySplitMonthly):           "mensile",
	string(models.PaySplitYearly):            "annuale",
	string(models.PaySplitSingleInstallment): "singolo",
}

type LifeDTO struct {
	Contractor       contractorDTO
	Channel          string
	Prizes           priceDTO
	PriceAnnuity     string
	ConsultancyValue consultancyDTO
	ValidityDate     validityDateDTO
	ProductorName    string
	ProposalNumber   string
}

type validityDateDTO struct {
	StartDate          string
	EndDate            string
	FirstAnnuityExpiry string
}

func formatDate(t time.Time) string {
	location, _ := time.LoadLocation("Europe/Rome")
	time := t.In(location)
	return time.In(location).Format(constants.DayMonthYearFormat) + " ore 24:00"
}

func getSplit(split string) string {
	if split, ok := splitPayment[split]; ok {
		return split
	}
	return ""
}

func NewLifeDto() LifeDTO {
	return LifeDTO{}
}

func (n *LifeDTO) FromPolicy(policy *models.Policy, network *models.NetworkNode) {
	n.Channel = policy.Channel
	n.ProposalNumber = fmt.Sprint(policy.ProposalNumber)
	(&n.Contractor).fromPolicy(policy.Contractor)

	n.Prizes = priceDTO{
		Split: getSplit(policy.PaymentSplit),
	}
	n.Prizes.Gross.ValueFloat = policy.PriceGross
	n.Prizes.Gross.Text = lib.HumanaizePriceEuro(policy.PriceGross)
	n.ConsultancyValue.Price = policy.ConsultancyValue.Price
	n.PriceAnnuity = lib.HumanaizePriceEuro(policy.PaymentComponents.PriceAnnuity.Total)

	n.ValidityDate = validityDateDTO{
		StartDate:          formatDate(policy.StartDate),
		EndDate:            formatDate(policy.EndDate),
		FirstAnnuityExpiry: formatDate(policy.StartDate.AddDate(1, 0, 0)),
	}

	if policy.Channel == models.ECommerceChannel {
		n.ProductorName = "Michele Lomazzi"
	}
	if policy.Channel == models.NetworkChannel {
		n.ProductorName = network.GetName()
	}
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
