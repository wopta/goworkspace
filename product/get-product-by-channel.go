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

const (
	basePath = "products-v2"
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

/*
Returns the requested product version for the specified channel based on the provided input parameters, including
product name, version, and channel.
*/
func GetProductV2(productName, productVersion, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) *models.Product {
	var (
		product *models.Product
	)

	log.Println("[GetProductV2] function start -----------------")

	filePath := fmt.Sprintf("%s/%s/%s/%s.json", basePath, productName, productVersion, channel)

	log.Printf("[GetProductV2] filePath: %s", filePath)

	productBytes := lib.GetFilesByEnv(filePath)

	log.Printf("[GetProductV2] retrieved product: %s", string(productBytes))

	err := json.Unmarshal(productBytes, &product)
	if err != nil {
		log.Printf("[GetProductV2] error unmarshaling product: %s", err.Error())
		return nil
	}

	err = replaceDatesInProduct(product, channel)
	if err != nil {
		log.Printf("[GetProductV2] error during replace dates in product: %s", err.Error())
		return nil
	}

	overrideProductInfo(product, networkNode, warrant)

	return product
}

/*
Returns the most recent active default product associated with the specified channel.
*/
func GetDefaultProduct(productName, channel string) *models.Product {
	var (
		result, product *models.Product
	)

	log.Println("[GetDefaultProduct] function start --------------")

	filesList, err := lib.ListGoogleStorageFolderContent(fmt.Sprintf("%s/%s/", basePath, productName))
	if err != nil {
		log.Printf("[GetProduct] error: %s", err.Error())
		return nil
	}

	log.Println("[GetDefaultProduct] filtering file list by channel")

	filesList = lib.SliceFilter(filesList, func(filePath string) bool {
		return strings.HasSuffix(filePath, fmt.Sprintf("%s.json", channel))
	})
	if len(filesList) == 0 {
		log.Println("[GetDefaultProduct] empty file list")
		return nil
	}

	log.Println("[GetDefaultProduct] sorting file list by version")

	sort.Slice(filesList, func(i, j int) bool {
		return strings.SplitN(filesList[i], "/", 4)[2] > strings.SplitN(filesList[j], "/", 4)[2]
	})

	for _, file := range filesList {
		productBytes := lib.GetFilesByEnv(file)

		err = json.Unmarshal(productBytes, &product)
		if err != nil {
			log.Printf("[GetDefaultProduct] error unmarshaling product: %s", err.Error())
			return nil
		}

		if product.IsActive {
			log.Printf("[GetDefaultProduct] product %s version %s is active", product.Name, product.Version)
			result = product
			break
		}
		log.Printf("[GetDefaultProduct] product %s version %s is not active", product.Name, product.Version)
	}

	if result == nil {
		log.Printf("[GetDefaultProduct] no active %s product found", productName)
		return nil
	}

	err = replaceDatesInProduct(result, channel)

	log.Println("[GetDefaultProduct] function end ---------------------")

	return result
}

/*
Returns the latest active default product linked to the specified channel. If a network node is provided and is
not nil, and it possesses a custom journey product, the defined steps will take precedence over the defaults.
Furthermore, if the network node has a warrant, payment providers will be filtered based on the specified flow
for the requested product.
*/
func GetLatestActiveProduct(productName, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) *models.Product {
	var (
		product *models.Product
	)

	log.Println("[GetLatestActiveProduct] function start ---------------------")

	log.Printf("[GetLatestActiveProduct] product: %s", productName)

	product = GetDefaultProduct(productName, channel)
	if product == nil {
		log.Printf("[GetLatestActiveProduct] no active product found")
		return nil
	}

	overrideProductInfo(product, networkNode, warrant)

	log.Println("[GetLatestActiveProduct] function end ---------------------")

	return product
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

func overrideProductInfo(product *models.Product, networkNode *models.NetworkNode, warrant *models.Warrant) {
	if networkNode == nil {
		return
	}

	if networkNode.HasAccessToProduct(product.Name, warrant) {
		for _, nodeProduct := range networkNode.Products {
			if nodeProduct.Name == product.Name && len(nodeProduct.Steps) > 0 {
				log.Printf("[GetLatestActiveProduct] overriding steps for product %s", product.Name)
				product.Steps = nodeProduct.Steps
			}
		}

		if warrant != nil {
			paymentProviders := make([]models.PaymentProvider, 0)
			warrantProduct := warrant.GetProduct(product.Name)
			if warrantProduct != nil {
				for _, paymentProvider := range product.PaymentProviders {
					if lib.SliceContains(paymentProvider.Flows, warrantProduct.Flow) {
						paymentProviders = append(paymentProviders, paymentProvider)
					}
				}
			}
			product.PaymentProviders = paymentProviders
		}
	}
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
	case models.PersonaProduct:
		err = replaceDatesInProduct(product, models.UserRoleAll)
	case models.LifeProduct:
		err = replaceDatesInProduct(product, models.UserRoleAll)
	}

	jsonOut, err = json.Marshal(product)

	return string(jsonOut), product, err
}

// DEPRECATED
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
