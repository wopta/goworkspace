package win

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type Client struct {
	basePath       string
	producer       string
	path           string
	headers        map[string]string
	externalConfig base.CallbackExternalConfig
}

func NewClient(producer string) *Client {
	return &Client{
		basePath: os.Getenv("WIN_CALLBACK_ENDPOINT"),
		producer: producer,
		// TODO: move me to external configuration
		externalConfig: base.CallbackExternalConfig{
			Events: map[base.CallbackoutAction]bool{
				base.Proposal:        true,
				base.RequestApproval: true,
				base.Emit:            true,
				base.Signed:          false,
				base.Paid:            true,
				base.EmitRemittance:  true,
				base.Approved:        false,
				base.Rejected:        false,
			},
			AuthType: "basic",
		},
	}
}

func (c *Client) post(body io.Reader) (*http.Request, *http.Response, error) {
	path := c.basePath + c.path
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err

	}
	req.SetBasicAuth(os.Getenv("WIN_CALLBACK_AUTH_USER"), os.Getenv("WIN_CALLBACK_AUTH_PASS"))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)

	return req, res, err
}

func (c *Client) Proposal(policy models.Policy) (callbackInfo base.CallbackInfo) {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA", c.producer)
	callbackInfo.ResAction = base.Emit
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}

	req, res, err := c.post(bytes.NewReader(body))
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}
	callbackInfo.FromRequestResponse(base.Proposal, res, req)
	return callbackInfo
}

func (c *Client) Emit(policy models.Policy) (callbackInfo base.CallbackInfo) {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA", c.producer)
	callbackInfo.ResAction = base.Emit
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}

	req, res, err := c.post(bytes.NewReader(body))
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}
	callbackInfo.FromRequestResponse(base.Emit, res, req)
	return callbackInfo
}

func (c *Client) RequestApproval(policy models.Policy) (callbackInfo base.CallbackInfo) {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "RICHIESTA_QUOTAZIONE", c.producer)
	callbackInfo.ResAction = base.RequestApproval
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}

	req, res, err := c.post(bytes.NewReader(body))
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}
	callbackInfo.FromRequestResponse(base.RequestApproval, res, req)
	return callbackInfo
}

func (c *Client) Paid(policy models.Policy) (callbackInfo base.CallbackInfo) {
	c.path = "restba/extquote/emissione"

	body, err := emissione(policy, c.producer)
	callbackInfo.ResAction = base.Paid
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}

	req, res, err := c.post(bytes.NewReader(body))
	if err != nil {
		callbackInfo.Error = err
		return callbackInfo
	}
	callbackInfo.FromRequestResponse(base.Paid, res, req)
	return callbackInfo
}

func (c *Client) Signed(models.Policy) base.CallbackInfo {
	return base.CallbackInfo{ResAction: base.Signed}
}

func (c *Client) Approved(models.Policy) base.CallbackInfo {
	return base.CallbackInfo{ResAction: base.Approved}
}

func (c *Client) Rejected(models.Policy) base.CallbackInfo {
	return base.CallbackInfo{ResAction: base.Rejected}
}

func (c *Client) DecodeAction(rawAction base.CallbackoutAction) []base.CallbackoutAction {
	actionEnabled, ok := c.externalConfig.Events[rawAction]
	if !actionEnabled || !ok {
		return nil
	}

	availableActions := base.GetAvailableActions()
	decodedActions := availableActions[rawAction]

	return decodedActions
}
