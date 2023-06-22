package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

type GetProductListResp struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Logo    string `json:"logo"`
}

func GetProductsListByEntitlementFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)
	log.Println("GetProductsListByEntitlement")

	productsList := make([]models.Product, 0)

	res := lib.GetFolderContentByEnv("products/agency")

	//res := lib.ReadStorageDirContent(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/agency/")

	for _, file := range res {
		var product models.Product
		err = json.Unmarshal(file, &product)
		productsList = append(productsList, product)
	}

	jsonOut, err := json.Marshal(productsList)

	return string(jsonOut), productsList, err
}
