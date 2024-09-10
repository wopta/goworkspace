package product

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func printSlice(stringSlice []string) int {
	tot := 0
	for _, v := range stringSlice {
		fmt.Printf("%v\n", v)
		tot++
	}
	return tot
}

func GetAllProductsByChannel(channel string) []models.ProductInfo {
	log.Println("[GetAllProductsByChannel] function start -----------------------")

	var products = make([]models.ProductInfo, 0)

	fileList := getProductsFileList()

	// filter only the files for the given channel
	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel)
	})

	products = getProductsFromFileList(fileList)

	log.Printf("[GetAllProductsByChannel] found %d products", len(products))
	log.Println("[GetAllProductsByChannel] function end -------------------------")

	return products
}

func GetProductsByChannel(productList []string, channel string) []models.ProductInfo {
	log.Println("[GetNetworkNodeProducts] function start ---------------------")

	var products = make([]models.ProductInfo, 0)

	fileList := getProductsFileList()

	// filter only the files for the network channel present on the product list
	fileList = lib.SliceFilter(fileList, func(file string) bool {
		filenameParts := strings.SplitN(file, "/", 4)
		return strings.HasPrefix(filenameParts[3], channel) &&
			lib.SliceContains[string](productList, filenameParts[1])
	})

	products = getProductsFromFileList(fileList)

	log.Printf("[GetNetworkNodeProducts] found %d products", len(products))
	log.Println("[GetNetworkNodeProducts] function end -----------------------")

	return products
}

func getProductsFromFileList(fileList []string) []models.ProductInfo {
	var (
		err            error
		products       = make([]models.ProductInfo, 0)
		fileChunks     = make([][]string, 0)
		currentProduct models.BaseProduct
	)

	if len(fileList) == 0 {
		log.Println("[getProductsFromFileList] error empty fileList")
		return products
	}

	log.Printf("[getProductsFromFileList] checking fileList %v", fileList)

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
		log.Printf("[getProductsFromFileList] sorted chunk %v", chunk)
		// loop each version
		for _, file := range chunk {
			// download file from bucket
			log.Printf("[getProductsFromFileList] downloading file %s", file)
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

	switch os.Getenv("env") {
	case "local":
		fileList, err = lib.ListLocalFolderContent(models.ProductsFolder)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(models.ProductsFolder)
	}

	if err != nil {
		log.Printf("[GetNetworkNodeProducts] error getting file list: %s", err.Error())
	}

	/*
		COUNT SEPARATORS
	*/
	var filteredFileList []string
	var removedElements []string
	log.Println("=== ORIGINAL LIST ===")
	tot := printSlice(fileList)
	log.Printf("There are %d elements", tot)
	start := time.Now()
	for _, s := range fileList {
		if strings.Count(s, "/") == 3 {
			filteredFileList = append(filteredFileList, s)
		} else {
			removedElements = append(removedElements, s)
		}
	}
	elapsed := time.Since(start)
	log.Println("=== COUNT SEPARATORS FILTERED LIST ===")
	tot = printSlice(filteredFileList)
	log.Printf("There are %d elements", tot)
	log.Printf("### ==> COUNT SEPARATORS took %s", elapsed)

	/*
		REGEXP
	*/
	pattern := regexp.MustCompile(`/v[1-9]+/`)
	filteredFileList = nil
	removedElements = nil
	start = time.Now()
	for _, s := range fileList {
		if pattern.MatchString(s) {
			filteredFileList = append(filteredFileList, s)
		} else {
			removedElements = append(removedElements, s)
		}
	}

	elapsed = time.Since(start)
	log.Println("=== REGEXP FILTERED LIST ===")
	tot = printSlice(filteredFileList)
	log.Printf("There are %d elements", tot)
	log.Printf("### ==> REGEXP took %s", elapsed)

	return filteredFileList
}
