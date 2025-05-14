package catnat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
)

const authEndpoint = "/Identity/connect/token"
const scope = "emettiPolizza_441-029-007"

type NetClient struct {
	*http.Client
}

type INetClient interface {
	Quote(dto RequestDTO) (response ResponseDTO, err error)
	Emit(dto RequestDTO) (response any, err error)
}

func NewNetClient() (client *NetClient) {
	client = &NetClient{}
	tokenUrl := os.Getenv("NET_BASEURL") + authEndpoint
	client.Client = lib.ClientCredentials(os.Getenv("NETINS_ID"), os.Getenv("NETINS_SECRET"), scope, tokenUrl)
	return client
}

func (c *NetClient) Quote(dto RequestDTO) (response ResponseDTO, err error) {
	url := os.Getenv("NET_BASEURL") + "/PolizzeGateway24/emettiPolizza/441-029-007"
	rBuff := new(bytes.Buffer)
	err = json.NewEncoder(rBuff).Encode(dto)

	if err != nil {
		return response, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, rBuff)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return response, err
		}
		return response, errors.New(string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	if response.Result != "OK" {
		return response, fmt.Errorf("%+v", response.Errors)
	}

	return response, nil
}

func (c *NetClient) Emit(dto RequestDTO) (response any, err error) {
	//TODO: to fix this
	dto.Emission = "si"
	url := os.Getenv("NET_BASEURL") + "/PolizzeGateway24/emettiPolizza/441-029-007"
	rBuff := new(bytes.Buffer)
	err = json.NewEncoder(rBuff).Encode(dto)

	if err != nil {
		return response, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, rBuff)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return response, err
		}
		return response, errors.New(string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	//TODO: to adpta when emit catnat will be merged
	//	if response.Result != "OK" {
	//		return response, fmt.Errorf("%+v", response.Errors)
	//	}

	return response, nil
}
