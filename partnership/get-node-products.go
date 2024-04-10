package partnership

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

type Response struct {
	Partnership models.PartnershipNode `json:"partnership"`
	Products    []models.ProductInfo   `json:"products"`
}

func GetPartnershipNodeAndProductsFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var response Response

	log.SetPrefix("[GetPartnershipNodeAndProductsFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	affinity := chi.URLParam(r, "partnershipName")
	jwtData := r.URL.Query().Get("jwt")
	key := lib.ToUpper(fmt.Sprintf("%s_SIGNING_KEY", affinity))

	_, err := lib.ParseJwt(jwtData, os.Getenv(key), encryptedPartnerships[affinity])
	if err != nil {
		log.Printf("error decoding jwt: %s", err.Error())
		return "", nil, err
	}

	node, err := network.GetNodeByUid(affinity)
	if err != nil {
		log.Printf("error getting node '%s': %s", affinity, err.Error())
		return "", nil, err
	}
	response.Partnership = *node.Partnership

	productList := lib.SliceMap(node.Products, func(p models.Product) string { return p.Name })
	response.Products = product.GetProductsByChannel(productList, lib.ECommerceChannel)

	responseJson, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}
