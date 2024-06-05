package win

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/models"
)

type Client struct {
	basePath string
	producer string
	path     string
	headers  map[string]string
}

func NewClient(producer string) *Client {
	return &Client{
		basePath: os.Getenv("WIN_CALLBACK_ENDPOINT"),
		producer: producer,
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

func (c *Client) Proposal(policy models.Policy) internal.CallbackInfo {
	return c.Emit(policy)
}

func (c *Client) Emit(policy models.Policy) internal.CallbackInfo {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA", c.producer)
	if err != nil {
		return internal.CallbackInfo{
			Request:     nil,
			RequestBody: nil,
			Response:    nil,
			Error:       err,
		}
	}

	req, res, err := c.post(bytes.NewReader(body))
	return internal.CallbackInfo{
		Request:     req,
		RequestBody: body,
		Response:    res,
		Error:       err,
	}
}

func (c *Client) RequestApproval(policy models.Policy) internal.CallbackInfo {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "RICHIESTA_QUOTAZIONE", c.producer)
	if err != nil {
		return internal.CallbackInfo{
			Request:     nil,
			RequestBody: nil,
			Response:    nil,
			Error:       err,
		}
	}

	req, res, err := c.post(bytes.NewReader(body))
	return internal.CallbackInfo{
		Request:     req,
		RequestBody: body,
		Response:    res,
		Error:       err,
	}
}

func (c *Client) Paid(policy models.Policy) internal.CallbackInfo {
	c.path = "restba/extquote/emissione"

	body, err := emissione(policy, c.producer)
	if err != nil {
		return internal.CallbackInfo{
			Request:     nil,
			RequestBody: nil,
			Response:    nil,
			Error:       err,
		}
	}

	req, res, err := c.post(bytes.NewReader(body))
	return internal.CallbackInfo{
		Request:     req,
		RequestBody: body,
		Response:    res,
		Error:       err,
	}
}
