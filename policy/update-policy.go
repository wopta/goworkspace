package policy

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/wopta/goworkspace/lib"
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
		// NOTE: remove qbe check for now.

		//err := checkQbe(policy, input)
		//if err != nil {
		//	return nil, err
		//}
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
	i["contractors"] = p.Contractors
	i["priceGroup"] = p.PriceGroup

	return nil
}

func checkQbeSignatory(p *models.Policy) error {
	if p.Contractors == nil || len(*p.Contractors) == 0 {
		return nil
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

	if p.Assets == nil {
		return nil
	}

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
	if e.WorkEmployersRemuneration < 1 {
		return fmt.Errorf("employer remuneration must be greater than 1")
	}
	if e.Revenue < 1 {
		return fmt.Errorf("revenue must be greater than 0")
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
	if p.Contractor.Name == "" {
		return nil
	}
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
	err := checkAddress(p.Contractor.CompanyAddress)
	if err != nil {
		return err
	}

	return nil
}

func checkAddress(a *models.Address) error {
	if a == nil {
		return fmt.Errorf("nil address")
	}
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
	if a.IsManualInput {
		err := verifyManualAddress(a.City, a.PostalCode, a.CityCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func verifyManualAddress(city, postalCode, cityCode string) error {
	fileName := "enrich/postal-codes.csv"
	file, err := lib.GetFilesByEnvV2(fileName)
	if err != nil {
		log.Printf("error reading file %s: %v", fileName, err)
		return err
	}
	df, err := lib.CsvToDataframeV2(file, ';', true)
	if err != nil {
		log.Printf("error reading df %v", err)
		return err
	}

	columns := []string{"postal code", "place name", "admin code2"}
	sel := df.Select(columns)

	fil := sel.Filter(dataframe.F{Colname: "admin code2", Comparator: series.Eq, Comparando: lib.ToUpper(cityCode)}).
		Filter(dataframe.F{Colname: "postal code", Comparator: series.Eq, Comparando: postalCode})

	var found []string
	if fil.Nrow() == 0 {
		return fmt.Errorf("can't find any city for postal code %s and city code %s", postalCode, cityCode)
	} else {
		for i := 0; i < fil.Nrow(); i++ {
			found = append(found, fil.Elem(i, 1).String())
		}
	}
	for _, v := range found {
		if lib.ToUpper(normalizeString(v)) == lib.ToUpper(normalizeString(city)) {
			return nil
		}
	}
	return fmt.Errorf("city %s doesn't match any postal code %s and city code %s", city, postalCode, cityCode)
}

func normalizeString(in string) string {
	var out string
	for _, r := range in {
		var s string

		switch r {
		case ' ', '\'', '.', '/', '_', '-', '’', '`', '‘', '´':
			s = ""
		case 'è', 'é':
			s = "e"
		case 'à', 'ä':
			s = "a"
		case 'ò', 'ö':
			s = "o"
		case 'ì':
			s = "i"
		case 'ù', 'ü':
			s = "u"
		case 'ß':
			s = "ss"
		default:
			s = string(r)
		}
		out += s
	}

	return out
}

func checkFiscalCode(fc string) bool {
	return true
}

func checkDeclaredClaims(d []models.DeclaredClaims) error {
	if d == nil {
		return nil
	}
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
