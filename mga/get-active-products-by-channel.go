package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/product"
	"log"
	"net/http"
)

func GetActiveProductsByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetActiveProductsByChannelFx]")

	channel := r.Header.Get("channel")

	products := product.GetAllProductsByChannel(channel)

	jsonOut, err := json.Marshal(products)

	return string(jsonOut), products, err
}
