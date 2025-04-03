package test

import (
	"bytes"
	"encoding/json"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
)

type fabrickPersistentTokenReq struct {
	ApiKey      string `json:"apiKey"`
	ApiKeyKey   string `json:"apiKeyKey"`   // "api-key"
	AuthSchema  string `json:"authSchema"`  // "S2S"
	ContentType string `json:"contentType"` // "application/json"
	XAuthToken  string `json:"xAuthToken"`
	Accept      string `json:"accept"` // application/json
	Body        string `json:"body"`
}

func fabrickPersistentToken(r *http.Request) (string, interface{}, error) {
	var (
		err  error
		body fabrickPersistentTokenReq
	)

	log.AddPrefix("[fabrickPersistentToken] ")
	defer log.PopPrefix()
	log.Println("Handler Start -----------------------------------------------")

	rawBody := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		log.Printf("Error Unmarshal Request Body: %s\n", err.Error())
		return "", nil, err
	}

	var urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/sessions/createSession"

	req, err := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer([]byte(body.Body)))
	if err != nil {
		log.Printf("Error request body: %s\n", err.Error())
		return "", nil, err
	}

	req.Header.Set(body.ApiKeyKey, body.ApiKey)
	req.Header.Set("Auth-Schema", body.AuthSchema)
	req.Header.Set("Content-Type", body.ContentType)
	req.Header.Set("Accept", body.Accept)
	//req.Header.Set("x-auth-token", body.XAuthToken)

	log.Printf("request: %+v", req)

	res, err := lib.RetryDo(req, 5, 10)
	if err != nil {
		log.Printf("Error getting response: %s\n", err.Error())
		return "", nil, err
	}

	log.Printf("response: %+v", res)

	rawResp, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error Read Response Body: %s\n", err.Error())
		return "", nil, err
	}

	log.Printf("rawResp: %+v", string(rawResp))

	log.Println("Handler End -------------------------------------------------")

	return string(rawResp), res, nil
}
