package win

import (
	"bytes"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
)

type winClient struct {
	path string
}

func (c *winClient) Post(body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", "value")
	log.Printf("win request: %v", req)

	return lib.RetryDo(req, 5, 10)
}
