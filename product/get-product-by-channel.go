package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	minAge         = "minAge"
	minReservedAge = "minReservedAge"
)

/*
Returns the requested product version for the specified channel based on the provided input parameters, including
product name, version, and channel.
*/
func GetProductV2(productName, productVersion, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) *models.Product {
	var (
		product *models.Product
	)

	log.Println("[GetProductV2] function start -----------------")

	filePath := fmt.Sprintf("%s/%s/%s/%s.json", models.ProductsFolder, productName, productVersion, channel)

	log.Printf("[GetProductV2] filePath: %s", filePath)

	productBytes := lib.GetFilesByEnv(filePath)
	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, productBytes)

	log.Printf("[GetProductV2] retrieved product: %s", buffer.String())

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

	overrideProductInfo(product, networkNode, warrant, channel)

	return product
}

/*
Returns the most recent active default product associated with the specified channel.
*/
func getDefaultProduct(productName, channel string) *models.Product {
	var (
		result, product *models.Product
	)

	log.Println("[GetDefaultProduct] function start --------------")

	filesList, err := lib.ListGoogleStorageFolderContent(fmt.Sprintf("%s/%s/", models.ProductsFolder, productName))
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

	product.Steps = loadProductSteps(product)

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

	product = getDefaultProduct(productName, channel)
	if product == nil {
		log.Printf("[GetLatestActiveProduct] no active product found")
		return nil
	}

	overrideProductInfo(product, networkNode, warrant, channel)

	err := replaceDatesInProduct(product, channel)
	if err != nil {
		log.Printf("[GetLatestActiveProduct] error replacing dates in product: %s", err.Error())
		return nil
	}

	log.Println("[GetLatestActiveProduct] function end ---------------------")

	return product
}

func GetLifeAgeInfo(productName, productVersion, channel string) (int, int) {
	var ageMap map[string]map[string]int

	rawMap := lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/age_info.json", models.ProductsFolder, productName, productVersion))
	err := json.Unmarshal(rawMap, &ageMap)
	if err != nil {
		return 0, 0
	}
	if ageMap[channel] != nil {
		return ageMap[channel]["minAge"], ageMap[channel][minReservedAge]
	}
	return 0, 0
}

func getGapAgeInfo(productName, productVersion, channel string) (minContractorAge, minAssetPersonAge int) {
	var ageMap map[string]map[string]int

	rawMap := lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/age_info.json", models.ProductsFolder, productName, productVersion))
	err := json.Unmarshal(rawMap, &ageMap)
	if err != nil {
		return 0, 0
	}
	if ageMap[channel] != nil {
		return ageMap[channel]["minContractorAge"], ageMap[channel]["minAssetPersonAge"]
	}
	return 0, 0
}

func replaceDatesInProduct(product *models.Product, channel string) error {
	if product == nil {
		return fmt.Errorf("no product found")
	}

	var err error

	log.Println("[replaceDatesInProduct] function start ----------------------")

	switch product.Name {
	case models.LifeProduct, models.PersonaProduct:
		err = replaceLifeDates(product, channel)
	case models.GapProduct:
		err = replaceGapDates(product, channel)
	default:
		log.Printf("[replaceDatesInProduct] product %s does not have dates to be replaced", product.Name)
	}

	log.Println("[replaceDatesInProduct] function end ------------------------")

	return err
}

func replaceLifeDates(product *models.Product, channel string) error {
	jsonOut, err := product.Marshal()
	if err != nil {
		return err
	}

	productJson := string(jsonOut)

	minAgeValue, minReservedAgeValue := GetLifeAgeInfo(product.Name, product.Version, channel)

	log.Printf("[replaceLifeDates] minAge: %d minReservedAge: %d", minAgeValue, minReservedAgeValue)

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

	return json.Unmarshal([]byte(productJson), &product)
}

