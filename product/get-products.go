package product

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"os"
	"sort"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetAllProductsByChannel(channel string) []models.ProductInfo {
	log.AddPrefix("GetAllProductsByChannel")
	defer log.PopPrefix()
	log.Println("function start -----------------------")

	var products = make([]models.ProductInfo, 0)

	fileList := getProductsFileList()

	// filter only the files for the given channel
	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel)
	})

	products = getProductsFromFileList(fileList)

	log.Printf("found %d products", len(products))
	log.Println("function end -------------------------")

	return products
}

func GetProductsByChannel(productList []string, channel string) []models.ProductInfo {
	log.AddPrefix("GetNetworkNodeProducts")
	defer log.PopPrefix()
	log.Println("function start ---------------------")

	var products = make([]models.ProductInfo, 0)

	fileList := getProductsFileList()

	// filter only the files for the network channel present on the product list
	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel) &&
			lib.SliceContains[string](productList, filenameParts[1])
	})

	products = getProductsFromFileList(fileList)

	log.Printf("found %d products", len(products))
	log.Println("function end -----------------------")

	return products
}

func getProductsFromFileList(fileList []string) []models.ProductInfo {
	var (
		err            error
		products       = make([]models.ProductInfo, 0)
		fileChunks     = make([][]string, 0)
		currentProduct models.BaseProduct
	)
	log.AddPrefix("GetProductsFromFileList")
	defer log.PopPrefix()
	if len(fileList) == 0 {
		log.Println("error empty fileList")
		return products
	}

	log.Printf("checking fileList %v", fileList)

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
		log.Printf("sorted chunk %v", chunk)
		// loop each version
		for _, file := range chunk {
			// download file from bucket
			log.Printf("downloading file %s", file)
			fileBytes := lib.GetFilesByEnv(file)
			if err != nil {
				log.ErrorF("error getting file: %s", err.Error())
				break
			}
			err = json.Unmarshal(fileBytes, &currentProduct)
			if err != nil {
				log.ErrorF("error unmarshaling file: %s", err.Error())
				break
			}

			if currentProduct.IsActive {
				products = append(products, currentProduct.ToProductInfo())
			}
		}
		if err != nil {
			break
		}
	}

	return products
}

func getProductsFileList() []string {
	var (
		err      error
		fileList = make([]string, 0)
	)
	log.AddPrefix("GetNetworkNodeProducts")
	switch os.Getenv("env") {
	case "local":
		fileList, err = lib.ListLocalFolderContent(models.ProductsFolder)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(models.ProductsFolder)
	}

	if err != nil {
		log.ErrorF("error getting file list: %s", err.Error())
	}

	checkedList := lib.SliceFilter(fileList, checkSlashes)

	return checkedList
}

func checkSlashes(s string) bool {
	// Correct path is: products/{{product_dir}}/{{version_number}}/{{filename.extension}}
	// but this function supports further nesting
	return strings.Count(s, "/") >= 3
}
