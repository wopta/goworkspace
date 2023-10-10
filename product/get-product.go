package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetProduct(name, version, channel string) (*models.Product, error) {
	var (
		product  *models.Product
		filePath = "products/"
	)

	log.Println("[GetProduct] function start ---------------------")

	filePath += channel + "/" + name + "-" + version + ".json"

	log.Printf("[GetProduct] product filePath: %s", filePath)

	jsonFile := lib.GetFilesByEnv(filePath)
	err := json.Unmarshal(jsonFile, &product)
	lib.CheckError(err)

	product, err = replaceDatesInProduct(product, channel)

	log.Println("[GetProduct] function end ---------------------")

	return product, err
}

func replaceDatesInProduct(product *models.Product, channel string) (*models.Product, error) {
	jsonOut, err := product.Marshal()
	if err != nil {
		return &models.Product{}, err
	}

	log.Println("[replaceDatesInProduct] function start -------------------")

	log.Printf("[replaceDatesInProduct] channel: %s", channel)

	productJson := string(jsonOut)

	minAgeValue, minReservedAgeValue := ageMap[channel][product.Name][minAge], ageMap[channel][product.Name][minReservedAge]

	log.Printf("[replaceDatesInProduct] minAgeValue: %d minReservedAgeValue: %d", minAgeValue, minReservedAgeValue)

	initialDate := time.Now().AddDate(-18, 0, 0).Format(models.TimeDateOnly)
	minDate := time.Now().AddDate(-minAgeValue, 0, 1).Format(models.TimeDateOnly)
	minReservedDate := time.Now().AddDate(-minReservedAgeValue, 0, 1).Format(models.TimeDateOnly)
	startDate := time.Now().Format(models.TimeDateOnly)
	maxStartDate := time.Now().AddDate(0, 0, 30).Format(models.TimeDateOnly)

	regexInitialDate := regexp.MustCompile("{{INITIAL_DATE}}")
	regexMinDate := regexp.MustCompile("{{MIN_DATE}}")
	regexMinAgentDate := regexp.MustCompile("{{MIN_RESERVED_DATE}}")
	regexStartDate := regexp.MustCompile("{{START_DATE}}")
	regexMaxStartDate := regexp.MustCompile("{{MAX_START_DATE}}")

	productJson = regexInitialDate.ReplaceAllString(productJson, initialDate)
	productJson = regexMinDate.ReplaceAllString(productJson, minDate)
	productJson = regexMinAgentDate.ReplaceAllString(productJson, minReservedDate)
	productJson = regexStartDate.ReplaceAllString(productJson, startDate)
	productJson = regexMaxStartDate.ReplaceAllString(productJson, maxStartDate)

	err = json.Unmarshal([]byte(productJson), product)

	log.Println("[replaceDatesInProduct] function end -------------------")

	return product, err
}

// TODO: remove this endpoint once in production
func GetNameFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	name := r.Header.Get("name")
	origin := r.Header.Get("origin")

	log.Println(r.RequestURI)

	product, err := GetName(origin, name, "v1")
	if err != nil {

		return "", nil, err
	}
	jsonOut, err := product.Marshal()
	if err != nil {
		return "", nil, err
	}

	switch name {
	case "persona":
		product, err = replaceDatesInProduct(product, models.UserRoleAll)
	case "life":
		product, err = replaceDatesInProduct(product, models.UserRoleAll)
	}

	jsonOut, err = json.Marshal(product)

	return string(jsonOut), product, err
}

// TODO: remove this endpoint once in production
func GetName(origin string, name string, version string) (*models.Product, error) {
	q := lib.Firequeries{
		Queries: []lib.Firequery{{
			Field:      "name",
			Operator:   "==",
			QueryValue: name,
		},
			{
				Field:      "version",
				Operator:   "==",
				QueryValue: version,
			},
		},
	}

	fireProduct := lib.GetDatasetByEnv(origin, "products")
	query, _ := q.FirestoreWherefields(fireProduct)
	products := models.ProductToListData(query)
	if len(products) == 0 {
		return &models.Product{}, fmt.Errorf("no product json file found for %s %s", name, version)
	}

	return &products[0], nil
}
