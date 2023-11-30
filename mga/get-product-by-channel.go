package mga

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

type GetProductReq struct {
	ProductName string `json:"name"`
	CompanyName string `json:"company"` // DEPRECATED
	Version     string `json:"version"` // DEPRECATED
}

func GetProductByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req         *GetProductReq
		networkNode *models.NetworkNode
		warrant     *models.Warrant
	)

	log.SetPrefix("[GetProductByChannelFx] ")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("body req: %s", string(body))

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("error unmarshaling request body")
		return "", nil, err
	}

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", nil, err
	}

	channel := authToken.GetChannelByRoleV2()
	log.Printf("channel: %s", channel)

	if strings.EqualFold(channel, models.NetworkChannel) {
		networkNode = network.GetNetworkNodeByUid(authToken.UserID)
		if networkNode != nil {
			warrant = networkNode.GetWarrant()
		}
		if warrant != nil && !networkNode.HasAccessToProduct(req.ProductName, warrant) {
			log.Printf("network node %s hasn't access to product %s", networkNode.Uid, req.ProductName)
			return "", nil, fmt.Errorf("network node hasn't access to product")
		}
	}

	log.Printf("getting last active action for product %s", req.ProductName)

	product := prd.GetLatestActiveProduct(req.ProductName, channel, networkNode, warrant)
	if product == nil {
		log.Printf("no active product found")
		return "", nil, fmt.Errorf("no product active found")
	}

	jsonOut, err := product.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), product, err
}
