package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models/dto/net"
)

type Client struct {
	httpC *http.Client
}

func NewNetClient() *Client {
	return &Client{}
}

func (c *Client) Authenticate() {
	const scope = "emettiPolizza_441-029-007"
	const basePath = "https://apigatewaydigital.netinsurance.it"
	const authEndpoint = "/Identity/connect/token"
	const tokenUrl = basePath + authEndpoint
	c.httpC = lib.ClientCredentials(os.Getenv("NETINS_ID"), os.Getenv("NETINS_SECRET"), scope, tokenUrl)
}

func (c *Client) Quote(dto net.RequestDTO) (net.ResponseDTO, *net.ErrorResponse, error) {
	url := "https://apigatewaydigital.netinsurance.it/PolizzeGateway24/emettiPolizza/441-029-007"
	rBuff := new(bytes.Buffer)
	err := json.NewEncoder(rBuff).Encode(dto)
	if err != nil {
		return net.ResponseDTO{}, nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, rBuff)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpC.Do(req)
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

func (c *Client) Emit() error {

	return nil
}
