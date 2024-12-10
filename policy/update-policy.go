package policy

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/models"
)

func UpdatePolicy(policy *models.Policy) (map[string]interface{}, error) {
	input := make(map[string]interface{}, 0)

	input["assets"] = policy.Assets
	input["contractor"] = policy.Contractor
	input["fundsOrigin"] = policy.FundsOrigin
	if policy.Surveys != nil {
		input["surveys"] = policy.Surveys
	}
	if policy.Statements != nil {
		input["statements"] = policy.Statements
	}
	input["step"] = policy.Step
	if policy.OfferlName != "" {
		input["offerName"] = policy.OfferlName
	}

	switch policy.Name {
	case models.PersonaProduct:
		input["taxAmount"] = policy.TaxAmount
		input["priceNett"] = policy.PriceNett
		input["priceGross"] = policy.PriceGross
		input["taxAmountMonthly"] = policy.TaxAmountMonthly
		input["priceNettMonthly"] = policy.PriceNettMonthly
		input["priceGrossMonthly"] = policy.PriceGrossMonthly
	case models.CommercialCombinedProduct:
		err := checkQbe(policy, input)
		if err != nil {
			return nil, err
		}
	}

	input["updated"] = time.Now().UTC()

	return input, nil
}

func checkQbe(p *models.Policy, i map[string]interface{}) error {
	var err error

	err = checkQbeAssets(p)
	if err != nil {
		return err
	}
	err = checkQbeContractor(p)
	if err != nil {
		return err
	}
	err = checkDeclaredClaims(p.DeclaredClaims)
	if err != nil {
		return err
	}
	err = checkQbeBonds(p)
	if err != nil {
		return err
	}
	err = checkQbeStatements(p)
	if err != nil {
		return err
	}
	err = checkQbeSignatory(p)
	if err != nil {
		return err
	}

	i["startDate"] = p.StartDate
	i["endDate"] = p.EndDate
	i["declaredClaims"] = p.DeclaredClaims
	i["hasBond"] = p.HasBond
	i["bond"] = p.Bond
	i["clause"] = p.Clause

	return nil
}

func checkQbeSignatory(p *models.Policy) error {
	if p.Contractors == nil || len(*p.Contractors) == 0 {
		return fmt.Errorf("no signatory provided")
	}
	sgn := (*p.Contractors)[0]
	if sgn.Name == "" {
		return fmt.Errorf("empty signatory name")
	}
	if sgn.Surname == "" {
		return fmt.Errorf("empty signatory surname")
	}
	if sgn.FiscalCode == "" {
		return fmt.Errorf("empty signatory fiscal code")
	}
	if checkFiscalCode(sgn.FiscalCode) == false {
		return fmt.Errorf("wrong signatory fiscal code")
	}
	if sgn.Phone == "" {
		return fmt.Errorf("wrong signatory phone number")
	}
	if sgn.Mail == "" {
		return fmt.Errorf("wrong signatory mail")
	}
	if sgn.Consens == nil || len(*sgn.Consens) == 0 {
		return fmt.Errorf("empty signatory consens")
	}

	return nil
}

func checkQbeStatements(p *models.Policy) error {
	if len(*p.Statements) != 5 {
		return fmt.Errorf("there must be 5 signed statements")
	}

	return nil
}

func checkQbeBonds(p *models.Policy) error {
	if (p.HasBond) && (p.Bond == "") {
		return fmt.Errorf("empty bond")
	}

	return nil
}

func checkQbeAssets(p *models.Policy) error {
	var err error

	for _, v := range p.Assets {
		if v.Uuid == "" {
			return fmt.Errorf("empty Uuid for asset %s", v.Name)
		}
		if v.Enterprise != nil {
			err = checkEnterprise(v.Enterprise)
			if err != nil {
				return err
			}
		}
		if v.Building != nil {
			err = checkBuilding(v.Building)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkEnterprise(e *models.Enterprise) error {
	if e.Employer < 1 || e.Employer > 100 {
		return fmt.Errorf("number of employers must be between 1 and 100")
	}
	if e.WorkEmployersRemuneration == "" {
		return fmt.Errorf("empty employer remuneration")
	}
	num, err := strconv.Atoi(e.WorkEmployersRemuneration)
	if err != nil {
		return fmt.Errorf("employer remuneration must consist of digits only")
	}
	if num < 1 {
		return fmt.Errorf("employer remuneration must be greater than 1")
	}
	if e.Revenue == "" {
		return fmt.Errorf("empty total revenue")
	}

	return nil
}

func checkBuilding(b *models.Building) error {
	if b.NaicsCategory == "" {
		return fmt.Errorf("empty naics category")
	}
	if b.NaicsDetail == "" {
		return fmt.Errorf("empty naics detail")
	}
	if b.Naics == "" {
		return fmt.Errorf("empty naics code")
	}
	if b.BuildingMaterial == "" {
		return fmt.Errorf("empty building material")
	}
	err := checkAddress(b.BuildingAddress)
	if err != nil {
		return err
	}

	return nil
}

func checkQbeContractor(p *models.Policy) error {
	if p.Contractor.Type != "legalEntity" {
		return fmt.Errorf("invalid contractor type")
	}
	if _, err := strconv.Atoi(p.Contractor.VatCode); err != nil {
		return fmt.Errorf("contractor Vat code must consist of digits only")
	}
	if len(p.Contractor.VatCode) != 11 {
		return fmt.Errorf("contractor Vat code must have 11 digits")
	}
	if checkFiscalCode(p.Contractor.FiscalCode) == false {
		return fmt.Errorf("wrong contractor fiscal code")
	}
	if p.Contractor.Name == "" {
		return fmt.Errorf("empty contractor name")
	}
	err := checkAddress(p.Contractor.CompanyAddress)
	if err != nil {
		return err
	}

	return nil
}

func checkAddress(a *models.Address) error {
	if a.StreetName == "" {
		return fmt.Errorf("empty address: street name")
	}
	if a.StreetNumber == "" {
		return fmt.Errorf("empty address: street number")
	}
	if a.Locality == "" {
		return fmt.Errorf("empty address: locality")
	}
	if a.City == "" {
		return fmt.Errorf("empty address: city")
	}
	if a.PostalCode == "" {
		return fmt.Errorf("empty address: postal code")
	}
	if a.CityCode == "" {
		return fmt.Errorf("empty address: city code")
	}

	return nil
}

func checkFiscalCode(fc string) bool {
	return true
}

func checkDeclaredClaims(d []models.DeclaredClaims) error {
	for _, v := range d {
		if v.GuaranteeSlug == "" {
			return fmt.Errorf("empty guarantee slug")
		}
		if len(v.History) == 0 {
			return fmt.Errorf("guarantee history must contain at least one year")
		}
	}

	return nil
}
