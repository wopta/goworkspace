package renew

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"os"
	"sort"
	"strings"
)

func getProductsMapByPolicyType(policyType, quoteType string) map[string]models.Product {
	products := make(map[string]models.Product)

	productsList := getProducts()

	for _, prd := range productsList {
		if strings.EqualFold(prd.PolicyType, policyType) && strings.EqualFold(prd.QuoteType, quoteType) {
			key := fmt.Sprintf("%s-%s", prd.Name, prd.Version)
			products[key] = prd
		}
	}

	return products
}

func getProducts() []models.Product {
	const channel = models.MgaChannel
	var products = make([]models.Product, 0)

	fileList := getProductsFileList()

	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel)
	})

	products = getProductsFromFileList(fileList)

	return products
}

func getProductsFileList() []string {
	var (
		err      error
		fileList = make([]string, 0)
	)

	switch os.Getenv("env") {
	case "local", "local-test":
		fileList, err = lib.ListLocalFolderContent(models.ProductsFolder)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(models.ProductsFolder)
	}

	if err != nil {
		log.Printf("[GetNetworkNodeProducts] error getting file list: %s", err.Error())
	}

	return fileList
}

func getProductsFromFileList(fileList []string) []models.Product {
	var (
		err        error
		products   = make([]models.Product, 0)
		fileChunks = make([][]string, 0)
	)

	if len(fileList) == 0 {
		return products
	}

	// create subarrays for each different product
	for _, file := range fileList {
		filenameParts := strings.SplitN(file, "/", 4)
		productName := filenameParts[1]
		if len(fileChunks) == 0 {
			fileChunks = append(fileChunks, make([]string, 0))
		}
		if len(fileChunks[len(fileChunks)-1]) > 0 {
			chunkProductName := strings.SplitN(fileChunks[len(fileChunks)-1][0], "/", 3)[1]
			if chunkProductName != productName {
				fileChunks = append(fileChunks, make([]string, 0))
			}
		}
		fileChunks[len(fileChunks)-1] = append(fileChunks[len(fileChunks)-1], file)
	}

	// loop each product
	for _, chunk := range fileChunks {
		// sort them by the last version
		sort.Slice(chunk, func(i, j int) bool {
			return strings.SplitN(chunk[i], "/", 4)[2] > strings.SplitN(chunk[j], "/", 4)[2]
		})
		// loop each version
		for _, file := range chunk {
			var currentProduct models.Product
			// download file from bucket
			fileBytes := lib.GetFilesByEnv(file)
			err = json.Unmarshal(fileBytes, &currentProduct)
			if err != nil {
				continue
			}

			products = append(products, currentProduct)
		}
		if err != nil {
			break
		}
	}

	return products
}
