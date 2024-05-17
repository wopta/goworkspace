package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wopta/goworkspace/lib"
)

type fabrickDeleteReq struct {
	Body        string `json:"body"` // `{"id":"27c44ee3-425f-4bca-88a6-9369efe9795b","newExpirationDate":"2023-07-14 00:00:00"}`
	ApiKey      string `json:"apiKey"`
	ApiKeyKey   string `json:"apiKeyKey"`   // "api-key"
	AuthSchema  string `json:"authSchema"`  // "S2S"
	ContentType string `json:"contentType"` // "application/json"
	XAuthToken  string `json:"xAuthToken"`
	Accept      string `json:"accept"` // application/json
}

func fabrickDelete(r *http.Request) (string, any, error) {
	var request fabrickDeleteReq
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(body, &request)

	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/payments/expirationDate"

	fabrickRequest, err := http.NewRequest(http.MethodPut, urlstring, strings.NewReader(request.Body))
	if err != nil {
		log.Printf("error creating request: %s", err.Error())
		return "", nil, err
	}

	fabrickRequest.Header.Set(request.ApiKeyKey, request.ApiKey)
	fabrickRequest.Header.Set("Auth-Schema", request.AuthSchema)
	fabrickRequest.Header.Set("Content-Type", request.ContentType)
	fabrickRequest.Header.Set("x-auth-token", request.XAuthToken)
	fabrickRequest.Header.Set("Accept", request.Accept)

	log.Printf("request: %v", fabrickRequest)

	res, err := lib.RetryDo(fabrickRequest, 5, 10)
	if err != nil {
		log.Printf("error getting response: %s", err.Error())
		return "", nil, err
	}

	log.Printf("reponse: %+v", res)

	respBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", nil, err
	}

	return string(respBody), nil, nil
}
