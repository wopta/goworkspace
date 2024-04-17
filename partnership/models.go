package partnership

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wopta/goworkspace/models"
)

type PartnershipNode struct {
	Name string       `json:"name"`
	Skin *models.Skin `json:"skin,omitempty"`
}

type PartnershipResponse struct {
	Policy      models.Policy   `json:"policy"`
	Partnership PartnershipNode `json:"partnership"`
	Product     models.Product  `json:"product"`
}

type BeprofClaims struct {
	UserBeprofid         int    `json:"user.beprofid"`
	UserFirstname        string `json:"user.firstname"`
	UserLastname         string `json:"user.lastname"`
	UserEmail            string `json:"user.email"`
	UserMobile           string `json:"user.mobile"`
	UserFiscalcode       string `json:"user.fiscalcode"`
	UserPiva             string `json:"user.piva"`
	UserProvince         string `json:"user.province"`
	UserCity             string `json:"user.city"`
	UserPostalcode       string `json:"user.postalcode"`
	UserAddress          string `json:"user.address"`
	UserMunicipalityCode string `json:"user.municipality_code"`
	UserEmploymentSector string `json:"user.employment_sector"`
	ProductCode          string `json:"product.code"`
	ProductPurchaseid    string `json:"product.purchaseid"`
	Price                string `json:"price"`
	jwt.RegisteredClaims
}

type BeprofLifeClaimsAdapter struct {
	beprofClaims *BeprofClaims
}

func (a *BeprofLifeClaimsAdapter) ExtractClaims() (models.LifeClaims, error) {
	data := make(map[string]interface{})
	b, err := json.Marshal(a.beprofClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return models.LifeClaims{}, err
	}

	return models.LifeClaims{
		Name:       a.beprofClaims.UserFirstname,
		Surname:    a.beprofClaims.UserLastname,
		Email:      a.beprofClaims.UserEmail,
		FiscalCode: a.beprofClaims.UserFiscalcode,
		Address:    a.beprofClaims.UserAddress,
		Postalcode: a.beprofClaims.UserPostalcode,
		City:       a.beprofClaims.UserCity,
		CityCode:   a.beprofClaims.UserMunicipalityCode,
		Work:       a.beprofClaims.UserEmploymentSector,
		VatCode:    a.beprofClaims.UserPiva,
		Data:       data,
	}, nil
}

type FacileClaims struct {
	CustomerName       string `json:"customerName"`
	CustomerFamilyName string `json:"customerFamilyName"`
	CustomerBirthDate  string `json:"customerBirthDate"`
	Gender             string `json:"gender"`
	Email              string `json:"email"`
	Mobile             string `json:"mobile"`
	IsSmoker           bool   `json:"isSmoker"`
	InsuredCapital     int    `json:"insuredCapital"`
	Duration           int    `json:"duration"`
	jwt.RegisteredClaims
}

type FacileLifeClaimsAdapter struct {
	facileClaims *FacileClaims
}

func (a *FacileLifeClaimsAdapter) ExtractClaims() (models.LifeClaims, error) {
	data := make(map[string]interface{})
	b, err := json.Marshal(a.facileClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return models.LifeClaims{}, err
	}

	birthDate, _ := time.Parse(models.TimeDateOnly, a.facileClaims.CustomerBirthDate)

	return models.LifeClaims{
		Name:      a.facileClaims.CustomerName,
		Surname:   a.facileClaims.CustomerFamilyName,
		Email:     a.facileClaims.Email,
		BirthDate: birthDate.Format(time.RFC3339),
		Phone:     fmt.Sprintf("+39%s", a.facileClaims.Mobile),
		Gender:    a.facileClaims.Gender,
		Guarantees: map[string]models.ClaimsGuarantee{
			"death": {
				Duration:                   a.facileClaims.Duration,
				SumInsuredLimitOfIndemnity: float64(a.facileClaims.InsuredCapital),
			},
		},
		Data: data,
	}, nil
}

type ELeadsClaims struct {
	Name                       string `json:"name"`
	Surname                    string `json:"surname"`
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	FiscalCode                 string `json:"fiscalCode"`
	BirthDate                  string `json:"birthDate"`
	SumInsuredLimitOfIndemnity int    `json:"sumInsuredLimitOfIndemnity"`
	Duration                   int    `json:"duration"`
	jwt.RegisteredClaims
}

type ELeadsLifeClaimsAdapter struct {
	eLeadsClaims *ELeadsClaims
}

func (a *ELeadsLifeClaimsAdapter) ExtractClaims() (models.LifeClaims, error) {
	data := make(map[string]interface{})
	b, err := json.Marshal(a.eLeadsClaims)
	if err != nil {
		return models.LifeClaims{}, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return models.LifeClaims{}, err
	}

	birthDate, _ := time.Parse(models.TimeDateOnly, a.eLeadsClaims.BirthDate)

	return models.LifeClaims{
		Name:       a.eLeadsClaims.Name,
		Surname:    a.eLeadsClaims.Surname,
		Email:      a.eLeadsClaims.Email,
		BirthDate:  birthDate.Format(time.RFC3339),
		Phone:      fmt.Sprintf("+39%s", a.eLeadsClaims.Phone),
		FiscalCode: a.eLeadsClaims.FiscalCode,
		Guarantees: map[string]models.ClaimsGuarantee{
			"death": {
				Duration:                   a.eLeadsClaims.Duration,
				SumInsuredLimitOfIndemnity: float64(a.eLeadsClaims.SumInsuredLimitOfIndemnity),
			},
		},
		Data: data,
	}, nil
}
