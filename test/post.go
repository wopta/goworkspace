package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
)

func TestPostFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var request interface{}

	log.SetPrefix("[TestPostFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	operation := chi.URLParam(r, "operation")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(body, &request)
	log.Printf("payload %v", request)

	if operation == "error" {
		return "", nil, GetErrorJson(400, "Bad Request", "Testing error POST")
	}

	if operation == "fabrick-01" {
		var req fabrickTestRequest
		json.Unmarshal(body, &req)
		return fabrickGetPaymentInstruments(req)
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

type fabrickTestRequest struct {
	SubjectXId string `json:"subjectXId"`
}

func fabrickGetPaymentInstruments(req fabrickTestRequest) (string, interface{}, error) {
	var (
		woptaMerchantId string = "wop134b31-5926-4b26-1411-726bc9f0b111"
		token                  = os.Getenv("FABRICK_TOKEN_BACK_API")
	)

	urlstring := fmt.Sprintf("%spayment-instruments?merchantId=%s&subjectXId=%s&status=ACTIVE", os.Getenv("FABRICK_BASEURL"), woptaMerchantId, req.SubjectXId)

	request, err := http.NewRequest(http.MethodGet, urlstring, nil)
	if err != nil {
		log.Printf("error generating fabrick payment request: %s", err.Error())
		return "", nil, err
	}

	request.Header.Set("api-key", token)
	request.Header.Set("Auth-Schema", "S2S")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-auth-token", token)
	request.Header.Set("Accept", "application/json")

	res, err := lib.RetryDo(request, 5, 10)
	if err != nil || res == nil || res.StatusCode != 200 {
		log.Println(err)
		log.Println(res)
		return "error", nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", nil, err
	}

	return string(body), nil, err
}
