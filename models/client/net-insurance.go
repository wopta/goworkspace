package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models/dto/net"
)

const scope = "emettiPolizza_441-029-007"
const basePath = "https://apigatewaydigital.netinsurance.it"
const authEndpoint = "/Identity/connect/token"

type NetClient struct {
	*http.Client
}

func NewNetClient() *NetClient {
	client := &NetClient{}
	const tokenUrl = basePath + authEndpoint
	client.Client = lib.ClientCredentials(os.Getenv("NETINS_ID"), os.Getenv("NETINS_SECRET"), scope, tokenUrl)
	return client
}

func (c *NetClient) Quote(dto net.RequestDTO) (net.ResponseDTO, error) {
	var result net.ResponseDTO
	url := os.Getenv("NET_BASEURL") + "/PolizzeGateway24/emettiPolizza/441-029-007"
	rBuff := new(bytes.Buffer)
	err := json.NewEncoder(rBuff).Encode(dto)

	if err != nil {
		return result, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, rBuff)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errByte, _ := io.ReadAll(resp.Body)
		return result, errors.New(string(errByte))
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	if len(result.Errors) != 0 {
		return result, errors.New(fmt.Sprintln(result.Errors))
	}
	return result, nil
}
