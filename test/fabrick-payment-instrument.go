package test

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"

	"github.com/wopta/goworkspace/lib"
)

type FabrickPaymentInstrumentReq struct {
	BaseUrl     string `json:"baseUrl"`
	ApiKey      string `json:"apiKey"`
	ApiKeyKey   string `json:"apiKeyKey"`
	XAuthToken  string `json:"xAuthToken"`
	Accept      string `json:"accept"`
	AuthSchema  string `json:"authSchema"`
	MerchantId  string `json:"merchantId"`
	SubjectXId  string `json:"subjectXId"`
	Status      string `json:"status"`
	ContentType string `json:"contentType"`
}

func fabrickPaymentInstrument(r *http.Request) (string, interface{}, error) {
	var (
		err  error
		resp map[string]interface{}
		body FabrickPaymentInstrumentReq
	)

	log.AddPrefix("FabrickPaymentInstrument")
	defer log.PopPrefix()
	log.Println("Handler start -----------------------------------------------")

	rawBody := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}

	url := fmt.Sprintf("%s/payment-instruments?merchantId=%s&subjectXId=%s&status=%s", body.BaseUrl,
		body.MerchantId, body.SubjectXId, body.Status)

	log.Printf("url: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set(body.ApiKeyKey, body.ApiKey)
	req.Header.Set("x-auth-token", body.XAuthToken)
	req.Header.Set("Accept", body.Accept)
	req.Header.Set("Auth-Schema", body.AuthSchema)
	req.Header.Set("Content-Type", body.ContentType)

	log.Printf("request: %v", req)

	res, err := lib.RetryDo(req, 5, 100)
	if err != nil {
		log.Println(err.Error())
	}

	fabrickStatus := res.StatusCode
	fabrickRespBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", nil, err
	}

	log.Printf("response with status %d body: %s", fabrickStatus, string(fabrickRespBody))

	err = json.Unmarshal(fabrickRespBody, &resp)

	return string(fabrickRespBody), resp, err
}
