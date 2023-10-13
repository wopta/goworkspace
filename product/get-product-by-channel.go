package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

// DEPRECATED
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

	err = replaceDatesInProduct(product, channel)

	log.Println("[GetProduct] function end ---------------------")

	return product, err
}

func GetProductV2(productName, companyName, channel string) *models.Product {
	var (
		result, product *models.Product
		basePath        = "products-v2"
	)

	log.Println("[GetProductV2] function start ---------------------")

	log.Printf("[GetProductV2] product: %s", productName)

	filesList, err := lib.ListGoogleStorageFolderContent(fmt.Sprintf("%s/%s/", basePath, productName))
	if err != nil {
		log.Printf("[GetProduct] error: %s", err.Error())
		return nil
	}

	log.Println("[GetProductV2] filtering file list by channel")

	filesList = lib.SliceFilter(filesList, func(filePath string) bool {
		return strings.HasSuffix(filePath, fmt.Sprintf("%s.json", channel))
	})
	if len(filesList) == 0 {
		log.Println("[GetProductV2] empty file list")
		return nil
	}

	log.Println("[GetProductV2] sorting file list by version")

	sort.Slice(filesList, func(i, j int) bool {
		return strings.SplitN(filesList[i], "/", 4)[2] > strings.SplitN(filesList[j], "/", 4)[2]
	})

outerLoop:
	for _, file := range filesList {
		productBytes := lib.GetFilesByEnv(file)

		err = json.Unmarshal(productBytes, &product)
		if err != nil {
			log.Printf("[GetProductV2] error unmarshaling product: %s", err.Error())
			return nil
		}

		for _, company := range product.Companies {
			if company.Name == companyName && company.IsActive {
				result = product
				break outerLoop
			}
		}
	}

	if result == nil {
		log.Printf("[GetProductV2] no active %s product for %s company found", productName, companyName)
		return nil
	}

	log.Printf("[GetProductV2] productName: %s productVersion: %s channel: %s", product.Name, product.Version, channel)

	err = replaceDatesInProduct(result, channel)

	log.Println("[GetProductV2] function end ---------------------")

	return result
}

func replaceDatesInProduct(product *models.Product, channel string) error {
	if product == nil {
		return fmt.Errorf("no product found")
	}

	jsonOut, err := product.Marshal()
	if err != nil {
		return err
	}

	log.Println("[replaceDatesInProduct] function start -------------------")

	log.Printf("[replaceDatesInProduct] channel: %s", channel)

	productJson := string(jsonOut)

	minAgeValue, minReservedAgeValue := ageMap[channel][product.Name][minAge], ageMap[channel][product.Name][minReservedAge]

	log.Printf("[replaceDatesInProduct] minAge: %d minReservedAge: %d", minAgeValue, minReservedAgeValue)

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

	err = json.Unmarshal([]byte(productJson), &product)

	log.Println("[replaceDatesInProduct] function end -------------------")

	return err
}

// DEPRECATED
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
		err = replaceDatesInProduct(product, models.UserRoleAll)
	case "life":
		err = replaceDatesInProduct(product, models.UserRoleAll)
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
