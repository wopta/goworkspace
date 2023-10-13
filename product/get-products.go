package product

import (
	"encoding/json"
	"log"
	"sort"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	productFolderPath = "products-v2"
)

func GetProductsByChannel(channel string) []models.ProductInfo {
	log.Println("[GetProductsByChannel] function start -----------------------------")
	var (
		err            error
		products       []models.ProductInfo = make([]models.ProductInfo, 0)
		fileChunks     [][]string           = make([][]string, 0)
		currentProduct models.Product
	)

	// TODO use lib function
	fileList, err := GetFileListMock(productFolderPath)

	// filter only the files for the given channel
	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel)
	})

	// create subarrays for each different product
	for _, file := range fileList {
		filenameParts := strings.SplitN(file, "/", 4)
		if len(fileChunks) == 0 {
			chunk := make([]string, 0)
			fileChunks = append(fileChunks, chunk)
		}
		if len(fileChunks[len(fileChunks)-1]) > 0 && strings.SplitN(fileChunks[len(fileChunks)-1][0], "/", 3)[1] != filenameParts[1] {
			chunk := make([]string, 0)
			fileChunks = append(fileChunks, chunk)
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
			// download file from bucket
			fileBytes := lib.GetFilesByEnv(file)
			if err != nil {
				log.Printf("[GetProductsByChannel] error getting file: %s", err.Error())
				break
			}
			err = json.Unmarshal(fileBytes, &currentProduct)
			if err != nil {
				log.Printf("[GetProductsByChannel] error unmarshaling file: %s", err.Error())
				break
			}
			// loop all companies for that product/version
			for _, company := range currentProduct.Companies {
				// if active add to response
				if company.IsActive {
					log.Printf("[GetProductsByChannel] adding %s %s %s", currentProduct.Name, currentProduct.Version, company.Name)
					products = append(products, models.ProductInfo{
						Name:         currentProduct.Name,
						NameTitle:    currentProduct.NameTitle,
						NameSubtitle: currentProduct.NameSubtitle,
						NameDesc:     *currentProduct.NameDesc,
						Version:      currentProduct.Version,
						Company:      company.Name,
						Logo:         currentProduct.Logo,
					})
				}
			}
		}
		if err != nil {
			break
		}
	}

	log.Printf("[GetProductsByChannel] found %d products", len(products))
	log.Println("[GetProductsByChannel] function end -------------------------------")

	return products
}

func GetNetworkNodeProducts(productList []string) []models.ProductInfo {
	log.Println("[GetNetworkNodeProducts] function start ---------------------")
	var (
		// err      error
		products []models.ProductInfo
	)

	log.Printf("[GetNetworkNodeProducts] found %d products", len(products))
	log.Println("[GetNetworkNodeProducts] function end -----------------------")

	return products
}

func GetFileListMock(folderPath string) ([]string, error) {
	return []string{
		"products-v2/life/v1/mga.json",
		"products-v2/life/v1/e-commerce.json",
		"products-v2/life/v1/network.json",
		"products-v2/life/v2/mga.json",
		"products-v2/life/v2/e-commerce.json",
		"products-v2/life/v2/network.json",
		"products-v2/gap/v1/mga.json",
		"products-v2/gap/v1/network.json",
		"products-v2/persona/v1/mga.json",
		"products-v2/persona/v1/e-commerce.json",
		"products-v2/persona/v1/network.json",
	}, nil
}
