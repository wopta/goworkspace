package dto

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

var splitPayment map[string]string = map[string]string{
	string(models.PaySplitMonthly):           "mensile",
	string(models.PaySplitYearly):            "annuale",
	string(models.PaySplitSingleInstallment): "singolo",
}

type LifeDTO struct {
	Contractor       *contractorDTO
	Channel          string
	Prizes           *priceDTO
	PriceAnnuity     string
	ConsultancyValue *consultancyDTO
	ValidityDate     *validityDateDTO
	ProductorName    string
	ProposalNumber   string
}

func NewLifeDto() LifeDTO {
	return LifeDTO{
		Contractor:       newContractorDTO(),
		Channel:          constants.EmptyField,
		Prizes:           newPriceDTO(),
		PriceAnnuity:     constants.EmptyField,
		ConsultancyValue: newConsultacyDTO(),
		ValidityDate:     newValidityDateDTO(),
		ProductorName:    constants.EmptyField,
		ProposalNumber:   constants.EmptyField,
	}
}

func (n *LifeDTO) FromPolicy(policy *models.Policy, network *models.NetworkNode) {
	n.Channel = policy.Channel
	n.ProposalNumber = fmt.Sprint(policy.ProposalNumber)

	n.Contractor.fromPolicy(policy.Contractor)

	n.Prizes.Split = getSplit(policy.PaymentSplit)
	n.Prizes.Gross.ValueFloat = policy.PriceGross
	n.Prizes.Gross.Text = lib.HumanaizePriceEuro(policy.PriceGross)

	n.ConsultancyValue.Price.FromValue(policy.ConsultancyValue.Price)

	n.PriceAnnuity = lib.HumanaizePriceEuro(policy.PaymentComponents.PriceAnnuity.Total)
	n.ValidityDate.fromPolicy(policy)

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
