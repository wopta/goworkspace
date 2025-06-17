package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func MailValidationFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
	)
	log.AddPrefix("MailValidation")
	defer log.PopPrefix()

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.ErrorF("error unmarshalling policy: %s", err.Error())
		return "", nil, err
	}
	err = validateAllAddresses(&policy)
	if err != nil {
		return "", nil, errors.New("Errore nell'indirizzo")
	}
	return "{}", nil, nil
}

func validateAllAddresses(policy *models.Policy) error {
	for _, a := range policy.Assets {
		if a.Building != nil {
			if err := validateAddress(a.Building.BuildingAddress); err != nil {
				return err
			}
		}
	}

	for _, contractor := range *policy.Contractors {
		if err := validateAddress(contractor.Residence); err != nil {
			return err
		}
		if err := validateAddress(contractor.Domicile); err != nil {
			return err
		}
	}

	if err := validateAddress(policy.Contractor.Residence); err != nil {
		return err
	}

	if err := validateAddress(policy.Contractor.Domicile); err != nil {
		return err
	}

	if err := validateAddress(policy.Contractor.CompanyAddress); err != nil {
		return err
	}
	return nil
}

func validateAddress(address *models.Address) error {
	if address == nil {
		return nil
	}
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
