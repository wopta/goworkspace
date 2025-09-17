package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type Client struct {
	basePath       string
	producer       string
	externalConfig CallbackExternalConfig
	network        string
}

func NewClient(networkNode *models.NetworkNode, network string) *Client {
	basePath := os.Getenv(fmt.Sprintf("%s_CALLBACK_ENDPOINT", lib.ToUpper(network)))
	if basePath == "" {
		return nil
	}

	var externalConfig CallbackExternalConfig
	configBytes := lib.GetFilesByEnv("callback-out/base.json")
	if err := json.Unmarshal(configBytes, &externalConfig); err != nil {
		return nil
	}

	return &Client{
		basePath:       basePath,
		producer:       networkNode.Code,
		network:        network,
		externalConfig: externalConfig,
	}
}

func (c *Client) baseRequest(policy models.Policy, action CallbackoutAction) (callbackInfo CallbackInfo) {
	rawBody, err := json.Marshal(policy)
	callbackInfo.ResAction = action
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo

	}

	req, err := http.NewRequest(http.MethodPost, c.basePath, bytes.NewReader(rawBody))
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	//TODO: insert action
	callbackInfo.FromRequestResponse(action, res, req)
	return callbackInfo
}

func (c *Client) Proposal(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Proposal)
}

func (c *Client) Emit(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Emit)
}

func (c *Client) RequestApproval(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, RequestApproval)
}

func (c *Client) Paid(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Paid)
}

func (c *Client) Signed(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Signed)
}

func (c *Client) Approved(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Approved)
}

func (c *Client) Rejected(policy models.Policy) CallbackInfo {
	return c.baseRequest(policy, Rejected)
}

func (c *Client) DecodeAction(rawAction CallbackoutAction) []CallbackoutAction {
	actionEnabled, ok := c.externalConfig.Events[rawAction]
	if !actionEnabled || !ok {
		return nil
	}

	availableActions := GetAvailableActions()
	decodedActions := availableActions[rawAction]

	return decodedActions
}

func (c *Client) setAuth(req *http.Request) {
	network := lib.ToUpper(c.network)
	switch c.externalConfig.AuthType {
	case "basic":
		req.SetBasicAuth(
			os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_USER", network)),
			os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_PASS", network)))
	}
}
