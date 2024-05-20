package win

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type Client struct {
	path    string
	headers map[string]string
}

func (c *Client) post(body io.Reader) (*http.Request, *http.Response, error) {
	// inject fixed base url
	req, err := http.NewRequest(http.MethodPost, c.path, body)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	res, err := lib.RetryDo(req, 5, 10)

	return req, res, err
}

func (c *Client) Emit(policy models.Policy) (*http.Request, *http.Response, error) {
	c.path = "/restba/extquote/inspratica"

	body, err := inspratica(policy, "QUOTAZIONE_ACCETTATA")
	if err != nil {
		return nil, nil, err
	}

	return c.post(bytes.NewReader(body))
}

func (c *Client) RequestApproval(policy models.Policy) (*http.Request, *http.Response, error) {
	c.path = "/restba/extquote/inspratica"

	body, err := inspratica(policy, "RICHIESTA_QUOTAZIONE")
	if err != nil {
		return nil, nil, err
	}

	return c.post(bytes.NewReader(body))
}

func (c *Client) Paid(policy models.Policy) (*http.Request, *http.Response, error) {
	c.path = "/restba/extquote/emissione"

	body, err := emissione(policy)
	if err != nil {
		return nil, nil, err
	}

	return c.post(bytes.NewReader(body))
}
