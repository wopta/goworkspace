package partnership

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
	"gopkg.in/square/go-jose.v2"
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

	err := decryptJwt(jwtData, os.Getenv(key), affinity)
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

// TODO: mode to lib
func decryptJwt(jwtData, key, partnershipName string) error {
	if partnershipName == "facile" {
		object, err := jose.ParseEncrypted(jwtData)
		if err != nil {
			log.Printf("[DecryptJwt] could not parse jwt - %s", err.Error())
			return fmt.Errorf("could not parse jwt")
		}

		decryptionKey, err := b64.StdEncoding.DecodeString(key)
		if err != nil {
			log.Printf("[DecryptJwt] could not decode signing key - %s", err.Error())
			return fmt.Errorf("could not decode jwt key")
		}

		_, err = object.Decrypt(decryptionKey)
		if err != nil {
			log.Printf("[DecryptJwt] could not decrypt jwt - %s", err.Error())
			return fmt.Errorf("could not decrypt jwt")
		}

		return nil
	}

	_, err := jwt.Parse(jwtData, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key, e := b64.StdEncoding.DecodeString(key)

		return []byte(key), e
	})

	return err
}
