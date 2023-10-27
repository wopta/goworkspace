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
		if networkNode != nil {
			warrant = networkNode.GetWarrant()
		}
		if warrant != nil && !networkNode.HasAccessToProduct(req.ProductName, warrant) {
			log.Printf("[GetProductByChannelFx] network node %s hasn't access to product %s", networkNode.Uid, req.ProductName)
			return "", nil, fmt.Errorf("network node hasn't access to product")
		}
	}

	log.Printf("[GetProductByChannelFx] getting last active action for product %s", req.ProductName)

	product := prd.GetLatestActiveProduct(req.ProductName, channel, networkNode, warrant)
	if product == nil {
		log.Printf("[GetProductByChannelFx] no active product found")
		return "", nil, fmt.Errorf("no product active found")
	}

	product.Steps = filterProductSteps(product, warrant)

	jsonOut, err := product.Marshal()

	log.Println("[GetProductByChannelFx] handler end -------------")

	return string(jsonOut), product, err
}

func filterProductSteps(product *models.Product, warrant *models.Warrant) []models.Step {
	var steps []models.Step
	rawSteps := lib.GetFilesByEnv(fmt.Sprintf("products-v2/%s/%s/builder_ui.json", product.Name, product.Version))
	_ = json.Unmarshal(rawSteps, &steps)

	outputSteps := make([]models.Step, 0)
	for _, step := range steps {
		if len(step.Flows) == 0 || lib.SliceContains(step.Flows, warrant.GetFlowName(product.Name)) {
			outputSteps = append(outputSteps, step)
		}
	}
	return outputSteps
}
