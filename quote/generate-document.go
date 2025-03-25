package quote

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/product"
)

type DraftResponse struct {
	RawDoc string `json:"rawDoc"`
}

func GenerateDocumentFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err           error
		policy        models.Policy
		nn            *models.NetworkNode
		w             *models.Warrant
		docBytes      []byte
		responseBytes []byte
	)

	log.SetPrefix("[GenerateQuoteDocumentFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	if err = json.NewDecoder(r.Body).Decode(&policy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	policy.Normalize()

	if nn = network.GetNetworkNodeByUid(policy.ProducerUid); nn != nil {
		w = nn.GetWarrant()
	}

	prd := product.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nn, w)

	if docBytes, err = document.Quote(&policy, prd); err != nil {
		log.Println("error generating document")
		return "", nil, err
	}

	response := DraftResponse{
		RawDoc: base64.StdEncoding.EncodeToString(docBytes),
	}

	if responseBytes, err = json.Marshal(response); err != nil {
		log.Println("error marshaling response")
		return "", nil, err
	}

	return string(responseBytes), response, err
}
