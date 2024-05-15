package win

import (
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
)

type winClient struct {
	path    string
	headers map[string]string
}

func (c *winClient) Post(body io.Reader) (*http.Request, *http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.path, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("key", "value")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	log.Printf("win request: %v", req)

	res, err := lib.RetryDo(req, 5, 10)

	log.Printf("win response: %v", res)

	return req, res, err
}
