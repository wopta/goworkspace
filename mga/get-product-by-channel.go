package mga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

type GetProductReq struct {
	ProductName     string `json:"name"`
	CompanyName     string `json:"company"` // DEPRECATED
	Version         string `json:"version"` // DEPRECATED
	PartnershipName string `json:"partnershipName"`
}

func GetProductByChannelFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req         *GetProductReq
		networkNode *models.NetworkNode
		warrant     *models.Warrant
	)

	log.AddPrefix("GetProductByChannelFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error unmarshaling request body")
		return "", nil, err
	}

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	channel := authToken.GetChannelByRoleV2()
	log.Printf("channel: %s", channel)

	nodeUid := req.PartnershipName
	if strings.EqualFold(channel, models.NetworkChannel) {
		nodeUid = authToken.UserID
	}

	networkNode = network.GetNetworkNodeByUid(nodeUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	if warrant != nil && !networkNode.HasAccessToProduct(req.ProductName, warrant) {
		log.Printf("network node %s hasn't access to product %s", networkNode.Uid, req.ProductName)
		return "", nil, fmt.Errorf("network node hasn't access to product")
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
