package fabrick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

type refreshTokenReq struct {
	Username        string `json:"username,omitempty"`
	AgencyExtId     string `json:"agencyExtId"`
	SubAgencyExtId  string `json:"subAgencyExtId,omitempty"`
	VisibilityExtId string `json:"visibilityExtId,omitempty"`
	MerchantExtId   string `json:"merchantExtId,omitempty"`
	UserId          string `json:"userId"`
	Name            string `json:"name"`
	Surname         string `json:"surname,omitempty"`
	Role            string `json:"role"`
	Node            string `json:"node,omitempty"`
	TokenHash       string `json:"tokenHash,omitempty"`
	IpRequest       string `json:"ipRequest,omitempty"`
	Tenant          string `json:"tenant"`
	Merchant        string `json:"merchant"`
}

var requestBody refreshTokenReq = refreshTokenReq{
	AgencyExtId: "test agency",
	UserId:      os.Getenv("FABRICK_MERCHANT_ID"),
	Name:        "WOPTA",
	Role:        "USER",
	Tenant:      os.Getenv("FABRICK_REFRESH_TENANT"),
	Merchant:    os.Getenv("FABRICK_MERCHANT_ID"),
}

func RefreshTokenFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.AddPrefix("RefreshTokenFx")
	log.Println("Handler Start -----------------------------------------------")
	defer func() {
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	err := refreshToken()
	if err != nil {
		log.Error(err)
	}

	return "", nil, err
}

func refreshToken() error {
	var (
		urlstring = os.Getenv("FABRICK_BASEURL") + "api/fabrick/pace/v4.0/mods/back/v1.0/sessions/createSession"
		body      []byte
		err       error
	)

	if body, err = json.Marshal(requestBody); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("api-key", os.Getenv("FABRICK_TOKEN_BACK_API"))
	req.Header.Set("Auth-Schema", "S2S")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	log.Printf("executing request: %+v", req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Printf("got response: %+v", res)

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("response status code: %d", res.StatusCode)
	}

	return nil
}
