package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

type InitReq struct {
	ProductName string `json:"productName"`
}

func InitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
		req InitReq
	)

	log.AddPrefix("InitFx")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		log.Println("Handler end ----------------------------------------------")
		log.PopPrefix()
	}()

	log.Println("Handler start -----------------------------------------------")

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	product = prd.GetLatestActiveProduct(req.ProductName, authToken.GetChannelByRoleV2(), nil, nil)
	if product == nil {
		return "", nil, fmt.Errorf("product %s not found", req.ProductName)
	}

	now := time.Now().UTC()
	channel := authToken.GetChannelByRoleV2()

	policy := models.Policy{
		Uid:            lib.NewDoc(lib.PolicyCollection),
		Annuity:        0,
		Channel:        channel,
		Company:        product.Companies[0].Name,
		CreationDate:   now,
		IsAutoRenew:    product.IsAutoRenew,
		IsRenewable:    product.IsRenewable,
		Name:           product.Name,
		NameDesc:       *product.NameDesc,
		PolicyType:     product.PolicyType,
		ProductVersion: product.Version,
		QuoteType:      product.QuoteType,
		Status:         models.PolicyStatusInit,
		StatusHistory:  []string{models.PolicyStatusInit},
		Updated:        now,
	}

	if channel == lib.NetworkChannel {
		networkNode = network.GetNetworkNodeByUid(authToken.UserID)
		if networkNode == nil {
			return "", nil, fmt.Errorf("network node %s not found", authToken.UserID)
		}
		policy.ProducerCode = networkNode.Code
		policy.ProducerType = networkNode.Type
		policy.ProducerUid = networkNode.Uid
	}

	err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
	if err != nil {
		return "", nil, err
	}

	policy.BigquerySave()

	rawPolicy, err := json.Marshal(policy)

	return string(rawPolicy), policy, err
}
