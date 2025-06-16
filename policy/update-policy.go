package policy

import (
	"fmt"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func UpdatePolicy(policy *models.Policy) (map[string]any, error) {
	input := make(map[string]any, 0)

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
	input["consultancyValue"] = map[string]any{
		"percentage": policy.ConsultancyValue.Percentage,
		"price":      lib.RoundFloat(policy.PriceGross*policy.ConsultancyValue.Percentage, 2),
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
		input["startDate"] = policy.StartDate
		input["endDate"] = policy.EndDate
		input["declaredClaims"] = policy.DeclaredClaims
		input["hasBond"] = policy.HasBond
		input["bond"] = policy.Bond
		input["clause"] = policy.Clause
		input["contractors"] = policy.Contractors
		input["priceGroup"] = policy.PriceGroup
	case models.CatNatProduct:
		input["startDate"] = policy.StartDate
		input["endDate"] = policy.EndDate
		input["quoteQuestions"] = policy.QuoteQuestions
		input["offersPrices"] = policy.OffersPrices
		input["contractors"] = policy.Contractors
		input["taxAmount"] = policy.TaxAmount
		input["priceNett"] = policy.PriceNett
		input["priceGross"] = policy.PriceGross
		input["paymentSplit"] = policy.PaymentSplit

		for _, a := range policy.Assets {
			if a.Building != nil && a.Building.BuildingAddress != nil {
				if err := validateAddress(*a.Building.BuildingAddress); err != nil {
					return nil, err
				}
			}
		}

		for _, c := range *policy.Contractors {
			if c.Residence != nil {
				if err := validateAddress(*c.Residence); err != nil {
					return nil, err
				}
			}
			if c.Domicile != nil {
				if err := validateAddress(*c.Domicile); err != nil {
					return nil, err
				}
			}
		}

		if policy.Contractor.Residence != nil {
			if err := validateAddress(*policy.Contractor.Residence); err != nil {
				return nil, err
			}
		}

		if policy.Contractor.Domicile != nil {
			if err := validateAddress(*policy.Contractor.Domicile); err != nil {
				return nil, err
			}
		}
	}

	input["updated"] = time.Now().UTC()

	return input, nil
}

func validateAddress(address models.Address) error {
	city := address.City
	postalCode := address.PostalCode
	cityCode := address.CityCode
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
