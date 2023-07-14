package product

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"regexp"
	"time"
)

func GetProduct(name, version, role string) (*models.Product, error) {
	var (
		product  *models.Product
		filePath = "products/"
	)

	switch role {
	case models.UserRoleAdmin:
		filePath += "mga"
	case models.UserRoleAgency, models.UserRoleAgent:
		filePath += role
	default:
		filePath += "e-commerce"
	}
	filePath += "/" + name + "-" + version + ".json"

	jsonFile := lib.GetFilesByEnv(filePath)
	err := json.Unmarshal(jsonFile, &product)

	product, err = replaceDatesInProduct(product, role)

	return product, err
}

func replaceDatesInProduct(product *models.Product, role string) (*models.Product, error) {
	jsonOut, err := product.Marshal()
	if err != nil {
		return &models.Product{}, err
	}

	productJson := string(jsonOut)

	minAgeValue, minReservedAgeValue := ageMap[role][product.Name][minAge], ageMap[role][product.Name][minReservedAge]

	initialDate := time.Now().AddDate(-18, 0, 0).Format("2006-01-02")
	minDate := time.Now().AddDate(-minAgeValue, 0, 1).Format("2006-01-02")
	minReservedDate := time.Now().AddDate(-minReservedAgeValue, 0, 1).Format("2006-01-02")

	regexInitialDate := regexp.MustCompile("{{INITIAL_DATE}}")
	regexMinDate := regexp.MustCompile("{{MIN_DATE}}")
	regexMinAgentDate := regexp.MustCompile("{{MIN_RESERVED_DATE}}")

	productJson = regexInitialDate.ReplaceAllString(productJson, initialDate)
	productJson = regexMinDate.ReplaceAllString(productJson, minDate)
	productJson = regexMinAgentDate.ReplaceAllString(productJson, minReservedDate)

	err = json.Unmarshal([]byte(productJson), product)

	return product, err
}
