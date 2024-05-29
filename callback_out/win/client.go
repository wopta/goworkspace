package win

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/callback_out/internal"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type Client struct {
	basePath string
	path     string
	headers  map[string]string
}

func NewClient() *Client {
	return &Client{
		basePath: os.Getenv("WIN_CALLBACK_ENDPOINT"),
	}
}

func (c *Client) post(body io.Reader) (*http.Request, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", c.basePath, c.path)
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	res, err := lib.RetryDo(req, 5, 10)

	return req, res, err
}

func (c *Client) Emit(policy models.Policy) internal.CallbackInfo {
	c.path = "restba/extquote/inspratica"

	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA")
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

	body, err := inspratica(policy, "RICHIESTA_QUOTAZIONE")
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

	body, err := emissione(policy)
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
