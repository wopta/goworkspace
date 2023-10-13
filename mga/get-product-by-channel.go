package mga

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
	"strings"
)

type GetProductReq struct {
	ProductName string `json:"productName"`
	CompanyName string `json:"companyName"`
}

func GetProductByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req         *GetProductReq
		networkNode *models.NetworkNode
	)

	log.Println("[GetProductByChannelFx] handler start -------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[GetProductByChannelFx] body req: %s", string(body))

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("[GetProductByChannelFx] error unmarshaling request body")
		return "", nil, err
	}

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[GetProductFx] error: %s", err.Error())
		return "", nil, err
	}

	channel := authToken.GetChannelByRoleV2()
	log.Printf("[GetProductByChannelFx] channel: %s", channel)

	if strings.EqualFold(channel, models.NetworkChannel) {
		networkNode = network.GetNetworkNodeByUid(authToken.UserID)
		if networkNode != nil && !networkNode.HasAccessToProduct(req.ProductName) {
			log.Printf("[GetProductByChannelFx] network node %s hasn't access to product %s for company %s", networkNode.Uid, req.ProductName, req.CompanyName)
			return "", nil, fmt.Errorf("network node hasn't access to product")
		}
	}

	log.Printf("[GetProductByChannelFx] getting last active action for product %s", req.ProductName)

	product := prd.GetProductV2(req.ProductName, channel, networkNode)
	if product == nil {
		log.Printf("[GetProductByChannelFx] no active product found")
		return "", nil, fmt.Errorf("no product active found")
	}

	jsonOut, err := product.Marshal()

	log.Println("[GetProductByChannelFx] handler end -------------")

	return string(jsonOut), product, err
}
