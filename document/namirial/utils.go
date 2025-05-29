package namirial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func doNamirialRequest(method, url string, body any) (*http.Request, error) {
	var (
		err error
		req *http.Request
	)

	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		requestJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestReader := bytes.NewReader(requestJson)
		req, err = http.NewRequest(method, url, requestReader)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("apiToken", os.Getenv("ESIGN_TOKEN_API"))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// check http response(status code) and unmarshal the body
func handleResponse[T any](r *http.Response, err error) (T, error) {
	var (
		req T
	)
	if err != nil {
		return req, err
	}

	if r.StatusCode != http.StatusOK {
		var body []byte
		r.Body.Read(body)
		return req, fmt.Errorf("ErrorNamirial: %s", string(body))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return *new(T), err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		return *new(T), err
	}
	return req, nil
}
