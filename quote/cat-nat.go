package quote

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/models/dto/net"
)

func CatNatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		reqPolicy *models.Policy
	)

	log.SetPrefix("[CatNatFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	_, err = lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	var cnReq net.RequestDTO
	err = cnReq.FromPolicy(reqPolicy)
	if err != nil {
		log.Printf("error building NetInsurance DTO: %s", err.Error())
		return "", nil, err
	}

	scope := "emettiPolizza_441-029-007"
	tokenUrl := "https://apigatewaydigital.netinsurance.it/Identity/connect/token"
	client := lib.ClientCredentials(os.Getenv("NETINS_ID"), os.Getenv("NETINS_SECRET"), scope, tokenUrl)

	resp, errResp, err := netInsuranceQuotation(client, cnReq)
	if err != nil {
		log.Printf("error calling NetInsurance api: %s", err.Error())
		return "", nil, err
	}
	var out []byte
	if errResp != nil {
		out, err = json.Marshal(errResp)
	} else {
		if resp.Result == "OK" {
			_ = resp.ToPolicy(reqPolicy)
			out, err = json.Marshal(reqPolicy)
		} else {
			out, err = json.Marshal(resp)
		}
	}

	if err != nil {
		log.Println("error encoding response %w", err.Error())
		return "", nil, err
	}

	return string(out), out, err
}

func netInsuranceQuotation(cl *http.Client, dto net.RequestDTO) (net.ResponseDTO, *net.ErrorResponse, error) {
	url := "https://apigatewaydigital.netinsurance.it/PolizzeGateway24/emettiPolizza/441-029-007"
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(dto)
	if err != nil {
		return net.ResponseDTO{}, nil, err
	}
	r := reqBodyBytes.Bytes()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(r))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.Do(req)
	if err != nil {
		return net.ResponseDTO{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := net.ErrorResponse{
			Errors: make(map[string]any),
		}
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Println("error decoding catnat error response")
			return net.ResponseDTO{}, nil, err
		}
		return net.ResponseDTO{}, &errResp, nil
	}
	cndto := net.ResponseDTO{}
	if err = json.NewDecoder(resp.Body).Decode(&cndto); err != nil {
		log.Println("error decoding catnat response")
		return net.ResponseDTO{}, nil, err
	}

	return cndto, nil, nil
}
