package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wopta/goworkspace/callback_out/internal"
	md "github.com/wopta/goworkspace/callback_out/models"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type Client struct {
	basePath       string
	producer       string
	externalConfig internal.CallbackExternalConfig
	network        string
}

func NewClient(networkNode *models.NetworkNode, network string) *Client {
	basePath := os.Getenv(fmt.Sprintf("%s_CALLBACK_ENDPOINT", lib.ToUpper(network)))
	if basePath == "" {
		return nil
	}

	var externalConfig internal.CallbackExternalConfig
	// configBytes := lib.GetFilesByEnv("/callback-out/base.json")
	configBytes := []byte(`{
		"events": {
			"Proposal": true,
			"RequestApproval": false,
			"Emit": true,
			"Paid": true,
			"EmitRemittance": false
		},
		"authType": "basic"
	}`)
	if err := json.Unmarshal(configBytes, &externalConfig); err != nil {
		return nil
	}

	return &Client{
		basePath:       basePath,
		producer:       networkNode.Code,
		externalConfig: externalConfig,
		network:        network,
	}
}

func (c *Client) baseRequest(policy models.Policy) internal.CallbackInfo {
	rawBody, err := json.Marshal(policy)
	if err != nil {
		return internal.CallbackInfo{
			Request:     nil,
			RequestBody: nil,
			Response:    nil,
			Error:       err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, c.basePath, bytes.NewReader(rawBody))
	if err != nil {
		return internal.CallbackInfo{
			Request:     nil,
			RequestBody: nil,
			Response:    nil,
			Error:       err,
		}
	}

	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)

	return internal.CallbackInfo{
		Request:     req,
		RequestBody: rawBody,
		Response:    res,
		Error:       err,
	}
}

func (c *Client) Proposal(policy models.Policy) internal.CallbackInfo {
	return c.baseRequest(policy)
}

func (c *Client) Emit(policy models.Policy) internal.CallbackInfo {
	return c.baseRequest(policy)
}

func (c *Client) RequestApproval(policy models.Policy) internal.CallbackInfo {
	return c.baseRequest(policy)
}

func (c *Client) Paid(policy models.Policy) internal.CallbackInfo {
	return c.baseRequest(policy)
}

func (c *Client) DecodeAction(rawAction string) []string {
	actionEnabled, ok := c.externalConfig.Events[rawAction]
	if !actionEnabled || !ok {
		return nil
	}

	availableActions := md.GetAvailableActions()
	decodedActions := availableActions[rawAction]

	return decodedActions
}

func (c *Client) setAuth(req *http.Request) {
	network := lib.ToUpper(c.network)
	switch c.externalConfig.AuthType {
	case "basicAuth":
		req.SetBasicAuth(
			os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_USER", network)),
			os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_PASS", network)))
	case "api-key":
		req.Header.Add("Authorization", os.Getenv(fmt.Sprintf("%s_CALLBACK_AUTH_APIKEY", network)))
	}
}
