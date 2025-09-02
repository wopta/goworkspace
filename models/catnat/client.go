package catnat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"golang.org/x/oauth2/clientcredentials"
)

const authEndpoint = "/Identity/connect/token"

type NetClient struct {
	*http.Client
}

type INetClient interface {
	Quote(dto QuoteRequest, policy *models.Policy) (response QuoteResponse, err error)
	Emit(dto QuoteRequest, policy *models.Policy) (response QuoteResponse, err error)
	Download(numeroPoliza string) (response DownloadResponse, err error)
}

func NewNetClient() (client *NetClient) {
	client = &NetClient{}
	tokenUrl := os.Getenv("NET_BASEURL") + authEndpoint
	config := clientcredentials.Config{
		ClientID:     os.Getenv("NETINS_ID"),
		ClientSecret: os.Getenv("NETINS_SECRET"),
		Scopes:       []string{"emettiPolizza_441-006-006", "emettiPolizza_441-027-056", "emettiPolizza_441-029-009", "IncassaTitolo_441", "InserisciAllegato_441", "StampaPolizza_441"},
		TokenURL:     tokenUrl,
		EndpointParams: url.Values{
			"grant_type": {"client_credentials"},
		}}

	client.Client = config.Client(context.Background())
	return client
}

func (c *NetClient) Quote(dto QuoteRequest, policy *models.Policy) (response QuoteResponse, err error) {
	response, err = c.quote(dto)
	if err != nil {
		return response, err
	}
	err = mappingQuoteResponseToPolicy(response, policy)
	if err != nil {
		return response, err
	}
	err = mappingQuoteResponseToGuarantee(response, policy)
	return response, err
}

func (c *NetClient) Emit(dto QuoteRequest, policy *models.Policy) (response QuoteResponse, err error) {
	response, err = c.emit(dto)
	if err != nil {
		return response, err
	}
	err = mappingQuoteResponseToPolicy(response, policy)
	if err != nil {
		return response, err
	}
	err = mappingQuoteResponseToGuarantee(response, policy)
	return response, err
}
func (c *NetClient) quote(dto QuoteRequest) (response QuoteResponse, err error) {
	url := os.Getenv("NET_BASEURL") + "/PolizzeGateway24/emettiPolizza/441-029-009"
	rBuff := new(bytes.Buffer)
	log.PrintStruct("request: ", dto)
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
		return response, errors.New(resp.Status + ":" + string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	if response.Result != "OK" {
		log.ErrorF("Errore quotazione %+v", response.Errors)
		for i := range response.Errors {
			if strings.Contains(response.Errors[i].Description, "Importo premio inferiore al premio minimo consentito") {
				return response, errors.New("Il premio non può essere inferiore a 100 euro annui, aumenta le somme assicurate per raggiungere il premio minimo.")
			}
			if strings.Contains(response.Errors[i].Description, "Indirizzo immobile non valido") || strings.Contains(response.Errors[i].Description, "Provincia bene assicurato non trovata") {
				return response, errors.New("Indirizzo immobile non valido")
			}
		}
		return response, errors.New("Errore quotazione")
	}
	return response, nil
}

func (c *NetClient) emit(dto QuoteRequest) (response QuoteResponse, err error) {
	dto.Emission = "si"
	url := os.Getenv("NET_BASEURL") + "/PolizzeGateway24/emettiPolizza/441-029-009"
	rBuff := new(bytes.Buffer)
	log.PrintStruct("request interagration api netensurance", dto)
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
		return response, errors.New(resp.Status + ":" + string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	log.PrintStruct("response api catnat", response)
	if response.Result != "OK" {
		log.ErrorF("Errore emissione %+v", response.Errors)
		return response, errors.New("Errore emissione")
	}
	return response, nil
}

func (c *NetClient) Download(numeroPolizza string) (response DownloadResponse, err error) {
	url := os.Getenv("NET_BASEURL") + "/OperationsGateway/StampaPolizzaAppendice/StampaPolizza"
	rBuff := new(bytes.Buffer)
	request := DownloadRequest{
		CodiceCompagnia: "441",
		NumeroPolizza:   numeroPolizza,
		TipoOperazione:  "E",
		DataOperazione:  time.Now(),
	}
	err = json.NewEncoder(rBuff).Encode(request)
	log.PrintStruct("request download", request)
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
		return response, errors.New(resp.Status + ":" + string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	if response.Result != "OK" {
		log.ErrorF("Errore emissione %+v", response.Errors)
		return response, errors.New("Error download")
	}
	return response, nil
}
func CalculateControlNumberFromVatCode(vatCode string) string {
	x := 0
	y := 0
	i := 0
	for {
		if i == 10 {
			break
		}
		if i%2 != 0 {
			y += (int(vatCode[i]) - 48) * 2
			if (int(vatCode[i])-48)*2 > 9 {
				y -= 9
			}
		} else {
			x += (int(vatCode[i]) - 48)
		}

		i++
	}
	sum := x + y
	sum = sum % 10
	sum = (10 - sum) % 10
	return fmt.Sprint(sum)
}
func (c *NetClient) EnrichAteco(fiscalCode string) (response AtecoResponse, err error) {
	if len(fiscalCode) == 11 { //partita iva
		checkDigit := CalculateControlNumberFromVatCode(fiscalCode)
		if checkDigit != string(fiscalCode[10]) {
			return response, errors.New("codice ateco non corretto, verificare il codice di controllo")
		}
	}
	url := os.Getenv("NET_BASEURL") + "/OperationsGateway/InformazioniCompagnia/" + fiscalCode
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}
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
		return response, errors.New(resp.Status + ":" + string(errBytes))
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.ErrorF("error decoding catnat response")
		return response, err
	}
	if response.Result != "OK" {
		return response, errors.New("Error Enrich ateco")
	}
	//Formatting ateco code
	var ateco string
	hasDot := func(i int) bool {
		if i == 0 {
			return false
		}
		return i%2 == 0
	}
	for i := range response.Ateco {
		if hasDot(i) {
			ateco += "."
		}
		ateco += string(response.Ateco[i])
	}
	response.Ateco = ateco
	return response, err
}
