package mga

import (
	"encoding/json"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"log"
	"net/http"
)

type GetActiveProductsByChannelResp struct {
	Products []models.Product `json:"products"`
}

func GetActiveProductsByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		resp GetActiveProductsByChannelResp
	)

	log.Println("[GetActiveProductsByChannelFx]")

	channel := r.Header.Get("channel")

	resp.Products = product.GetAllProductsByChannel(channel)

	jsonOut, err := json.Marshal(resp)

	return string(jsonOut), resp, err
}