func replaceGapDates(product *models.Product, channel string) error {
	jsonOut, err := product.Marshal()
	if err != nil {
		return err
	}

	productJson := string(jsonOut)

	minContractorAgeValue, minAssetPersonAgeValue := getGapAgeInfo(product.Name, product.Version, channel)

	log.Printf("[replaceGapDates] minContractorAge: %d minAssetPersonAge: %d", minContractorAgeValue, minAssetPersonAgeValue)

	maxContractorBirthDate := time.Now().AddDate(-minContractorAgeValue, 0, 0).Format(models.TimeDateOnly)
	maxAssetPersonBirthDate := time.Now().AddDate(-minAssetPersonAgeValue, 0, 0).Format(models.TimeDateOnly)

	regexMaxContractorBirthDate := regexp.MustCompile("{{MAX_CONTRACTOR_BIRTH_DATE}}")
	regexMaxAssetPersonBirthDate := regexp.MustCompile("{{MAX_ASSET_PERSON_BIRTH_DATE}}")

	productJson = regexMaxContractorBirthDate.ReplaceAllString(productJson, maxContractorBirthDate)
	productJson = regexMaxAssetPersonBirthDate.ReplaceAllString(productJson, maxAssetPersonBirthDate)

	return json.Unmarshal([]byte(productJson), &product)
}

func overrideProductInfo(product *models.Product, networkNode *models.NetworkNode, warrant *models.Warrant, channel string) {
	if networkNode == nil {
		product.Steps = filterProductStepsByFlow(product.Steps, channel)
		return
	}

	var flow = channel

	if networkNode.HasAccessToProduct(product.Name, warrant) {
		if warrant != nil {
			paymentProviders := make([]models.PaymentProvider, 0)
			if warrantProduct := warrant.GetProduct(product.Name); warrantProduct != nil {
				for _, paymentProvider := range product.PaymentProviders {
					if lib.SliceContains(paymentProvider.Flows, warrantProduct.Flow) {
						paymentProviders = append(paymentProviders, paymentProvider)
					}
				}
				if warrantProduct.ConsultancyConfig != nil {
					product.ConsultancyConfig = warrantProduct.ConsultancyConfig
				}
			}

			// TODO: this need to be removed in the future
			paymentProviders = removeFacilePaymentRate(paymentProviders, warrant.Name)
			product.PaymentProviders = paymentProviders
		}

		for _, nodeProduct := range networkNode.Products {
			if nodeProduct.Name != product.Name {
				continue
			}
			if len(nodeProduct.Steps) > 0 {
				product.Steps = nodeProduct.Steps
			}
			if nodeProduct.ConsultancyConfig != nil {
				product.ConsultancyConfig = nodeProduct.ConsultancyConfig
			}
		}

		if networkNode.Type != models.PartnershipNetworkNodeType {
			flow = warrant.GetFlowName(product.Name)
		}
		product.Steps = filterProductStepsByFlow(product.Steps, flow)
	}
}

func filterProductStepsByFlow(steps []models.Step, flowName string) []models.Step {
	outputSteps := make([]models.Step, 0)
	for _, step := range steps {
		if len(step.Flows) == 0 || lib.SliceContains(step.Flows, flowName) {
			outputSteps = append(outputSteps, step)
		}
	}
	return outputSteps
}

func loadProductSteps(product *models.Product) []models.Step {
	var steps []models.Step
	rawSteps := lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/builder_ui.json", models.ProductsFolder, product.Name, product.Version))
	_ = json.Unmarshal(rawSteps, &steps)

	return steps
}

func removeFacilePaymentRate(paymentProviders []models.PaymentProvider, warrantName string) []models.PaymentProvider {
	if lib.SliceContains([]string{"facile_agent"}, warrantName) {
		for index, paymentProvider := range paymentProviders {
			configs := make([]models.PaymentConfig, 0)
			for _, config := range paymentProvider.Configs {
				if config.Rate != string(models.PaySplitMonthly) {
					configs = append(configs, config)
				}
			}
			paymentProviders[index].Configs = configs
		}
	}
	return paymentProviders
}
